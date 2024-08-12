package main

import (
	"os"

	"github.com/joho/godotenv"

	"flowspell/cron"
	"flowspell/server"
	"flowspell/worker"

	"github.com/urfave/cli"
)

var (
	app *cli.App
)

func init() {
	// Initialise a CLI app
	app = cli.NewApp()
	app.Name = "FlowSpell"
	app.Usage = "FlowSpell is a magical tool to handle the flows of your applications"
	app.Version = "1.0.0"
}

func main() {
	godotenv.Load()

	// Set the CLI app commands
	app.Commands = []cli.Command{
		{
			Name:  "worker",
			Usage: "Start FlowSpell worker",
			Action: func(c *cli.Context) error {
				worker.Worker()
				return nil
			},
		},
		{
			Name:  "server",
			Usage: "Start FlowSpell server",
			Action: func(c *cli.Context) error {
				server.FlowSpellServer()
				return nil
			},
		},
		{
			Name:  "cron",
			Usage: "Start FlowSpell scheduler",
			Action: func(c *cli.Context) error {
				cron.Start()
				return nil
			},
		},
	}

	// Run the CLI app
	_ = app.Run(os.Args)
}
