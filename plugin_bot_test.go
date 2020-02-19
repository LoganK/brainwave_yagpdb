package brainwave_yagpdb

import (
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/commands"
)


var _ bot.BotInitHandler = (*Plugin)(nil)
var _ commands.CommandProvider = (*Plugin)(nil)

