package main

import (
	"fmt"
	"log"
	"os"

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

func runApp(c *cli.Context)error  {
	fmt.Println("anb")
	return nil
}