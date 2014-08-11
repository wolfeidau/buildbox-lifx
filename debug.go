// +build !release

package main

import (
	"github.com/juju/loggo"
)

func init() {
	loggo.GetLogger("").SetLogLevel(loggo.DEBUG)
}
