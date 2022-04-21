package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const uninitializedRepoDir = ".beacon__repo"

// findRootDirectory returns the directory path
// where .butil is located starting from searchDir and looking
// back at parent dirs. Returns an empty string if it can't find it.
func findRootDirectory(searchDir string) string {
	var butilPath string

	for {
		possiblePath := filepath.Join(searchDir, ".butil")
		if fileExists(possiblePath) {
			butilPath = possiblePath
			break
		}

		oldSearchPath := searchDir
		searchDir = filepath.Dir(searchDir)
		if oldSearchPath == searchDir {
			break
		}
	}

	if butilPath != "" {
		return searchDir
	}
	return ""
}

func exeName(name string) string {
   ext := ""
   if runtime.GOOS == "windows" {
      ext = ".exe"
   }
   return name + ext
}

func callRealUtil(beaconDir string, workingDir string) {
	maybeBuildTools(beaconDir)
        exePath := filepath.Join(beaconDir, "tools", "bin", exeName("realbutil"))
	mustExist(exePath)
	_ = run(exePath, os.Args[1:]...)
	return
}

func main() {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if fileExists(uninitializedRepoDir) {
		callRealUtil(filepath.Join(uninitializedRepoDir), workingDir)
		return
	}

	// Either find its path by finding .butil
	// or bootstrap it.
	rootDir := findRootDirectory(workingDir)
	if rootDir == "" {
		bootstrap(workingDir)
		return
	}

	callRealUtil(filepath.Join(rootDir, "src", "beacon"), rootDir)
}

func argOrDefault(id int, defaultValue string) string {
	if len(os.Args) > id {
		return os.Args[id]
	}
	return defaultValue
}

// bootstrap clones and builds butil returning its path
func bootstrap(dir string) string {
	beaconDir := cloneBeacon(dir)
	maybeBuildTools(beaconDir)
	butilPath := filepath.Join(beaconDir, "tools", "bin", exeName("realbutil"))
	mustExist(butilPath)
	return butilPath
}

func shouldRebuild(dir string, lastBuild time.Time) bool {
	rebuild := false
	_ = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.ModTime().After(lastBuild) {
			rebuild = true
			return errors.New("stop")
		}
		return nil
	})
	return rebuild
}

func maybeBuildTools(beaconDir string) {
	toolsDir := filepath.Join(beaconDir, "tools")
	toolsSrc := filepath.Join(toolsDir, "src")
	mustExist(toolsSrc)

	toolsBin := filepath.Join(toolsDir, "bin")
	must(os.MkdirAll(toolsBin, 0700))
	mustExist(toolsBin)

	// find all projects in tools/src
	tools, err := os.ReadDir(toolsSrc)
	if err != nil {
		log.Fatalf("failed listing tools dir `%s`: %v", toolsDir, err)
	}

	for _, tool := range tools {
                if !tool.IsDir() {
                    continue
                }
		toolSrc := filepath.Join(toolsSrc, tool.Name())
		toolBin := filepath.Join(toolsBin, tool.Name())
	        toolBin = exeName(toolBin)
		binStat, err := os.Stat(toolBin)
		if err == nil && !shouldRebuild(toolSrc, binStat.ModTime()) {
			continue
		}

		fmt.Println("Building tools/" + tool.Name())
		c := newCmd("go", "build", "-o", toolBin)
		c.Dir = toolSrc
		must(c.Run())
	}
}

func cloneBeacon(dir string) string {
	if argOrDefault(1, "") != "clone" {
		fmt.Println("Beacon utility is not initialized in this directory")
		fmt.Println("Run: butil clone [url] [options]")
		os.Exit(1)
	}

	if empty, err := isDirEmpty(dir); err != nil || !empty {
		fmt.Println("Initialize failed directory must be empty")
		os.Exit(1)
	}

	cloneUrl := argOrDefault(2, "https://github.com/imperviousinc/beacon")
	var additionalArgs []string
	if len(os.Args) > 3 {
		additionalArgs = os.Args[3:]
	}

	args := []string{"clone", cloneUrl, uninitializedRepoDir}

	if len(additionalArgs) > 0 {
		args = append(args, additionalArgs...)
	}

	fmt.Println("Cloning ", cloneUrl)
	must(run("git", args...))
	mustExist(uninitializedRepoDir)
	beaconDir, err := filepath.Abs(uninitializedRepoDir)
	must(err)
	return beaconDir
}
