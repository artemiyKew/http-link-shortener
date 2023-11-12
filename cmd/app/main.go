package main

import (
	"github.com/artemiyKew/http-link-shortener/internal/app"
	"github.com/sirupsen/logrus"
)

const configPath = "config/config.yaml"

func main() {
	if err := app.Run(configPath); err != nil {
		logrus.Fatal(err)
	}
}
