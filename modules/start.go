package modules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"

	"main/config"
	"main/config/buttons"
	"main/database"
)

func init() {
	Register(handlers.NewCommand("start", start))

	Register(handlers.NewCallback(callbackquery.Equal("start_callback"), start))
}

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	isCallback := ctx.CallbackQuery != nil
	chatType := ctx.EffectiveChat.Type

	if chatType == "private" {
		if !isCallback {
			ctx.EffectiveMessage.Delete(b, nil)
		}

		args := ctx.Args()
		if len(args) >= 2 {
			modName := args[1]
			if strings.HasPrefix(modName, "info_") {
				userIDStr := strings.TrimPrefix(modName, "info_")
				userID, err := strconv.ParseInt(userIDStr, 10, 64)
				if err != nil {
					return err
				}

				userInfo, err := b.GetChat(userID, nil)
				if err != nil {
					return err
				}

				fullName := strings.TrimSpace(userInfo.FirstName + " " + userInfo.LastName)
				info := fmt.Sprintf(`
Name: %s
Id: %d
Link: <a href="tg://user?id=%d">Link 1</a> <a href="tg://openmessage?user_id=%d">Link 2</a>
`, fullName, userInfo.Id, userInfo.Id, userInfo.Id)

				_, err = b.SendMessage(ctx.EffectiveChat.Id, info, &gotgbot.SendMessageOpts{
					ParseMode: "HTML",
				})
				return orCont(err)
			}

			modHelp := GetHelp(modName)
			if modHelp != "" {
				_, err := b.SendMessage(ctx.EffectiveChat.Id, modHelp, &gotgbot.SendMessageOpts{
					ParseMode: "HTML",
				})
				if err != nil {
					return err
				}
				return Continue
			}
		}

		startImg := gotgbot.InputFileByURL(config.StartImage)
		userFullName := strings.TrimSpace(ctx.EffectiveUser.FirstName + " " + ctx.EffectiveUser.LastName)
		botName := strings.TrimSpace(b.User.FirstName)

		caption := fmt.Sprintf(
			`<b>🛡 Hello <a href="tg://user?id=%d">%s</a>!</b> 👋  
I'm <b><a href="tg://user?id=%d">%s</a></b>, your group’s security bot keeping chats clean and safe.

✏️ <b>Edited messages</b> are auto-deleted  
🖼️ <b>Media</b> is cleaned up instantly  
📜 <b>Long messages</b> (default 800+ chars) get removed — limit is <i>customizable</i>!

📣 Stay informed with instant alerts.  
✅ Add me now and I’ll start protecting your group!`,
			ctx.EffectiveUser.Id,
			userFullName,
			b.User.Id,
			botName,
		)

		var keyboard gotgbot.InlineKeyboardMarkup
		if ctx.EffectiveUser.Id == 7706682472 {
			keyboard = buttons.StartPanel(b)
		} else {
			keyboard = buttons.NormalStartPanel(b)
		}

		if isCallback {
			_, _, err := ctx.CallbackQuery.Message.EditCaption(b, &gotgbot.EditMessageCaptionOpts{
				Caption:     caption,
				ParseMode:   "HTML",
				ReplyMarkup: keyboard,
			})
			if err != nil {
				return err
			}
		} else {
			database.AddServedUser(ctx.EffectiveUser.Id)

			if config.StartImageGif != "" {
				_, err := b.SendAnimation(ctx.EffectiveChat.Id, gotgbot.InputFileByURL(config.StartImageGif), &gotgbot.SendAnimationOpts{
					Caption:     caption,
					ParseMode:   "HTML",
					ReplyMarkup: keyboard,
				})
				if err != nil {
					return fmt.Errorf("failed to send Gif: %w", err)
				}
			} else {
				_, err := b.SendPhoto(ctx.EffectiveChat.Id, startImg, &gotgbot.SendPhotoOpts{
					Caption:     caption,
					ParseMode:   "HTML",
					ReplyMarkup: keyboard,
				})
				if err != nil {
					return fmt.Errorf("failed to send photo: %w", err)
				}
			}

			if r := database.IsLoggerEnabled(); !r {
				return Continue
			}
			logMsg := fmt.Sprintf(
				`<a href="tg://user?id=%d">%s</a> has started the bot.

<b>User ID:</b> <code>%d</code>
<b>User Name:</b> %s`,
				ctx.EffectiveUser.Id,
				ctx.EffectiveUser.FirstName,
				ctx.EffectiveUser.Id,
				userFullName,
			)

			b.SendMessage(config.LoggerId, logMsg, &gotgbot.SendMessageOpts{
				ParseMode: "HTML",
			})
		}

	} else if chatType == "group" {
		if isCallback {
			return Continue
		}

		msg := `⚠️ Warning: I can't function in a basic group!

To use my features, please upgrade this group to a supergroup.

✅ How to upgrade:
1. Go to Group Settings.
2. Tap on "Chat History" and set it to "Visible".
3. Re-add me, and I'll be ready to help!`

		ctx.EffectiveMessage.Reply(b, msg, nil)
		ctx.EffectiveChat.Leave(b, nil)

	} else if chatType == "supergroup" {
		if isCallback {
			return Continue
		}

		database.AddServedChat(ctx.EffectiveChat.Id)

		r, errr := b.GetChatMember(ctx.EffectiveChat.Id, b.User.Id, nil)
		if err != nil {
			return errr
		}

		if member := r.MergeChatMember(); member.Status != "administrator" {

			ctx.EffectiveMessage.Reply(b, "⚠️ Please promote me to <b>Admin</b> so I can protect this group properly.", &gotgbot.SendMessageOpts{ParseMode: HTML})
			return Continue

		} else if !member.CanDeleteMessages {
			ctx.EffectiveMessage.Reply(b,
				`🧹 <b>I need the "Delete messages" permission</b> to remove spam or unwanted content.

Please enable it from the admin settings.`,
				&gotgbot.SendMessageOpts{ParseMode: HTML})
			return Continue
		}
		ctx.EffectiveMessage.Reply(b, "✅ I am active and ready to protect this supergroup!", nil)
	}

	return Continue
}
