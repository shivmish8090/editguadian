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

		return slices.ContainsFunc(m.Entities, func(entity gotgbot.MessageEntity) bool {
			return slices.Contains([]string{"text_link", "url"}, entity.Type)
		})
	}, DeleteLinkMessage))
}

func DeleteLinkMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	m := ctx.EffectiveMessage

	m.Delete(b, nil)
return nil
}
