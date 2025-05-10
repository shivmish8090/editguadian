package modules


func init(){
    Register(handlers.NewCommand("stats", stats))
}
func stats(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveUser.Id != config.OwnerId {
	    
	    return Continue
	}
	var text
	if chats, err := database.GetServedChats(); err !=nil {
	    return err
	} else {
	    text = fmt.Sprintf("Total Chats: %d\n", len(chats))
	    
	}
	
	if users, err := database.GetServedChats(); err !=nil {
	    return err
	} else {
	    text = fmt.Sprintf("Total Users: %d\n", len(users))
	    
	}
}