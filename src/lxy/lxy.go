package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

func main() {
	initLogging()
	Control()
}

func initLogging() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
}
