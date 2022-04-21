package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var debug = false

func main() {
	args := os.Args[1:]
	beaconPath := replaceCC(args)

	if wrapper := os.Getenv("BEACON_CC_WRAPPER"); wrapper != "" {
		args = append([]string{wrapper}, args...)
	}

	if ccDebug := os.Getenv("BEACON_CREDIRECT_DEBUG"); strings.EqualFold(ccDebug, "true") {
		debug = true
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
			return
		}
		os.Exit(1)
		return
	}

	// To check the redirected file timestamp, it should be marked as a dependency
	// for ninja. Linux/MacOS gcc deps format includes this file properly.
	// Windows msvc deps format does not include it, so we do it manually here.
	if runtime.GOOS == "windows" && beaconPath != "" {
		// This is a specially crafted string that ninja will look for to create
		// deps.
		fmt.Fprintf(os.Stderr, "Note: including file: %s\n", beaconPath)
	}
}

func replaceCC(args []string) string {
	// find -c <path>.cc
	cArgIdx := -1
	for i, arg := range args {
		if arg == "/c" || arg == "-c" {
			cArgIdx = i
			break
		}
	}
	if len(args) == 0 || cArgIdx == -1 || len(args)-1 == cArgIdx {
		return ""
	}

	// get <path>.cc
	ccFile := args[cArgIdx+1]
	ccFile, err := filepath.Abs(ccFile)
	if err != nil {
		return ""
	}

	exePath, err := os.Executable()
	if err != nil {
		return ""
	}

	binDir, exeName := filepath.Split(exePath)
	if !strings.HasPrefix(exeName, "credirect") {
		log.Fatalf("bad executable name: %s, want: %s", exeName, "credirect")
	}

	// Path of this executable MUST be:
	// /some/path/to/chromium_root/tools/credirect/bin/credirect
	wantSuffix := filepath.Clean("tools/bin") + string(os.PathSeparator)
	if !strings.HasSuffix(binDir, wantSuffix) {
		log.Fatalf("got path: `%s`, want credirect wrapper to have this path suffix: `%s`", binDir, wantSuffix)
	}

	// Getting chromium root from: /some/path/to/chromium_root/src/beacon/tools/bin/credirect
	// chromium dir is at /some/path/to/chromium_root/src
	chromiumDir := filepath.Join(filepath.Clean(binDir + "../../../"))
	if debug {
		log.Printf("attempt redirect file: `%s`, with chromium dir: `%s`", ccFile, chromiumDir)
	}

	if len(chromiumDir) >= len(ccFile) {
		return ""
	}
	if !strings.HasPrefix(ccFile, chromiumDir) {
		// source file doesn't share a common path with chromium dir
		if debug {
			log.Printf("file: `%s` does not share a common path with: `%s`", ccFile, chromiumDir)
		}
		return ""
	}

	ccRel := ccFile[len(chromiumDir)+1:]
	beaconPath := filepath.Join(chromiumDir, "beacon", "chromium_src", ccRel)

	if _, err := os.Stat(beaconPath); err == nil {
		if debug {
			log.Printf("successfully redirecting `%s` to `%s`", ccFile, beaconPath)
		}
		args[cArgIdx+1] = beaconPath
		return beaconPath
	}

	return ""
}
