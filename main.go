package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

const QRSize = 300

func main() {
	useOS := len(os.Args) > 1 && os.Args[1] == "live"

	a := &App{
		Port:   8080,
		Logger: logrus.StandardLogger(),
		Live:   useOS,
	}
	a.Logger.SetLevel(logrus.DebugLevel)

	a.Initialize()

	if err := a.StartServer(); err != nil {
		a.Logger.Fatal(err)
	}
}
