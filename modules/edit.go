package modules

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func DeleteEditedMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EditedMessage
	if message == nil || ctx.EffectiveChat.Type == "private" {
		return Continue
	}

	Chat, err := b.GetChat(ctx.EffectiveChat.Id, nil)
	if err != nil {
		return err
	}

	if message.SenderChat != nil {
		if message.Chat.Id == message.SenderChat.Id || Chat.LinkedChatId == message.SenderChat.Id {
			return Continue
		}
	}

	if _, err = ctx.EffectiveMessage.Delete(b, nil); err != nil {
		return orCont(err)
	}

	reason := "<b>🚫 Editing messages is prohibited in this chat.</b> Please refrain from modifying your messages to maintain the integrity of the conversation."

	switch {
	case message.Text != "":
		reason = "<b>🚫 Editing text is not allowed.</b> Please avoid changing messages once sent to keep conversations clear."

	case message.Caption != "":
		reason = "<b>✍️ Caption edits are restricted.</b> Changing them affects clarity and is not permitted."

	case message.Photo != nil:
		reason = "<b>📷 Photo edits are blocked.</b> Images must stay unchanged to preserve context."

	case message.Video != nil:
		reason = "<b>🎥 Video edits aren't allowed.</b> Videos must remain as originally shared."

	case message.Document != nil:
		reason = "<b>📄 Document edits are restricted.</b> Keep documents unchanged for reliability."

	case message.Audio != nil:
		reason = "<b>🎵 Audio edits aren't permitted.</b> Audio files must remain unaltered."

	case message.VideoNote != nil:
		reason = "<b>📹 Video note edits are not allowed.</b> They must stay as sent."

	case message.Voice != nil:
		reason = "<b>🎙️ Voice edits are restricted.</b> Voice messages should remain original."

	case message.Animation != nil:
		reason = "<b>🎞️ GIF edits are blocked.</b> Keep animations unchanged for context."

	case message.Sticker != nil:
		reason = "<b>🖼️ Sticker edits are not permitted.</b> Stickers must stay unaltered."
	}

	_, err = b.SendMessage(
		ctx.EffectiveChat.Id,
		reason,
		&gotgbot.SendMessageOpts{ParseMode: "HTML"},
	)
	if err != nil {
		return err
	}

	return Continue
}
