package modules

import (
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
)

func init() {
	Register(handlers.NewCallback(callbackquery.Equal("help"), helpCB))
}

func helpCB(b *gotgbot.Bot, ctx *ext.Context) error {
btn := &buttons.Button{RowWidth: 2}
for name, mod := range ModulesHelp {
btn.Add(name, mod.Callback)

}
btn.Row("⬅️ Back", "start_callback")

	helpText := `📚 <b>Bot Command Help</b>

Here you'll find details for all available plugins and features.

👇 Tap the buttons below to view help for each module:`

	_, _, err := ctx.CallbackQuery.Message.EditCaption(b, &gotgbot.EditMessageCaptionOpts{
		Caption:     helpText,
		ReplyMarkup: btn.Build(),
		ParseMode:   "HTML",
	})
	if err != nil {
		log.Println("Failed to edit caption:", err)
	}
	return nil
}

func helpModuleCB(b *gotgbot.Bot, ctx *ext.Context) error {
	cbData := ctx.CallbackQuery.Data

	var helpText string
	for _, module := range ModulesHelp {
		if module.Callback == cbData {
			helpText = module.Help
			break
		}
	}

	if helpText == "" {
		helpText = "❌ No help found for this module."
	}

	_, _, err := ctx.CallbackQuery.Message.EditCaption(b, &gotgbot.EditMessageCaptionOpts{
		Caption:   helpText,
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{Text: "⬅️ Back", CallbackData: "help"},
				},
			},
		},
	})
	return err
}
