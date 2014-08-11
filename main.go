package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/juju/loggo"
	"github.com/wolfeidau/lifx"
)

type Projects struct {
	XMLName xml.Name `xml:"Projects"`
	Project []Project
}

type Project struct {
	XMLName         xml.Name `xml:"Project"`
	Name            string   `xml:"name,attr"`
	Activity        string   `xml:"activity,attr"`
	LastBuildStatus string   `xml:"lastBuildStatus,attr"`
	LastBuildLabel  string   `xml:"lastBuildLabel,attr"`
	LastBuildTime   string   `xml:"lastBuildTime,attr"`
	WebURL          string   `xml:"webUrl,attr"`
}

func main() {
	os.Exit(realMain())
}

func realMain() int {

	// configure the logger
	log := loggo.GetLogger("buildbox-lifx")

	// flags
	app := cli.NewApp()

	app.Name = "buildbox-lifx"
	app.Usage = "Monitors buildbox and changes lifx bulbs to reflect success or failure"
	app.Version = Version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "apikey",
			Value: "",
			Usage: "buildbox api key",
		},
		cli.StringFlag{
			Name:  "branch",
			Value: "master",
			Usage: "branch to filter builds",
		},
		cli.StringFlag{
			Name:  "bulb",
			Value: "build",
			Usage: "the label of the bulb you want to control",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug logging",
		},
	}

	app.Action = func(c *cli.Context) {

		if c.String("apikey") == "" {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		if c.Bool("debug") {
			loggo.GetLogger("").SetLogLevel(loggo.DEBUG)
		}

		log.Infof("starting service version: %s", Version)

		client := lifx.NewClient()

		client.StartDiscovery()

		for {
			url := fmt.Sprintf("https://cc.buildbox.io/ninja-blocks-inc.xml?api_key=%s&branch=%s", c.String("apikey"), c.String("branch"))

			resp, err := http.Get(url)

			if err != nil {
				log.Errorf("Woops issue contacting buildbox: %s", err)
			}

			if resp.StatusCode == http.StatusOK {
				var p Projects

				err := xml.NewDecoder(resp.Body).Decode(&p)

				if err != nil {
					log.Errorf("Woops issue decoding response: %s", err)
				}

				log.Infof("Found %d Projects", len(p.Project))

				sc := getStatusCount(&p)

				log.Infof("Status Counts %+v", sc)

				for _, bulb := range client.GetBulbs() {

					log.Debugf("checking bulb label: %s", bulb.GetLabel())

					if bulb.GetLabel() == c.String("bulb") {

						// ensure globe is on
						client.LightOn(bulb)

						if sc.hasFailures() {
							// red
							log.Debugf("setting bulb : %s to red", bulb.GetLabel())
							client.LightColour(bulb, 65170, 55299, 65535, 3500, 1000)
						} else {
							log.Debugf("setting bulb : %s to green", bulb.GetLabel())
							// green
							client.LightColour(bulb, 20570, 65535, 33095, 3500, 1000)
						}

					}
				}

			}

			time.Sleep(30 * time.Second)
		}

	}

	// go forth and poll
	app.Run(os.Args)

	return 0
}

type statusCount struct {
	Success, Failure, Unknown, Exception int
}

func (sc *statusCount) hasFailures() bool {
	return (sc.Failure > 0)
}

func getStatusCount(p *Projects) *statusCount {
	s := &statusCount{}

	for _, proj := range p.Project {

		switch proj.LastBuildStatus {
		case "Success":
			s.Success++

		case "Failure":
			s.Failure++

		case "Unknown":
			s.Unknown++

		case "Exception":
			s.Exception++
		}

	}

	return s
}
