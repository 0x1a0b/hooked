package exampleCheckSimple

import (
	"github.com/0x1a0b/hooked/config"
	"github.com/0x1a0b/hooked/discordSender"
	log "github.com/sirupsen/logrus"
)

type Instance struct {
	Message string
	sender *discordSender.Sender
}

func Setup() (i *Instance) {
	log.Debugf("setup started")

	i = &Instance{}
    i.Message = InitialMessage
    secret := config.GetConf().Hooks.ExampleSender
    i.sender = discordSender.New(secret)

    log.Debugf("setup ended")
	return
}

const (
	InitialMessage = "init"
	FireMessage = "yess hello there"
)

func (i *Instance) Run() () {
	log.Debugf("started to run")

	if i.Message == InitialMessage {
		log.Debugf("there is something todo")
		i.Message = FireMessage
		i.Fire()
	} else {
		log.Debugf("nothing to do")
	}

	return
}


func (i *Instance) Fire() () {

	hook := i.SetHook()
	_ = i.sender.Send(hook)

	log.WithField("hook", hook).Debugf("fired hook")
	return

}

func (i *Instance) SetHook() (h discordSender.Hook) {
	h = discordSender.Hook{
		Content: i.Message,
		AvatarUrl: "https://i.kym-cdn.com/photos/images/newsfeed/000/925/493/19f.jpg",
		Username: "Kappa",
	}

	return
}
