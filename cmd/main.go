package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/fstanis/digitalblasphe.me/internal/changer"
	"github.com/fstanis/digitalblasphe.me/internal/config"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	conf, err := config.LoadConfig()
	if err != nil {
		conf, err = config.FromSurvey()
		if err != nil {
			log.Fatal(err)
		}
		conf.Save()
	}
	if err := changer.Apply(conf); err != nil {
		log.Fatal(err)
	}
}
