package main

import (
	"github.com/0x1a0b/hooked/config"
	"github.com/0x1a0b/hooked/exampleCheck"
	"github.com/0x1a0b/hooked/geekhack_rss"
	"github.com/0x1a0b/hooked/steam_economy"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"time"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	config.Read()
	example = exampleCheck.Setup()
	time.Sleep(1 * time.Second)
	log.Debugf("initialized")
}

var (
	example *exampleCheck.Instance
)

func main() {
	// https://godoc.org/github.com/robfig/cron
	c := cron.New()
	c.AddFunc("@every 10m", func() { go geekhack_rss.Update() })
	c.AddFunc("@every 10m", func() { go steam_economy.UpdateShop() })
	c.AddFunc("@every 1m", func() { go example.Run() })
	log.Debugf("starting cron in foreground")
	c.Run()
}
