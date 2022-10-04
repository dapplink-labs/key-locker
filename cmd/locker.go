package cmd

import (
	"github.com/savour-labs/key-locker/api"
	"os"

	"github.com/savour-labs/key-locker/config"
	"github.com/savour-labs/key-locker/db"
	"github.com/savour-labs/key-locker/model"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var cfg config.Config

func Execute() {
	app := cli.NewApp()
	app.Name = "key locker"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Usage: "Path of config yaml file",
			Value: "config.yml",
		},
	}

	app.Before = func(c *cli.Context) error {
		return config.LoadConfigFile(c.String("config"), &cfg)
	}

	app.Commands = []*cli.Command{
		{
			Name:  "start",
			Usage: "start api server",
			Action: func(c *cli.Context) error {
				server := api.NewServer(db.InitDB(cfg.Database), cfg.Server)
				server.Run()
				return nil
			},
		},
		{
			Name:  "init",
			Usage: "migrate database",
			Action: func(c *cli.Context) error {
				dba := db.InitDB(cfg.Database)
				if err := dba.AutoMigrate(&model.Key{}); err != nil {
					log.WithError(err).Fatal("Failed to migrate database")
					return err
				}
				return nil
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal("Failed to start application: %v", err)
	}
}
