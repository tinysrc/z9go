package main

import (
	"github.com/tinysrc/z9go/tools/utils"
)

var plugins = []plugin{
	{
		name: "protoc-gen-grpc-gateway",
		url:  "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway",
	},
	{
		name: "protoc-gen-openapiv2",
		url:  "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2",
	},
	{
		name: "protoc-gen-go",
		url:  "google.golang.org/protobuf/cmd/protoc-gen-go",
	},
	{
		name: "protoc-gen-go-grpc",
		url:  "google.golang.org/grpc/cmd/protoc-gen-go-grpc",
	},
}

type plugin struct {
	name string
	url  string
}

func checkPlugins() (err error) {
	for _, plugin := range plugins {
		if err = utils.RunCmd("go", "get", "-u", plugin.url); err != nil {
			return
		}
	}
	return
}
