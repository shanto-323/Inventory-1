package main

import (
	"log"
	"time"

	"inventory/internal"
	"inventory/internal/storage"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	IpAddr string `envconfig:"IP_ADDR"`
	Dsn    string `envconfig:"DSN"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	var (
		repository storage.Repository
		err        error
	)

	retry.ForeverSleep(
		2*time.Second,
		func(_ int) error {
			repository, err = storage.NewRepository([]string{cfg.Dsn})
			if err != nil {
				log.Println(err)
				return err
			}
			return nil
		},
	)

	service := storage.NewService(repository)
	server := internal.NewServer(cfg.IpAddr, service)
	log.Fatal(server.Start())
}
