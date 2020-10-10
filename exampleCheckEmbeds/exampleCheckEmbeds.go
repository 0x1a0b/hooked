package exampleCheckEmbeds

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
    secret := config.GetConf().Hooks.ExampleSenderEmbeds
    i.sender = discordSender.New(secret)

    log.Debugf("setup ended")
	return
}

const (
	InitialMessage = "init"
	FireMessage = "yess hello there embeds"
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
		Embeds: []discordSender.Embed{
			{
				Title: "bla bla bla",
				Url: "https://keeb.io/",
				Color: 14521290,
                Author: discordSender.Author{
					Name: "dem bot",
					IconUrl: "http://placekitten.com/g/200/300",
					Url: "http://placekitten.com",
				},
				Thumbnail: discordSender.Thumbnail{
					Url: "https://www.gdargaud.net/Antarctica/Life/RacingPenguin.jpg",
				},
				Fields: []discordSender.Field{
					{
						Name: "hello",
						Value: "there",
						Inline: true,
					},
				},
				Footer: discordSender.Footer{
					Text: "Footer",
				},
			},
		},
	}

	return
}
