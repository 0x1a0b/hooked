package exampleCheck

import "github.com/0x1a0b/hooked/discordSender"

type Instance struct {
	Message string
	sender *discordSender.Sender
}

func Setup() (i *Instance) {
    i.Message = InitialMessage
    i.sender = discordSender.New("test")
	return
}

const (
	InitialMessage = "init"
	FireMessage = "yess"
)

func (i *Instance) Run() () {

	if i.Message == InitialMessage {
		i.Message = FireMessage
		i.Fire()
	}

	return
}


func (i *Instance) Fire() () {

	hook := i.SetHook()
	_ = i.sender.Send(hook)
	return

}

func (i *Instance) SetHook() (h discordSender.Hook) {

	return
}