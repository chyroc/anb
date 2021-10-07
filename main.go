package main

import (
	"log"
	"os"

	"github.com/chyroc/anb/internal/app"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "anb",
		Commands: nil,
		Flags:    nil,
		Action:   runApp,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func runApp(c *cli.Context) error {
	config := c.Args().First()

	return app.Run(&app.RunRequest{
		Config: config,
	})
}
