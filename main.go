package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/juju/loggo"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {

	// configure the logger
	log := loggo.GetLogger("buildbox-lifx")

	log.Infof("starting service version: ", Version)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	s := <-c
	fmt.Println("Got signal:", s)

	return 0
}
