package modules

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"main/config"
	"main/config/helpers"
)

func init() {
	Register(handlers.NewCommand("ping", pingHandler))
}

func pingHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	uptime := time.Since(config.StartTime)
	uptimeStr := helpers.FormatUptime(uptime)
	ctx.EffectiveMessage.Delete(b, nil)

	_, err := ctx.EffectiveChat.SendMessage(b, "Bot has been running for: "+uptimeStr, nil)
	return orCont(err)
}
