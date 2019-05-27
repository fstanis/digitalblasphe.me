package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/fstanis/digitalblasphe.me/internal/changer"
	"github.com/fstanis/digitalblasphe.me/internal/config"
	"github.com/fstanis/digitalblasphe.me/pkg/digitalblasphemy"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	digitalblasphemy.LoadCache()

	conf, err := config.Load()
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

	digitalblasphemy.SaveCache()
}
