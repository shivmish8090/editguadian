package modules 


func init(){

Register(handlers.NewMessage(func(m *gotgbot.Message) bool { return m.EditDate != 0}, DeleteNudePhoto))

}


func DeleteNudePhoto(b *gotgbot.Bot, ctx *ext.Context) {}