package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// setupPatches creates dummy chrome and beacon repos
// in a tmp directory to test patching
// TODO: stub git calls instead
func setupPatches(t *testing.T) *Patcher {
	tmpDir, err := os.MkdirTemp("", "patcher")
	if err != nil {
		log.Fatalf("failed making tmp tmpDir: %v", err)
	}

	// register cleanup function
	// to remove the tmp directory
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	chromeDir := filepath.Join(tmpDir, "src")
	beaconDir := filepath.Join(chromeDir, "beacon")
	patchesDir := filepath.Join(beaconDir, "patches")

	check(os.MkdirAll(chromeDir, 0700))
	check(os.MkdirAll(beaconDir, 0700))
	check(os.MkdirAll(patchesDir, 0700))

	gitInit := func(dir string) {
		c := exec.Command("git", "init")
		c.Dir = dir
		check(c.Run())
		// set fake user/email for this test repo
		c = exec.Command("git", "config", "user.name", "Dummy user")
		c.Dir = dir
		check(c.Run())
		c = exec.Command("git", "config", "user.email", "foo@example.invalid")
		c.Dir = dir
		check(c.Run())
	}

	// dummy chrome repo
	gitInit(chromeDir)
	updateFile(chromeDir, "example.cc", "example file\nfor testing\n....\n")
	updateFile(chromeDir, "example2.cc", "example 2 file\nfor testing\n....\n")
	// This file should be ignored by shouldSaveAsPatch
	updateFile(chromeDir, "foo.png", "totally valid png format")
	gitAdd(chromeDir, ".")
	gitCommit(chromeDir, "-m", "test files")

	// Dummy beacon repo
	gitInit(beaconDir)
	updateFile(beaconDir, "test.cc", "baz bar")
	updateFile(beaconDir, "hello.cc", "foo bar")
	gitAdd(beaconDir, ".")
	gitCommit(beaconDir, "-m", "dummy files")

	p, err := LoadPatcher(tmpDir)
	check(err)

	return p
}

func TestPatcher_Apply(t *testing.T) {
	p := setupPatches(t)
	// no patches yet apply should succeed
	if err := p.Apply(); err != nil {
		t.Fatal(err)
	}

	// make some patches
	updateFile(p.chromePath, "example.cc", "foo\nbar")
	check(p.Update())

	// reset
	c := exec.Command("git", "reset", "--hard")
	c.Dir = p.chromePath
	check(c.Run())

	if err := p.Apply(); err != nil {
		t.Fatal(err)
	}

	if readFile(p.chromePath, "example.cc") != "foo\nbar" {
		t.Fatal("patch wasn't applied")
	}

	// modify file in chrome so that patch is no longer applicable
	updateFile(p.chromePath, "example.cc", "bar dddd")
	if err := p.Apply(); err == nil {
		t.Fatal("got no error on non-applicable patch")
	}

	// re-write the patch file and re-apply
	check(p.Update())
	check(p.Apply())
}

func TestPatcher_DiffSinceLastPatched(t *testing.T) {
	p := setupPatches(t)
	if p.LastPatched != "" {
		t.Fatalf("got last patched = %s, but it shouldn't be set initially", p.LastPatched)
	}
	// Empty patches in beacon repo (created in setupPatches call)
	// apply will set LastPatched = <last commit hash/HEAD>
	check(p.Apply())
	prevLastPatched := p.LastPatched
	if prevLastPatched == "" {
		t.Fatalf("got last patched = '', want a commit hash")
	}

	// create example.cc.patch & commit changes
	updateFile(p.chromePath, "example.cc", "foo")
	check(p.Update())

	// diff should be empty we only want
	// it to change when some patches were "committed" into the repo
	// pulled or modified from upstream repo and need to be
	// applied
	added, modified, err := p.DiffSinceLastPatched()
	if err != nil || len(added) > 0 || len(modified) > 0 {
		t.Fatalf("want no diff")
	}

	gitAdd(p.beaconPath, ".")
	gitCommit(p.beaconPath, "-m", "add more")

	// Some patches were committed, but they were already applied
	// so no diff
	added, modified, err = p.DiffSinceLastPatched()
	if err != nil || len(added) > 0 || len(modified) > 0 {
		t.Fatalf("want no diff")
	}

	check(p.Apply())
	added, modified, err = p.DiffSinceLastPatched()
	if err != nil || len(added) > 0 || len(modified) > 0 {
		t.Fatalf("want no diff")
	}

	// replace foo with bar in the patch file
	// to simulate the patch being modified
	content := readFile(p.patchesPath, "example.cc.patch")
	content = strings.Replace(content, "foo", "bar", 1)
	updateFile(p.patchesPath, "example.cc.patch", content)
	gitAdd(p.beaconPath, ".")
	gitCommit(p.beaconPath, "-m", "modify patch")

	added, modified, err = p.DiffSinceLastPatched()
	if err != nil || len(added) > 0 || len(modified) != 1 {
		t.Fatalf("got modified = 0, want 1")
	}
}

func updateFile(dir, name, content string) {
	check(os.WriteFile(
		filepath.Join(dir, name),
		[]byte(content),
		0644,
	))
}

func readFile(dir, name string) string {
	raw, err := os.ReadFile(filepath.Join(dir, name))
	check(err)
	return string(raw)
}

func gitAdd(repo string, args ...string) {
	c := exec.Command("git", append([]string{"add"}, args...)...)
	fmt.Println(c.String())
	c.Dir = repo
	check(c.Run())
}

func gitCommit(repo string, args ...string) {
	c := exec.Command("git", append([]string{"commit"}, args...)...)
	c.Dir = repo
	check(c.Run())
}
