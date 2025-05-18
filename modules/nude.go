package modules

import (
	"fmt"
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"main/config"
)

func init() {
	Register(handlers.NewMessage(func(m *gotgbot.Message) bool { return m.Photo != nil }, DeleteNudePhoto))
}

func DeleteNudePhoto(b *gotgbot.Bot, ctx *ext.Context) error {
	if !slices.Contains(config.OwnerId, ctx.EffectiveUser.Id) {
		return Continue
	}

	photos := ctx.EffectiveMessage.Photo

	var msg string
	for _, p := range photos {

		file, err := b.GetFile(p.FileId, nil)
		if err != nil {
			return err
		}

		msg += fmt.Sprintf("%s\n", file.URL(b, nil))

	}

	ctx.EffectiveMessage.Reply(b, msg, nil)

	return Continue
}
