package modules

import (
	"fmt"
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"main/config"
	"main/utils"
)

func init() {
	Register(handlers.NewMessage(func(m *gotgbot.Message) bool { return m.Photo != nil }, DeleteNudePhoto))
}

func DeleteNudePhoto(b *gotgbot.Bot, ctx *ext.Context) error {
	m := ctx.EffectiveMessage

	if !slices.Contains(config.OwnerId, ctx.EffectiveUser.Id) {
		return Continue
	}

	photo := m.Photo[len(m.Photo)-1]

	file, err := b.GetFile(photo.FileId, nil)
	if err != nil {
		return err
	}

	path, e := utils.DownloadFile(file.URL(b, nil))

	if e != nil {
		return e
	}

	m.Reply(b, fmt.Sprintf("Successfully Downloaded in %s", path), nil)

	return Continue
}
