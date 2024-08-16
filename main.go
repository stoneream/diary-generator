package main

import (
	"diary-generator/cmd/archive"
	"diary-generator/cmd/initialize"
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
						Name:     "base-directory",
						Usage:    "base directory path",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "template-file",
						Usage:    "template file path",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					cmd := initialize.InitializeCmd{
						BaseDirectory: c.String("base-directory"),
						TemplateFile:  c.String("template-file"),
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
						Name:     "base-directory",
						Usage:    "base directory path",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "starts-with",
						Usage:    "starts with string",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					cmd := archive.ArchiveCmd{
						BaseDirectory: c.String("base-directory"),
						StartsWith:    c.String("starts-with"),
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
