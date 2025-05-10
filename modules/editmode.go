package modules

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"main/config/helpers"
)

func init() {
	Register(handlers.NewCommand("editmode", EditMode))
	AddHelp("✍️ Edit Mode", "editmode", `<b>Command:</b>
<blockquote>
<b>/editmode</b> – Show current settings  
<b>/editmode --set-mode=&lt;off|user|admin&gt;</b> – Set deletion mode  
<b>/editmode --set-duration=&lt;number&gt;</b> – Set deletion duration (in minutes)
</blockquote>
<b>Description:</b>  
Controls how the bot handles <b>edited messages</b> by deleting them based on mode and duration.

<b>Modes:</b>  
• <b>/editmode</b> <code>--set-mode=off</code> – No deletion  
• <b>/editmode</b> <code>--set-mode=user</code> – Delete edits from users (default)  
• <b>/editmode</b> <code>--set-mode=admin</code> – Delete edits from users & admins <i>(owner only)</i>

<b>Duration (via --set-duration):</b>  
• <b>0</b> – Deletes immediately <i>(default)</i> and warns users  
• <b>1-10</b> – Deletes if edited after the set number of minutes (no warning)  
• <b>&gt;10</b> – Disables deletion of edits entirely

<b>Examples:</b>  
<blockquote>
<code>/editmode --set-mode=user</code>  
<code>/editmode --set-duration=5</code>  
<code>/editmode --set-mode=admin --set-duration=2</code>
</blockquote>`, nil)
}

func EditMode(b *gotgbot.Bot, ctx *ext.Context) error {
	Message := ctx.EffectiveMessage
	if Message.SenderChat != nil {
		Message.Reply(
			b,
			"You are anonymous Admin you can't use this command.",
			nil,
		)
		return Continue
	}

	if ctx.EffectiveChat.Type != "supergroup" {
		Message.Reply(
			b,
			"This command is meant to be used in supergroups, not in private messages!",
			nil,
		)
		return Continue
	}

	args := ctx.Args()

	if len(args) < 2 {
		ctx.EffectiveMessage.Reply(b,
			fmt.Sprintf("Usage: <code>/editmode &lt;off|admin|user&gt;</code>\n<b>For more help, check out:</b> <a href=\"%s\">Edit Mode Help</a>",
				fmt.Sprintf("https://t.me/%s?start=editmode", b.User.Username)),
			&gotgbot.SendMessageOpts{ParseMode: "HTML"})
		ctx.EffectiveMessage.Delete(b, nil)
		return Continue
	}
	admins, err := helpers.GetAdmins(b, ctx.EffectiveChat.Id)
	if err != nil {
		return err
	}

	if !helpers.Contains(admins, ctx.EffectiveUser.Id) {
		b.SendMessage(ctx.EffectiveChat.Id, "Access denied: Only group admins can use this command.\n\nIf you are an admin, please use /reload to refresh the admin list.", nil)
		return Continue
	}
	keys := []string{"set-mode", "set-duration"}
	_, res := helpers.ParseFlags(keys, Message.Text)

	if res["set-mode"] != "off" && res["set-mode"] != "user" && res["set-mode"] != "admin" {
		ctx.EffectiveMessage.Reply(b,
			fmt.Sprintf("Usage: <code>/editmode &lt;off|admin|user&gt;</code>\n<b>For more help, check out:</b> <a href=\"%s\">Edit Mode Help</a>",
				fmt.Sprintf("https://t.me/%s?start=editmode", b.User.Username)),
			&gotgbot.SendMessageOpts{ParseMode: "HTML"})
		ctx.EffectiveMessage.Delete(b, nil)
		return Continue
	}

	if res["set-mode"] == "admin" {
		var ownerID int64
		ownerID, err = helpers.GetOwner(b, ctx.EffectiveChat.Id)
		if err != nil {
			return err
		}

		if ctx.EffectiveUser.Id != ownerID {
			ctx.EffectiveMessage.Reply(b, "Only for owner", nil)
			return Continue
		}

	}
	ctx.EffectiveMessage.Reply(b, "This command will be available soon..", nil)
	ctx.EffectiveMessage.Delete(b, nil)
	return Continue
}
