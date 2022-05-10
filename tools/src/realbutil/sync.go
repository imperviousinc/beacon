package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var errSyncRequiresCleanWorkingTree = errors.New("sync failed: you have unstaged changes")

// verifySynced verifies that ChromiumTag == git describe --tags
func verifySynced(chromePath string) error {
	c := exec.Command("git", "describe", "--tags")
	c.Dir = chromePath
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed checking tags: %v:%s", err, string(out))
	}

	tag := strings.TrimSpace(string(out))
	if tag != ChromiumTag {
		return fmt.Errorf("current version %s != supported version %s", tag, ChromiumTag)
	}

	return nil
}

// sync syncs chromium to the latest supported tag
// It will return an error if the chrome repo working tree isn't clean
//	or when some action fails.
// Otherwise, it will perform the following actions:
//  1. Switch to the newly supported tag in build.go
//     If the tag doesn't exist, it should try git pull
//  2. call `gclient sync` to update dependencies.
//  3. Re-apply all patches and overrides
func sync(chromePath string) error {
	if err := verifyWorkingTreeClean(chromePath); err != nil {
		return err
	}
	if err := run("Fetching tags", chromePath, "git", "fetch", "--tags"); err != nil {
		return fmt.Errorf("fetching tags failed: %v", err)
	}
	if err := run("Checkout tag "+ChromiumTag, chromePath,
		"git", "checkout", "tags/"+ChromiumTag); err != nil {
		return fmt.Errorf("checkout tag failed: %v", err)
	}
	if err := run("gclient sync --revision src@refs/tags/"+ChromiumTag, chromePath,
		"gclient", "sync", "--revision", "src@refs/tags/"+ChromiumTag); err != nil {
		return fmt.Errorf("sync failed: %v", err)
	}
	return nil
}

// verifyWorkingTreeClean verifies that the git working tree
// is clean based on https://github.com/git/git/blob/master/git-sh-setup.sh#L202
func verifyWorkingTreeClean(repo string) error {
	c := exec.Command("git", "status", "--short")
	c.Dir = repo
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git status failed: %v: %s", err, string(out))
	}

	changes := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(changes) != 1 {
		return errSyncRequiresCleanWorkingTree
	}
	if !strings.Contains(changes[0], "beacon/") {
		return errSyncRequiresCleanWorkingTree
	}
	return nil
}
