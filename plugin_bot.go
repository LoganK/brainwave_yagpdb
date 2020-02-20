package brainwave_yagpdb

import (
	"encoding/json"
	"fmt"

	"github.com/jbowens/codenames"
	"github.com/jinzhu/gorm"
	"github.com/jonas747/dcmd"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
)

var logger = common.GetPluginLogger(&Plugin{})

type brainwaveGame struct {
	gorm.Model
	GuildID   int64           `gorm:"primary_key;auto_increment:false"`
	ChannelID string          `gorm:"primary_key;auto_increment:false"`
	Game      codenames.Game  `gorm:"-"`
	GameSave  json.RawMessage `sql:"type:json"`
}

func (g *brainwaveGame) BeforeSave() error {
	gameObj, err := json.Marshal(g.Game)
	if err != nil {
		logger.WithError(err).Errorf("failed to save game")
		return err
	}

	g.GameSave = gameObj
	return nil
}

func (g *brainwaveGame) AfterFind() {
	if err := json.Unmarshal(g.GameSave, &g.Game); err != nil {
		logger.WithError(err).Errorf("failed to restore game")
		return
	}

	g.GameSave = json.RawMessage{}
}

func RegisterPlugin() {
	if err := common.GORM.AutoMigrate(&brainwaveGame{}).Error; err != nil {
		logger.WithError(err).Fatal("brainwave failed to initialize DB")
	}

	p := &Plugin{}
	common.RegisterPlugin(p)
}

type Plugin struct{}

func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "Brainwave",
		SysName:  "brainwave",
		Category: common.PluginCategoryMisc,
	}
}

func (p *Plugin) AddCommands() {
	commands.AddRootCommands(&commands.YAGCommand{
		CmdCategory:  commands.CategoryFun,
		Name:         "brainwave",
		Description:  "The main command to play Brainwave, a Codenames-like game.",
		Aliases:      []string{"bw"},
		RequiredArgs: 1,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Action", Type: ActionArg{}},
			&dcmd.ArgDef{Name: "TeamOrWord", Type: dcmd.String},
		},
		RunFunc: runAction,
	})
}

func (p *Plugin) BotInit() {
}

func runAction(parsed *dcmd.Data) (interface{}, error) {
	action := parsed.Args[0].Value.(Action)
	arg1 := parsed.Args[1].Str()
	switch action {
	case Lead:
		return runLead(arg1)
	case Touch:
		return runTouch(arg1)
	default:
		return fmt.Sprintf("I don't know '%s'. Care to try one of my many other fine commands?", action), nil
	}
}

func runLead(team string) (interface{}, error) {
	return nil, nil
}

func runTouch(word string) (interface{}, error) {
	return nil, nil
}
