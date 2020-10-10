package main

import (
	"github.com/0x1a0b/hooked/config"
	"github.com/0x1a0b/hooked/exampleCheckEmbeds"
	"github.com/0x1a0b/hooked/exampleCheckSimple"
	"github.com/0x1a0b/hooked/geekHackRss"
	"github.com/0x1a0b/hooked/steam_economy"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"time"
)

func init() {
	log.SetLevel(config.GetLogLevel())
	log.SetReportCaller(true)
	config.Read()
	exampleSimple = exampleCheckSimple.Setup()
	exampleEmbeds = exampleCheckEmbeds.Setup()
	geekHack = geekHackRss.Setup()

	time.Sleep(1 * time.Second)
	log.Debugf("initialized")
}

var (
	geekHack *geekHackRss.Instance
	exampleEmbeds *exampleCheckEmbeds.Instance
	exampleSimple *exampleCheckSimple.Instance
)

func main() {
	// https://godoc.org/github.com/robfig/cron
	c := cron.New()
	c.AddFunc("@every 10m", func() { go geekHack.Run() })
	c.AddFunc("@every 10m", func() { go steam_economy.UpdateShop() })
	if log.GetLevel() == log.TraceLevel {
		c.AddFunc("@every 1m", func() { go exampleSimple.Run() })
		c.AddFunc("@every 1m", func() { go exampleEmbeds.Run() })
	}
	log.Debugf("starting cron in foreground")
	c.Run()
}
