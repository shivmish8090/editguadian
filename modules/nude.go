package modules

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func init() {
	Register(handlers.NewMessage(func(m *gotgbot.Message) bool { return m.EditDate != 0 }, DeleteNudePhoto))
}

func DeleteNudePhoto(b *gotgbot.Bot, ctx *ext.Context) {}
