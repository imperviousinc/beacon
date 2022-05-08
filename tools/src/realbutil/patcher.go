// based on https://github.com/brave/brave-core/blob/master/build/commands/lib/updatePatches.js

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Patcher is used for storing some metadata
// in .patcher file to maintain patches for Beacon
type Patcher struct {
	// Version is the version of this metadata file
	Version int `json:"version"`

	// LastPatched the hash of the last git commit
	// since patches were applied. This can be used
	// to detect if patches need to be re-applied
	// if patch files were modified upstream & git pulled
	// using: git diff --name-status <last-patched-commit> HEAD -- *.patch
	LastPatched string `json:"last_patched"`

	// rootPath the root directory contains .patcher file
	rootPath string
	// chromePath Chrome browser repo
	chromePath string
	// beaconPath Beacon repo
	beaconPath  string
	patchesPath string
}

// LoadPatcher loads .patcher metadata file or creates
// a default one.
func LoadPatcher(rootDir string) (*Patcher, error) {
	dotFile := filepath.Join(rootDir, ".patcher")

	raw, err := os.ReadFile(dotFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed reading .butil file: %w", err)
	}

	// empty or doesn't exist
	if len(raw) == 0 {
		return getDefaultPatcher(rootDir), nil
	}

	p := getDefaultPatcher(rootDir)
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("failed reading .butil file: %v", err)
	}

	return p, nil
}

func (p *Patcher) applyPatch(name string) error {
	patchPath := filepath.Join(p.patchesPath, name)
	c := exec.Command("git", "apply", patchPath)
	c.Dir = p.chromePath

	out, err := c.CombinedOutput()
	if err != nil {
		if p.patchApplied(name) {
			return nil
		}

		fmt.Println(string(out))
		return err
	}
	return nil
}

func (p *Patcher) patchApplied(name string) bool {
	patchPath := filepath.Join(p.patchesPath, name)
	c := exec.Command("git", "apply", "--reverse", "--check", patchPath)
	c.Dir = p.chromePath

	return c.Run() == nil
}

// Apply applies patches
func (p *Patcher) Apply() error {
	patches, err := os.ReadDir(p.patchesPath)
	if err != nil {
		return fmt.Errorf("failed reading patches dir: %v", err)
	}

	failed := 0
	for _, patch := range patches {
		if patch.IsDir() {
			continue
		}

		fmt.Println("Applying patch " + patch.Name())
		if err = p.applyPatch(patch.Name()); err != nil {
			fmt.Printf("[ERROR] Failed applying %s: %v\n", patch.Name(), err)
			failed++
		}
	}

	if failed > 0 {
		return fmt.Errorf("failed applying patches ["+
			"total: %d, failed: %d]", len(patches), failed)
	}

	if err := p.UpdateLastPatched(); err != nil {
		return err
	}
	return p.PersistConfig()
}

func (p *Patcher) Update() error {
	updatePatches(p.chromePath, filepath.Join(p.beaconPath, "patches"))
	return nil
}

// Reverse clears patches
func (p *Patcher) Reverse() error {
	return nil
}

func (p *Patcher) UpdateLastPatched() error {
	c := exec.Command("git", "rev-parse", "HEAD")
	c.Dir = p.beaconPath

	commitHash, err := c.Output()
	if err != nil {
		return fmt.Errorf("call failed: %s: %v", c.String(), err)
	}

	p.LastPatched = strings.TrimSpace(string(commitHash))
	return nil
}

// PersistConfig stores the metadata .patcher file
func (p *Patcher) PersistConfig() error {
	raw, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("persist failed: %v", err)
	}

	dotFile := filepath.Join(p.rootPath, ".patcher")
	if err = writeFile(dotFile, raw); err != nil {
		return fmt.Errorf("persist failed: %v", err)
	}
	return nil
}

func (p *Patcher) readAllPatches() ([]string, error) {
	entires, err := os.ReadDir(p.patchesPath)
	if err != nil {
		return nil, err
	}

	var patchFiles []string
	for _, entry := range entires {
		if filepath.Ext(entry.Name()) == ".patch" {
			patchName := filepath.Join(p.patchesPath, entry.Name())
			patchFiles = append(patchFiles, patchName)
		}
	}
	return patchFiles, nil
}

