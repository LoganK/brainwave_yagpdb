package brainwave_yagpdb

import (
	"fmt"
	//"github.com/jinzhu/gorm"
	"github.com/jonas747/dcmd"
	"github.com/jonas747/yagpdb/commands"
)

//var logger = common.GetPluginLogger(&Plugin{})

type Plugin struct{}

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
