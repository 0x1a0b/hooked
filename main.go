package main

import (
	"github.com/0x1a0b/hooked/config"
	log "github.com/sirupsen/logrus"
	"time"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	config.Read()
	time.Sleep(3 * time.Second)
	log.Debugf("initialized")
}

func main() {
	log.Debugf("main")
}