package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var clean bool

const (
	goPath  = "GOPATH"
	pkgPath = "PACKAGE_PATH"
	srcPath = "SRC_PATH"
)

const usage = `
Small utility to setup a golang project for CI environments.

Example:
	GOPATH=/tmp/golang PACKAGE_PATH=github.com/zbindenren/go-setup go-setup

creates:
	- directory /tmp/golang/src/github.com/zbindenren
	- link: /tmp/golang/src/github.com/zbindenren -> CWD

Flags:
`

func init() {
	flag.BoolVar(&clean, "clean", false, "cleanup")
}
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
		flag.PrintDefaults()
	}
	flag.Parse()
	gopath := os.Getenv(goPath)
	packagePath := os.Getenv(pkgPath)
	if len(gopath) == 0 {
		log.Fatalf("environment variable %s is not defined or empty string", goPath)
	}
	if len(packagePath) == 0 {
		log.Fatalf("environment variable %s is not defined or empty string", pkgPath)
	}
	newName := filepath.Join(gopath, "src", packagePath)

	linkExists, err := isLink(newName)
	if err != nil {
		log.Fatalf("could not determine if %s is link: %s", newName, err)
	}
	if clean {
		if linkExists {
			os.RemoveAll(gopath)
		}
		log.Printf("removed directory %s", gopath)
		os.Exit(0)
	}
	src := "."
	if len(os.Getenv(srcPath)) > 0 {
		src = os.Getenv(srcPath)
	}

	oldName, err := filepath.Abs(src)
	if err != nil {
		log.Fatalf("could not detect working directory: %s", err)
	}
	if linkExists {
		log.Printf("link %s -> %s already exists, nothing to do", newName, oldName)
		os.Exit(0)
	}
	parentDir := filepath.Join(filepath.Dir(newName))
	err = os.MkdirAll(parentDir, 0755)
	if err != nil {
		log.Fatalf("could not create directory %s: %s", parentDir, err)
	}
	log.Printf("created directory %s", parentDir)
	err = os.Symlink(oldName, newName)
	if err != nil {
		log.Fatalf("could not create symlink %s -> %s: %s", newName, oldName, err)
	}
	log.Printf("created symbolic link %s -> %s", newName, oldName)
}

func isLink(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}
	fi, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return fi.Mode()&os.ModeSymlink == os.ModeSymlink, nil
}
