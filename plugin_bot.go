package brainwave_yagpdb

import (
	"fmt"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/common/templates"
)

var logger = common.GetPluginLogger(&Plugin{})

const (
	KeywordShort = "bw"
)

func RegisterPlugin() {
	if err := common.GORM.AutoMigrate(&Game{}).Error; err != nil {
		logger.WithError(err).Fatal("brainwave failed to initialize DB")
		return
	}

	p := &Plugin{}
	common.RegisterPlugin(p)
}

type Plugin struct {
}

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
		Aliases:      []string{KeywordShort},
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

func renderResponse(parsed *dcmd.Data, out interface{}) (interface{}, error) {
	outStr, ok := out.(string)
	if !ok {
		return out, nil
	}

	tmplCtx := templates.NewContext(parsed.GS, parsed.CS, nil)
	tmplCtx.Data["Keyword"] = KeywordShort
	return tmplCtx.Execute(outStr)
}

func runAction(parsed *dcmd.Data) (interface{}, error) {
	action := parsed.Args[0].Value.(Action)
	arg1 := parsed.Args[1].Str()

	game, err := loadGameFromDB(parsed.GS.ID, parsed.CS.ID)
	if err != nil {
		return "Life is cruel, and your game has been lost. Please start a new one.", nil
	}

	var out interface{}
	switch action {
	case Lead:
		out, err = game.runLead(parsed.Msg.Author, arg1)
	case Start:
		out, err = game.runStart()
	case Touch:
		out, err = game.runTouch(arg1)
	default:
		out, err = fmt.Sprintf("I don't know '%s'. Care to try one of my many other fine commands?", action), nil
	}
	if err != nil {
		return out, err
	}

	return renderResponse(parsed, out)
}
