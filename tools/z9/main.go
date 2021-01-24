package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "z9"
	app.Usage = "z9 tools"
	app.Commands = []cli.Command{
		{
			Name:   "update",
			Usage:  "Update engine",
			Action: runUpdate,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
