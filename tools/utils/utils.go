package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd run command
func RunCmd(name string, args ...string) (err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	cmd := exec.Command(name, args...)
	cmd.Dir = wd
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println(cmd.String())
	return cmd.Run()
}
