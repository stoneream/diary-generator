package main

import (
	"diary-generator/cmd/archive"
	"diary-generator/cmd/initialize"
	"diary-generator/config"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "diary-generator",
		Usage: "",
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "Initialize a diary",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Usage:    "config file path",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					config, err := config.LoadFile(c.String("config"))

					if err != nil {
						return err
					}

					cmd := initialize.InitializeCmd{
						BaseDirectory: config.BaseDirectory,
						TemplateFile:  config.TemplateFile,
						Name:          config.Name,
					}
					return cmd.Execute()
				},
			},
			{
				Name:    "archive",
				Aliases: []string{"a"},
				Usage:   "Archive a diary",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Usage:    "config file path",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "starts-with",
						Usage:    "starts with string",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					config, err := config.LoadFile(c.String("config"))

					if err != nil {
						return err
					}

					cmd := archive.ArchiveCmd{
						BaseDirectory: config.BaseDirectory,
						StartsWith:    c.String("starts-with"), // TODO メタデータをもとにアーカイブする
					}
					return cmd.Execute()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
