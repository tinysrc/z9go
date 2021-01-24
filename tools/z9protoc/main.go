package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

var (
	check       bool
	withGRPC    bool
	withGateway bool
	withSwagger bool
)

func main() {
	app := cli.NewApp()
	app.Name = "z9protoc"
	app.Usage = "z9 protoc"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "check",
			Destination: &check,
		},
		&cli.BoolFlag{
			Name:        "grpc",
			Usage:       "Whether to generate grpc.",
			Destination: &withGRPC,
		},
		&cli.BoolFlag{
			Name:        "gateway",
			Usage:       "Whether to generate gateway.",
			Destination: &withGateway,
		},
		&cli.BoolFlag{
			Name:        "swagger",
			Usage:       "Whether to generate swagger.",
			Destination: &withSwagger,
		},
	}
	app.Action = func(ctx *cli.Context) error {
		if check {
			err := checkPlugins()
			return err
		}
		files := ctx.Args()
		if len(files) == 0 {
			files, _ = filepath.Glob("*.proto")
		}
		if !withGRPC && !withGateway && !withSwagger {
			withGRPC = true
			withGateway = true
			withSwagger = true
		}
		if withGRPC {
			genGRPC(files)
		}
		if withGateway {
			genGateway(files)
		}
		if withSwagger {
			getSwagger(files)
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
