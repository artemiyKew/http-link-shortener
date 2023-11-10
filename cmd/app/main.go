package main

import (
	"context"

	"github.com/artemiyKew/http-link-shortener/internal/app"
	"github.com/sirupsen/logrus"
)

const configPath = "config/config.yaml"

func main() {
	if err := app.Run(context.Background(), configPath); err != nil {
		logrus.Fatal(err)
	}
}
