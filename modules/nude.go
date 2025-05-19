package modules

import (
	"fmt"
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/koyachi/go-nude"

	"main/config"
	"main/utils"
)

func init() {
	Register(handlers.NewMessage(func(m *gotgbot.Message) bool { return m.Photo != nil || m.Sticker != nil }, DeleteNudePhoto))
}

func DeleteNudePhoto(b *gotgbot.Bot, ctx *ext.Context) error {
	m := ctx.EffectiveMessage

	if !slices.Contains(config.OwnerId, ctx.EffectiveUser.Id) {
		return Continue
	}

	var fileid string
	if m.Photo != nil {
		fileid = m.Photo[len(m.Photo)-1].FileId
	} else {
		fileid = m.Sticker.FileId
	}
	file, err := b.GetFile(fileid, nil)
	if err != nil {
		return err
	}

	var path string
	path, err = utils.DownloadFile(file.URL(b, nil))
	if err != nil {
		return err
	}

	if m.Sticker != nil && !m.Sticker.IsVido && !m.Sticker.IsAnimated {
		err = utils.Webp2Png(path)
		if err != nil {
			return err
		}
	} else if m.Sticker.IsVideo && !m.Sticker.IsAnimated {
		return nil
	}

	var isNude bool

	isNude, err = nude.IsNude(path)
	if err != nil {
		return err
	}
	m.Reply(b, fmt.Sprintf("Your image contains nudity: %t", isNude), nil)

	return Continue
}
