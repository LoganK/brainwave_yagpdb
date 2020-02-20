package brainwave_yagpdb

import (
        "fmt"
        "strings"

	"github.com/jonas747/dcmd"
)

type Action string

const (
        Invalid = "XXinvalidXX"
        Lead = "lead"
        Start = "start"
        Teams = "teams"
        Touch = "touch"
)

func fromString(s string) (Action, bool) {
        s = strings.ToLower(s)
        switch s {
        case "lead":
                return Lead, true
        case "start":
                return Start, true
        case "teams":
                return Teams, true
        case "touch":
                return Touch, true
        default:
                return Invalid, false
        }
}

type ActionArg struct{}

func (_ ActionArg) Matches(def *dcmd.ArgDef, part string) bool {
        _, ok := fromString(part)
        return ok
}

func (_ ActionArg) Parse(def *dcmd.ArgDef, part string, data *dcmd.Data) (val interface{}, err error) {
        a, ok := fromString(part)
        if !ok {
                return Action(Invalid), fmt.Errorf("invalid string for Action: %s", part)
        }
        return a, nil
}

func (_ ActionArg) HelpName() string {
        return "Action"
}

