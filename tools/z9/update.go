package main

import (
	"github.com/tinysrc/z9go/tools/utils"
	"github.com/urfave/cli"
)

func runUpdate(ctx *cli.Context) (err error) {
	if err = utils.RunCmd("go", "get", "-u", "github.com/tinysrc/z9go/tools/z9protoc"); err != nil {
		return err
	}
	if err = utils.RunCmd("z9protoc", "--check"); err != nil {
		return err
	}
	return
}
