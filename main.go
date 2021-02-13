package main

import (
	"github.com/Prajithp/gchat-notifier/app"
	"github.com/Prajithp/gchat-notifier/config"
	"log"
)

func main() {
	c, err := config.ReadConfig("alerts.yaml")
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	app := &app.App{}
	app.Initialize(c)
	app.Run(":3000")
}
