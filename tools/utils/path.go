package utils

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// GoPath GOPATH
func GoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

// LatestMod get lastest mod
func LatestMod(prefix string, mod string) (string, error) {
	gopath := GoPath()
	base := path.Join(gopath, "pkg/mod", prefix)
	files, err := ioutil.ReadDir(base)
	if err != nil {
		fmt.Printf("LatestMod failed error=%v\n", err)
		return "", err
	}
	for i := len(files) - 1; i >= 0; i-- {
		if strings.HasPrefix(files[i].Name(), mod+"@") {
			return path.Join(base, files[i].Name()), nil
		}
	}
	err = fmt.Errorf("not found mod=%s", path.Join(base, mod))
	fmt.Printf("LatestMod failed error=%v\n", err)
	return "", err
}

// Z9Root z9 root
func Z9Root() (string, error) {
	z9root := os.Getenv("Z9ROOT")
	if z9root != "" {
		return z9root, nil
	}
	mod, err := LatestMod("github.com/tinysrc", "z9go")
	if err != nil {
		return "", err
	}
	return mod, nil
}
