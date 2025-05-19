package modules

import (
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func init() {
	Register(handlers.NewMessage(func(m *gotgbot.Message) bool {
		if m.Entities == nil {
			return false
		}

		return entity == "url"
	}, DeleteLinkMessage))
}

func DeleteLinkMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	m := ctx.EffectiveMessage

	_, err := m.Delete(b, nil)

if err != nil {
	return err

}
b.SendMessage(
    m.Chat.Id,
    "⚠️ Direct URLs aren't allowed.\nUse format: <a href='https://t.me/durov'>this</a>",
    &gotgbot.SendMessageOpts{ParseMode: "HTML"},
)


return Continue
}
