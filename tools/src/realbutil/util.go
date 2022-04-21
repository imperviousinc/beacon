package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func runAction(actionName, exePath string, args ...string) error {
	fmt.Println(actionName)
	return run(exePath, args...)
}

func runActionWithDir(actionName, dir, exePath string, args ...string) error {
	fmt.Println(actionName)
	cmd := newCmd(exePath, args...)
	cmd.Dir = dir
	return cmd.Run()
}

func run(name string, args ...string) error {
	return newCmd(name, args...).Run()
}

func newCmd(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func mustExist(path string) {
	if _, err := os.Stat(path); err != nil {
		log.Fatalf("path `%s` doesn't exist: %v", path, err)
	}
}

func mustNotExist(path string) {
	if _, err := os.Stat(path); err == nil {
		log.Fatalf("path `%s` already exists", path)
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// writeFile creates path if it doesn't exist and writes
// data to path.
func writeFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}
