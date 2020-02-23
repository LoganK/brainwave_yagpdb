package brainwave_yagpdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/jbowens/codenames"
	"github.com/jinzhu/gorm"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/common"
)

var (
	// We make a best-effort attempt to clean up after ourselves. Maps
	// channels to sent messages. (We don't persist this because it greatly
	// complicates the logic.)
	// TODO: Avoid slow, but boundless, growth over time.
	sentBoardMessages = make(map[int64]int64)
)

// The captains are Discord userids but stored as strings for future migration
// back to the core library.
type Captains struct {
	Red  int64
	Blue int64
}

type Game struct {
	gorm.Model
	GuildID   int64 `gorm:"primary_key;auto_increment:false"`
	ChannelID int64 `gorm:"primary_key;auto_increment:false"`

	// Data that should probably be moved into codenames.Game
	Game     codenames.Game  `gorm:"-"`
	GameSave json.RawMessage `sql:"type:json"`
	Captains Captains        `gorm:"embedded;embedded_prefix:captain_"`
}

var (
	ErrNoCaptains = errors.New("game is missing a captain")
)

func (g *Game) Start(wordList []string) error {
	if g.Captains.Red == 0 || g.Captains.Blue == 0 {
		return ErrNoCaptains
	}

	// TODO: Words would be useful...
	g.Game = *codenames.NewGame(wordList)

	return nil
}

func (g *Game) TableName() string {
	return "brainwave_games"
}

func (g *Game) BeforeSave() error {
	gameObj, err := json.Marshal(g.Game)
	if err != nil {
		logger.WithError(err).Errorf("failed to save game")
		return err
	}

	g.GameSave = gameObj
	return nil
}

func (g *Game) AfterFind() {
	if err := json.Unmarshal(g.GameSave, &g.Game); err != nil {
		logger.WithError(err).Errorf("failed to restore game")
		return
	}

	g.GameSave = json.RawMessage{}
}

func (g *Game) runStart() (interface{}, error) {
	if err := g.Start(defaultWordsEnUs); err != nil {
		if err == ErrNoCaptains {
			return template.Must(template.New("NoCaptainStart").
				Parse("You can't start a game without both captains. `{{.Data.Keyword}} lead (red|blue)`?")), nil
		}
	}

	return g.updateBoard()
}

func (g *Game) updateBoard() (interface{}, error) {
	asJpeg := func(v codenames.Viewer) io.Reader {
		img, err := g.Game.RenderGameBoard(v)
		if err != nil {
			logger.WithError(err).Errorf("failed to render board")
			return nil
		}

		var b bytes.Buffer
		if err := jpeg.Encode(&b, img, &jpeg.Options{90}); err != nil {
			logger.WithError(err).Errorf("failed to encode board")
			return nil
		}

		return &b
	}

	captainBoard := make(chan io.Reader)
	defer close(captainBoard)
	go func() { captainBoard <- asJpeg(codenames.Spymaster) }()

	playerBoard := make(chan io.Reader)
	go func() { playerBoard <- asJpeg(codenames.Player) }()
	defer close(playerBoard)

	// We don't cache the private channels because the captains may have changed.
	redChn, err := common.BotSession.UserChannelCreate(g.Captains.Red)
	if err != nil {
		logger.WithError(err).Errorf("failed to create private channel")
	}
	blueChn, err := common.BotSession.UserChannelCreate(g.Captains.Blue)
	if err != nil {
		logger.WithError(err).Errorf("failed to create private channel")
	}

	turnMsg := fmt.Sprintf("It's %s's turn!", g.Game.CurrentTeam())

	sendBoard := func(channelID int64, r io.Reader) {
		if msgID := sentBoardMessages[channelID]; msgID != 0 {
			// Clear out any boards we've sent. This isn't necessary, but let's try
			// not to fill up all of Discord's hard drives.
			common.BotSession.ChannelMessageDelete(channelID, msgID)
		}
		msg, _ := common.BotSession.ChannelMessageSendComplex(
			channelID, &discordgo.MessageSend{
				File: &discordgo.File{Reader: r}, Content: turnMsg})
		sentBoardMessages[channelID] = msg.ID
	}

	captainBytes, _ := ioutil.ReadAll(<-captainBoard)
	sendBoard(redChn.ID, bytes.NewReader(captainBytes))
	sendBoard(blueChn.ID, bytes.NewReader(captainBytes))

	sendBoard(g.ChannelID, <-playerBoard)

	return nil, nil
}

func (g *Game) runLead(user *discordgo.User, team string) (interface{}, error) {
	// Let's be intentionally loose with the parsing so people can have fun.
	var t codenames.Team
	switch team[0] {
	case 'r':
		g.Captains.Red = user.ID
		t = codenames.Red
	case 'b':
		g.Captains.Blue = user.ID
		t = codenames.Blue
	default:
		return "You must specify red or blue.", nil
	}

	return fmt.Sprintf("{{.Msg.Author.Username}} is now the %s captain", t), nil
}

func (g *Game) runTouch(word string) (interface{}, error) {
	if len(g.Game.Words) == 0 {
		return "Love the enthusiasm, but there's no active game.", nil
	}
	if len(word) == 0 {
		return "You've touched my heart. But not a word. Because you didn't specify one.", nil
	}
	word = strings.ToLower(word)

	// TODO: Check the player is on the correct team.

	var wordIdx int
	for i, bWord := range g.Game.Words {
		if word == strings.ToLower(bWord) {
			wordIdx = i
			break
		}
	}

	err := g.Game.Guess(wordIdx)
	if err != nil {
		if err == codenames.ErrAlreadyRevealed {
			return fmt.Sprintf("You touch it again but nothing changes. Still your turn, %s.", g.Game.CurrentTeam()), nil
		}
		logger.WithError(err).Error("failed to touch word")
		return "Something went terribly wrong...", err
	}

	// TODO: Report the result.

	return nil, nil
}

func loadGameFromDB(guildID, channelID int64) (*Game, error) {
	var g Game
	err := common.GORM.Where(&Game{GuildID: guildID, ChannelID: channelID}).
		First(&g).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.WithError(err).Errorf("failed to load game g[%d] c[%d]", guildID, channelID)
		return nil, err
	}

	return &g, nil
}
