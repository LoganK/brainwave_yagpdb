package brainwave_yagpdb
import (
        "encoding/json"

	"github.com/jbowens/codenames"
	"github.com/jinzhu/gorm"
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

        return nil, nil
}

func (g *Game) runLead(team string) (interface{}, error) {
	return nil, nil
}

func (g *Game) runTouch(word string) (interface{}, error) {
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

