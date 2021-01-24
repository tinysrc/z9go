package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/tinysrc/z9go/tools/utils"
)

func protoc() (protoc string, err error) {
	z9root, err := utils.Z9Root()
	if err != nil {
		return
	}
	switch runtime.GOOS {
	case "darwin":
		protoc = path.Join(z9root, "third_party/protobuf/protoc-3.14.0-osx-x86_64")
	case "linux":
		protoc = path.Join(z9root, "third_party/protobuf/protoc-3.14.0-linux-x86_64")
	case "windows":
		protoc = path.Join(z9root, "third_party/protobuf/protoc-3.14.0-win64.exe")
	default:
		err = fmt.Errorf("unsupported os=%s", runtime.GOOS)
	}
	return
}

func firstInc() (inc string, err error) {
	z9root, err := utils.Z9Root()
	if err != nil {
		return
	}
	inc = path.Join(z9root, "third_party/protobuf/include")
	return
}

func secondInc() (inc string, err error) {
	mod, err := utils.LatestMod("github.com/grpc-ecosystem/grpc-gateway", "v2")
	if err != nil {
		return
	}
	inc = path.Join(mod, "third_party/googleapis")
	return
}

func baseArgs() (args []string, err error) {
	firstInc, err := firstInc()
	if err != nil {
		return
	}
	secondInc, err := secondInc()
	if err != nil {
		return
	}
	thirdInc, err := os.Getwd()
	if err != nil {
		return
	}
	args = []string{
		"-I", firstInc,
		"-I", secondInc,
		"-I", thirdInc,
	}
	return
}

func genGRPC(files []string) (err error) {
	protoc, err := protoc()
	if err != nil {
		return
	}
	args, err := baseArgs()
	if err != nil {
		return
	}
	extraArgs := []string{
		"--go_out", ".",
		"--go_opt", "paths=source_relative",
		"--go-grpc_out", ".",
		"--go-grpc_opt", "paths=source_relative",
	}
	args = append(args, extraArgs...)
	args = append(args, files...)
	err = utils.RunCmd(protoc, args...)
	return
}

func genGateway(files []string) (err error) {
	protoc, err := protoc()
	if err != nil {
		return
	}
	args, err := baseArgs()
	if err != nil {
		return
	}
	extraArgs := []string{
		"--grpc-gateway_out", ".",
		"--grpc-gateway_opt", "logtostderr=true",
		"--grpc-gateway_opt", "paths=source_relative",
	}
	args = append(args, extraArgs...)
	args = append(args, files...)
	err = utils.RunCmd(protoc, args...)
	return
}

func getSwagger(files []string) (err error) {
	protoc, err := protoc()
	if err != nil {
		return
	}
	args, err := baseArgs()
	if err != nil {
		return
	}
	extraArgs := []string{
		"--openapiv2_out", ".",
		"--openapiv2_opt", "logtostderr=true",
	}
	args = append(args, extraArgs...)
	args = append(args, files...)
	err = utils.RunCmd(protoc, args...)
	return
}
