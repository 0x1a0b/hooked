package main

import (
	"github.com/0x1a0b/hooked/config"
	"github.com/0x1a0b/hooked/steam_economy"
	"github.com/kz/discordrus"
	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func init() {
	log.SetLevel(logrus.DebugLevel)
	log.SetReportCaller(true)
	config.Read()
	time.Sleep(1 * time.Second)
	LastItemLink = ""
	discordlog.SetFormatter(&logrus.TextFormatter{})
	discordlog.SetOutput(os.Stderr)
	discordlog.SetLevel(logrus.TraceLevel)
	discordlog.AddHook(discordrus.NewHook(
		config.GetConf().Discord.Url,
		logrus.TraceLevel,
		&discordrus.Opts{
			Username: "geekhack hook",
			Author:             "",
			DisableTimestamp:   false,
			TimestampFormat:    "Jan 2 15:04:05.00000 MST",
			TimestampLocale:    nil,
			EnableCustomColors: true,
			CustomLevelColors: &discordrus.LevelColors{
				Trace: 3092790,
				Debug: 10170623,
				Info:  3581519,
				Warn:  14327864,
				Error: 13631488,
				Panic: 13631488,
				Fatal: 13631488,
			},
			DisableInlineFields: false,
		},
	))
	log.Debugf("initialized")
}

const (
	// https://wiki.simplemachines.org/smf/XML_feeds
	GeekHackUrl = "https://geekhack.org/index.php?action=.xml;type=rss;limit=8;sa=news;board=70.0"
)

var (
	LastItemLink string
    log = logrus.New()
    discordlog = logrus.New()
)

func main() {
	// https://godoc.org/github.com/robfig/cron
	c := cron.New()
	c.AddFunc("@every 3m", func() { go update() })
	c.AddFunc("@every 3m", func() { go steam_economy.UpdateShop() })
	c.Run()
}

func update() () {
	// https://godoc.org/github.com/mmcdole/gofeed#Item
	log.Debugf("starting update at %v", time.Now())
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(GeekHackUrl)
	if LastItemLink == "" {
		LastItemLink = feed.Items[0].Link
		log.Printf("freshly started app, assuming %v as top item", LastItemLink)
		discordlog.Trace("resuming operations after restart")
	} else if feed.Items[0].Link == LastItemLink {
		log.Printf("top item is still %v, doing noting", LastItemLink)
	} else {
		for index, item := range feed.Items {
			if item.Link == LastItemLink {
				log.Printf("items delta finished with last seen item %v", LastItemLink)
				LastItemLink = item.Link
				log.Printf("new top item is now %v", LastItemLink)
				break
			} else {
				log.Printf("new item %v at index %v", item.Link, index)
				discordlog.WithFields(logrus.Fields{"Published": item.Published, "Link": item.Link}).Warnf("New geekhack post: %v \n %v", item.Title, item.Description)
			}
		}
	}
	log.Debugf("ending update at %v", time.Now())
}