package main

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/nozzlium/belimang/internal/client"
	"github.com/nozzlium/belimang/internal/config"
)

func main() {
	setupApp()
}

func setupApp() error {
	var cfg config.Config
	opts := env.Options{
		TagName: "json",
	}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		log.Fatalf("%+v\n", err)
		return err
	}

	_, err := client.InitDB(cfg.DB)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
