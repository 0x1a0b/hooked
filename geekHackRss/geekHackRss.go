package geekHackRss

import (
	"github.com/0x1a0b/hooked/config"
	"github.com/0x1a0b/hooked/discordSender"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

const (
	// https://wiki.simplemachines.org/smf/XML_feeds
	GeekHackUrl = "https://geekhack.org/index.php?action=.xml;type=rss;limit=8;sa=news;board=70.0"
)

type Instance struct {
	LastItem string
	sender *discordSender.Sender
	parser *gofeed.Parser
	logger *logrus.Logger
}

func Setup() (i *Instance) {

	i = &Instance{
		LastItem: "",
	}

	secret := config.GetConf().Hooks.GeekHackNews
	i.sender = discordSender.New(secret)

	i.logger = logrus.New()
	i.logger.SetLevel(config.GetLogLevel())
	i.logger.SetReportCaller(true)

	i.parser = gofeed.NewParser()

	return

}

func (i *Instance) Run() () {
	i.logger.Debugf("starting run, lastitem is %v", i.LastItem)

	feed, err := i.parser.ParseURL(GeekHackUrl)
	if err != nil {
		i.logger.Errorf("error parsing rss: %v", err)
		i.LastItem = ""
		i.logger.Errorf("resetting Lastitem, now - %v - ", i.LastItem)
		return
	}

	if i.LastItem == "" {
		i.logger.Debugf("LastItem Empty")
		i.LastItem = feed.Items[0].Link
		i.logger.Debugf("Resuming with empty HEAD or after error, setting %v", i.LastItem)
		i.Fire(feed.Items[0])
		return
	} else if i.LastItem == feed.Items[0].Link {
		i.logger.Debugf("same HEAD, doing nothing")
	} else {
		i.logger.Debugf("HEAD and LastItem are different. HEAD is %v while the last item seen was %v", feed.Items[0].Link, i.LastItem)
		for index, item := range feed.Items {
			if item.Link == i.LastItem {
				i.LastItem = item.Link
				i.logger.Debugf("delta ended with %v at %v", item.Link, index)
				break
			} else {
				i.logger.Debugf("firing hook for %v", item.Link)
				i.Fire(item)
			}
		}
		i.logger.Debugf("new lastitem is %v", i.LastItem)
	}

	i.logger.Debugf("ending run")
	return
}

func (i *Instance) Fire(item *gofeed.Item) () {
	hook := i.SetHook(item)
	if err := i.sender.Send(hook); err != nil {
		i.logger.WithField("hook", hook).Errorf("error firing hook: %v", err)
	} else {
		i.logger.WithField("hook", hook).Debugf("fired hook")
	}
	return
}

func (i *Instance) SetHook(item *gofeed.Item) (h discordSender.Hook) {

	h = discordSender.Hook{
		Content: "New Geekhack post: " + item.Title,
		AvatarUrl: "https://wiki.geekhack.org/images/thumb/5/51/Cherry_MX_White_Plate_Mount_Switch.jpg/300px-Cherry_MX_White_Plate_Mount_Switch.jpg",
		Username: "Geekhack",
		Embeds: []discordSender.Embed{
		{
			Title: item.Title,
			Description: item.Content,
            Color: 3581519,
            Url: item.Link,
            Thumbnail: discordSender.Thumbnail{
				Url: "https://cdn.mos.cms.futurecdn.net/qhWUafmdRWYFhFSyAiAFPh-970-80.jpg",
			},
		},
		},
	}

	return
}
