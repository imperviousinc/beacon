// based on https://github.com/brave/brave-core/blob/master/build/commands/lib/updatePatches.js

package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func touchOverriddenLogos(themeDir string) {
	_ = filepath.Walk(themeDir, func(path string, overrideInfo fs.FileInfo, err error) error {
		if err != nil || overrideInfo.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".png" && ext != ".icon" {
			return nil
		}

		if err := os.Chtimes(path, time.Now(), time.Now()); err != nil {
			log.Printf("[WARNING] failed changing mod time for logo override `%s`:%v", path, err)
			return nil
		}

		log.Printf("Touching original logo file `%s`", path)
		return nil
	})
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

func compareFixOverrideTime(chromePath, beaconPath string, chromeTime time.Time, beaconTime time.Time) error {
	if chromeTime.Equal(beaconTime) {
		return nil
	}
	// ensures both have the same time
	fixTime := time.Now().Round(time.Second)
	log.Printf("Touching override `%s`", chromePath)
	if err := os.Chtimes(chromePath, fixTime, fixTime); err != nil {
		return fmt.Errorf("failed changing mod time for chromium file `%s`:%v", chromePath, err)

	}
	if err := os.Chtimes(beaconPath, fixTime, fixTime); err != nil {
		return fmt.Errorf("failed changing mod time for beacon file `%s`:%v", beaconPath, err)
	}

	return nil
}

func touchCCOverrides(chromiumSrcDir, beaconChromiumSrcDir string) {
	_ = filepath.Walk(beaconChromiumSrcDir, func(path string, overrideInfo fs.FileInfo, err error) error {
		if err != nil || overrideInfo.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".cc" && ext != ".h" && ext != ".mm" {
			return nil
		}

		relPath, err := filepath.Rel(beaconChromiumSrcDir, path)
		if err != nil {
			log.Printf("[WARNING] failed getting relative path for chromium_src from: %v", err)
			return nil
		}

		chromePath := filepath.Join(chromiumSrcDir, relPath)
		originalInfo, err := os.Stat(chromePath)
		if err != nil {
			// if ext == ".h" {
			// 	return nil
			// }
			log.Printf("[WARNING] no matching chromium override found for `%s` since `%s` isn't available: %v",
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

func getModifiedPaths(gitRepoDir string) []string {
	cmd := exec.Command("git", "diff",
		"--diff-filter", "M", "--name-only", "--ignore-space-at-eol")
	cmd.Dir = gitRepoDir

	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed reading diff for repo: %v", gitRepoDir)
	}

	return strings.Split(string(out), "\n")
}

func writePatchFiles(modifiedPaths []string, gitRepoPath, patchDirPath string) []string {
	var patchFiles []string
	for _, m := range modifiedPaths {
		patch := filepath.ToSlash(filepath.Clean(m))
		patch = strings.ReplaceAll(patch, "/", "-") + ".patch"

		cmd := exec.Command("git", "diff",
			"--src-prefix=a/", "--dst-prefix=b/", "--full-index", m)
		cmd.Dir = gitRepoPath

		diff, err := cmd.Output()
		if err != nil {
			log.Fatalf("failed reading diff for file `%s`: %v", m, err)
		}

		dest := filepath.Join(patchDirPath, patch)
		log.Printf("writing patch file %s", patch)
		if err := os.WriteFile(dest, diff, 0644); err != nil {
			log.Fatalf("failed writing patch `%s`:%v", patch, err)
		}

		patchFiles = append(patchFiles, dest)
	}

	return patchFiles
}

func hasString(arr []string, s string) bool {
	for _, c := range arr {
		if c == s {
			return true
		}
	}
	return false
}

func removeStalePatchFiles(patchFileNames []string, patchDirPath string, keepPatchFilenames []string) {
	allPatches, err := os.ReadDir(patchDirPath)
	if err != nil {
		log.Fatalf("failed removing stale patches: %v", err)
	}

	for _, oldPatch := range allPatches {
		oldPatchPath := filepath.Join(patchDirPath, oldPatch.Name())
		if oldPatch.IsDir() ||
			hasString(patchFileNames, oldPatchPath) ||
			hasString(keepPatchFilenames, oldPatchPath) {
			continue
		}

		// remove stale patch
		log.Printf("[WARNING] removing stale patch `%s`", oldPatch.Name())
		if err := os.Remove(oldPatchPath); err != nil {
			log.Fatalf("failed removing stale patch `%s`: %v", oldPatch.Name(), err)
		}
	}
}

func updatePatches(gitRepoPath, patchDirPath string, repoPathFilter func(path string) bool, keepPatchFilenames []string) {
	modified := getModifiedPaths(gitRepoPath)
	var filtered []string
	for _, m := range modified {
		m := strings.TrimSpace(m)
		if repoPathFilter(m) {
			filtered = append(filtered, m)
		}
	}

	patchFileNames := writePatchFiles(filtered, gitRepoPath, patchDirPath)
	removeStalePatchFiles(patchFileNames, patchDirPath, keepPatchFilenames)
}

var chromiumPathFilter = func(s string) bool {
	return len(s) > 0 &&
		!strings.HasPrefix(s, "chrome/app/theme/default") &&
		!strings.HasPrefix(s, "chrome/app/theme/beacon") &&
		!strings.HasPrefix(s, "chrome/app/theme/chromium") &&
		!strings.HasSuffix(s, ".png") &&
		!strings.HasSuffix(s, ".xtb") &&
		!strings.HasSuffix(s, ".grdp") &&
		!strings.HasSuffix(s, ".grd") &&
		!strings.HasSuffix(s, ".svg") &&
		!strings.HasSuffix(s, ".icon") &&
		!strings.HasSuffix(s, "channel_constants.xml") &&
		!strings.HasSuffix(s, ".min.js") &&
		!strings.Contains(s, "google_update_idl")

}
