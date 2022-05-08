package main

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"time"
)

func run(actionName, dir, exePath string, args ...string) error {
	fmt.Println(actionName)
	cmd := newCmd(exePath, args...)
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: %v", cmd.String(), err)
	}
	return nil
}

func hasString(arr []string, s string) bool {
	for _, c := range arr {
		if c == s {
			return true
		}
	}
	return false
}

func fixOverrideTime(chromePath, beaconPath string) error {
	chromeStat, err := os.Stat(chromePath)
	if err != nil {
		return err
	}
	beaconStat, err := os.Stat(beaconPath)
	if err != nil {
		return err
	}
	return compareFixOverrideTime(chromePath, beaconPath, chromeStat.ModTime(), beaconStat.ModTime())
}

func touchCCOverrides(chromePath, beaconPath string) {
	_ = filepath.Walk(beaconPath, func(path string, overrideInfo fs.FileInfo, err error) error {
		if err != nil || overrideInfo.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".cc" && ext != ".h" && ext != ".mm" {
			return nil
		}

		relPath, err := filepath.Rel(beaconPath, path)
		if err != nil {
			log.Printf("[WARNING] failed getting relative path for chrome_src from: %v", err)
			return nil
		}

		chromePath := filepath.Join(chromePath, relPath)
		originalInfo, err := os.Stat(chromePath)
		if err != nil {
			log.Printf("[WARNING] no matching chrome override found for `%s` since `%s` isn't available: %v",
				path, chromePath, err)
			return nil
		}
		if err = compareFixOverrideTime(chromePath, path, originalInfo.ModTime(),
			overrideInfo.ModTime()); err != nil {
			log.Printf("[WARNING] %v", err)
		}

		return nil
	})
}

func compareFixOverrideTime(chromePath, beaconPath string, chromeTime time.Time, beaconTime time.Time) error {
	if chromeTime.Equal(beaconTime) {
		return nil
	}
	// ensures both have the same time
	fixTime := time.Now().Round(time.Second)
	log.Printf("Touching override `%s`", chromePath)
	if err := os.Chtimes(chromePath, fixTime, fixTime); err != nil {
		return fmt.Errorf("failed changing mod time for chrome file `%s`:%v", chromePath, err)

	}
	if err := os.Chtimes(beaconPath, fixTime, fixTime); err != nil {
		return fmt.Errorf("failed changing mod time for beacon file `%s`:%v", beaconPath, err)
	}

	return nil
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

func check(err error) {
	if err != nil {
		fmt.Println(err)
		debug.PrintStack()
		os.Exit(1)
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
