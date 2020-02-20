package brainwave_yagpdb
import (
        "encoding/json"
        "fmt"
        "strings"

	"github.com/jbowens/codenames"
	"github.com/jinzhu/gorm"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/common"
)

/// The captains are Discord userids
type captains struct {
        Red  int64
        Blue int64
}

type Game struct {
	gorm.Model
	GuildID   int64           `gorm:"primary_key;auto_increment:false"`
	ChannelID int64           `gorm:"primary_key;auto_increment:false"`
        Captains  captains        `gorm:"embedded;embedded_prefix:captain_"`
	Game      codenames.Game  `gorm:"-"`
	GameSave  json.RawMessage `sql:"type:json"`
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
        if g.Captains.Red == 0 {
                return "You can't start a game without both captains. `bw lead red`?", nil
        }
        if g.Captains.Blue == 0 {
                return "You can't start a game without both captains. `bw lead blue`?", nil
        }

        // TODO: Words would be useful...
        g.Game = *codenames.NewGame([]string{})

        return fmt.Sprintf("%s starts!", g.Game.CurrentTeam()), nil
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

        return fmt.Sprintf("%s is now the %s captain", user.Username, t), nil
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
        for i, bWord := range(g.Game.Words) {
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
        if err != nil && err != gorm.ErrRecordNotFound{
                logger.WithError(err).Errorf("failed to load game g[%d] c[%d]", guildID, channelID)
                return nil, err
        }

        return &g, nil
}

