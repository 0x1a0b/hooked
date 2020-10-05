package main

import (
	"github.com/0x1a0b/hooked/config"
	"github.com/mmcdole/gofeed"
	cron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"time"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	config.Read()
	time.Sleep(1 * time.Second)
	LastItemLink = ""
	log.Debugf("initialized")
}

const (
	// https://wiki.simplemachines.org/smf/XML_feeds
	// GeekHackUrl = "https://geekhack.org/index.php?action=.xml;type=rss;limit=100"
	GeekHackUrl = "https://geekhack.org/index.php?action=.xml;type=rss;limit=100;sa=news;board=70.0"
)

var (
	LastItemLink string
)

func main() {
	c := cron.New()
	c.AddFunc("@every 3m", func() { go update() })
	c.Run()
}

func update() () {
	log.Debugf("starting update at %v", time.Now())
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(GeekHackUrl)
	if LastItemLink == "" {
		LastItemLink = feed.Items[0].Link
		log.Printf("freshly started app, assuming %v as top item", LastItemLink)
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
			}
		}
	}
	log.Debugf("ending update at %v", time.Now())
}