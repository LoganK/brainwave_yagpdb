package brainwave_yagpdb

import (
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
)

var _ bot.BotInitHandler = (*Plugin)(nil)
var _ commands.CommandProvider = (*Plugin)(nil)
var _ common.Plugin = (*Plugin)(nil)

