package main

import (
	"log"
	"os"
	"time"

	"github.com/stoneream/diary-generator/v2/cmd/archive"
	"github.com/stoneream/diary-generator/v2/cmd/initialize"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "diary-generator",
		Usage: "",
		Flags: []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "Initialize a diary",
				Action: func(c *cli.Context) error {
					cmd := initialize.InitializeCmd{
						Now: time.Now(),
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
						Name:     "target-ym",
						Usage:    "target year and month (e.g. 2024-01)",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					targetYM := c.String("target-ym")

					cmd := archive.ArchiveCmd{
						TargetYM: targetYM,
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
