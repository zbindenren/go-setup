package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
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
	- directory /tmp/golang/src/github.com
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

	if clean {
		err := cleanUp(gopath, packagePath)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	src := "."
	if len(os.Getenv(srcPath)) > 0 {
		src = os.Getenv(srcPath)
	}
	msg, err := setup(gopath, packagePath, src)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(msg)

}

func cleanUp(gopath, packagePath string) error {
	new := newName(gopath, packagePath)
	linkExists, err := isLink(new)
	if err != nil {
		return fmt.Errorf("could not determine if %s is link: %s", new, err)
	}
	if linkExists {
		os.RemoveAll(gopath)
	}
	return nil
}

func setup(gopath, packagePath, src string) (string, error) {
	new := newName(gopath, packagePath)
	oldName, err := filepath.Abs(src)
	if err != nil {
		return "", fmt.Errorf("could not detect working directory: %s", err)
	}
	linkExists, err := isLink(new)
	if err != nil {
		return "", fmt.Errorf("could not determine if %s is link: %s", new, err)
	}
	if linkExists {
		return fmt.Sprintf("link %s -> %s already exists, nothing to do", new, oldName), nil
	}
	parentDir := filepath.Join(filepath.Dir(new))
	err = os.MkdirAll(parentDir, 0755)
	if err != nil {
		return "", fmt.Errorf("could not create directory %s: %s", parentDir, err)
	}
	err = os.Symlink(path.Dir(oldName), new)
	if err != nil {
		return "", fmt.Errorf("could not create symlink %s -> %s: %s", new, oldName, err)
	}
	return fmt.Sprintf("created symbolic link %s -> %s", new, oldName), nil
}

// newName is the link name.
func newName(gopath, packagePath string) string {
	return path.Dir(filepath.Join(gopath, "src", packagePath))
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
