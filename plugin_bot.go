package brainwave_yagpdb

import (
	"fmt"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
)

var logger = common.GetPluginLogger(&Plugin{})

func RegisterPlugin() {
	if err := common.GORM.AutoMigrate(&Game{}).Error; err != nil {
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

        game, err := loadGameFromDB(parsed.GS.ID, parsed.CS.ID)
        if err != nil {
                return "Life is cruel, and your game has been lost. Please start a new one.", nil
        }

	switch action {
	case Lead:
		return game.runLead(parsed.Msg.Author, arg1)
        case Start:
                return game.runStart()
	case Touch:
		return game.runTouch(arg1)
	default:
		return fmt.Sprintf("I don't know '%s'. Care to try one of my many other fine commands?", action), nil
	}
}