func (p *Patcher) ShouldReapply() (bool, error) {
	added, modified, err := p.DiffSinceLastPatched()
	return len(added) > 0 || len(modified) > 0, err
}

// DiffSinceLastPatched returns a list of modified patch files and a list
// of newly added patches since the last time Apply or Update were called.
// This may happen when the repo containing the patch files was updated
// upstream and git pulled.
func (p *Patcher) DiffSinceLastPatched() ([]string, []string, error) {
	if p.LastPatched == "" {
		added, err := p.readAllPatches()
		return added, nil, err
	}
	c := exec.Command("git", "diff", "--name-status", p.LastPatched, "HEAD", "--", "*.patch")
	c.Dir = p.beaconPath
	fmt.Println(p.beaconPath)
	if p.beaconPath == "" {
		panic("wat")
	}
	out, err := c.CombinedOutput()
	if err != nil {
		return nil, nil, fmt.Errorf("call failed: %s: %v:%s", c.String(), err, string(out))
	}

	output := strings.Split(strings.TrimSpace(string(out)), "\n")
	var added, modified []string
	for _, diff := range output {
		parts := strings.Fields(diff)
		if len(parts) != 2 {
			continue
		}

		if parts[0] == "A" {
			added = append(added, parts[1])
		}
		if parts[0] == "M" {
			modified = append(modified, parts[1])
		}
	}

	added = p.removeApplied(added)
	modified = p.removeApplied(modified)

	return added, modified, nil
}

func (p *Patcher) removeApplied(patches []string) []string {
	var notApplied []string
	for _, patch := range patches {
		_, file := filepath.Split(patch)
		if !p.patchApplied(file) {
			notApplied = append(notApplied, patch)
		}
	}

	return notApplied
}

func getDefaultPatcher(rootDir string) *Patcher {
	chromeDir := filepath.Join(rootDir, "src")
	mustExist(chromeDir)
	beaconDir := filepath.Join(chromeDir, "beacon")
	mustExist(beaconDir)
	patchesDir := filepath.Join(beaconDir, "patches")
	mustExist(patchesDir)
	return &Patcher{
		Version:     0,
		LastPatched: "",
		rootPath:    rootDir,
		chromePath:  chromeDir,
		beaconPath:  beaconDir,
		patchesPath: patchesDir,
	}
}

// getChromeDiffNamesOnly gets all files modified in chrome
// excluding with shouldSaveAsPatch
func getChromeDiffNamesOnly(gitRepoDir string) []string {
	cmd := exec.Command("git", "diff",
		"--diff-filter", "M", "--name-only", "--ignore-space-at-eol")
	cmd.Dir = gitRepoDir

	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed reading diff for repo: %v", gitRepoDir)
	}

	paths := strings.Split(string(out), "\n")
	var filtered []string
	for _, m := range paths {
		m := strings.TrimSpace(m)
		if shouldSaveAsPatch(m) {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

// getPatchNameFromPath flattens a path into a file
// name ending with a .patch extension
func getPatchNameFromPath(path string) string {
	patch := filepath.ToSlash(filepath.Clean(path))
	return strings.ReplaceAll(patch, "/", "-") + ".patch"
}

func writePatchFiles(modifiedPaths []string, gitRepoPath, patchDirPath string) []string {
	var patchFiles []string
	for _, m := range modifiedPaths {
		patch := getPatchNameFromPath(m)

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
		log.Printf("removing stale patch `%s`", oldPatch.Name())
		if err := os.Remove(oldPatchPath); err != nil {
			log.Fatalf("failed removing stale patch `%s`: %v", oldPatch.Name(), err)
		}
	}
}

func updatePatches(gitRepoPath, patchDirPath string) {
	diff := getChromeDiffNamesOnly(gitRepoPath)
	patchFileNames := writePatchFiles(diff, gitRepoPath, patchDirPath)
	removeStalePatchFiles(patchFileNames, patchDirPath, []string{} /*no exclusions*/)
}

var shouldSaveAsPatch = func(s string) bool {
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
