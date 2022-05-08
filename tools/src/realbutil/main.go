package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const uninitializedRepo = ".beacon__repo"

func fetchChromium() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed getting working dir: %v", workDir)
	}

	if !fileExists(uninitializedRepo) {
		return
	}

	srcDir := filepath.Join(workDir, "src")
	// create checkout .gclient dir and sync
	initBuildDir(workDir)
	mustExist("src")
	mustExist(".gclient")

	beaconDir := filepath.Join(srcDir, "beacon")
	if !fileExists(beaconDir) {
		// dummy butil should've cloned the repo here
		repoPath := filepath.Join(workDir, uninitializedRepo)
		mustExist(repoPath)
		check(os.Rename(repoPath, beaconDir))
		mustExist(beaconDir)
	}

	// used to mark this path as the root project dir
	check(ioutil.WriteFile(filepath.Join(workDir, ".butil"), []byte{}, 0644))
}

func handleInitCmd(c *cli.Context) error {
	fetchChromium()

	butilPath, err := os.Executable()
	check(err)
	workDir, err := os.Getwd()
	check(err)

	if err := run("Applying patches", workDir, butilPath, "patches", "apply"); err != nil {
		return fmt.Errorf("run failed: fix and call `butil init` again: %v", err)
	}
	return run("Updating strings", workDir, butilPath, "strings", "rebase")
}

// mustBeInRootAndInitialized verifies the butil is being called
// from the root project directory, and that .butil file is present.
// returns the project's root directory
func mustBeInRootAndInitialized() string {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed getting working dir: %v", workDir)
	}

	butilPath := filepath.Join(workDir, ".butil")
	if !fileExists(butilPath) {
		// try alternate path maybe its called from src/
		butilPath = filepath.Join(workDir, "..", ".butil")
		if !fileExists(butilPath) {
			log.Fatalf("This command check be called from the root project directory (parent of src/).")
		}

		workDir = filepath.Dir(butilPath)
		mustExist(filepath.Join(workDir, ".butil"))
	}

	mustExist(filepath.Join(workDir, "src"))

	return workDir
}

func prepareBuild(c *cli.Context) (string, error) {
	workDir := mustBeInRootAndInitialized()
	srcDir := filepath.Join(workDir, "src")
	mustExist(srcDir)
	beaconChromiumSrcDir := filepath.Join(srcDir, "beacon", "chromium_src")
	mustExist(beaconChromiumSrcDir)

	// First touch any changed .h, .cc or .mm files
	touchCCOverrides(srcDir, beaconChromiumSrcDir)

	// Copy all overridden resources in overrides ... etc.
	if err := applyResourceOverridesCmd(c); err != nil {
		return "", err
	}
	appendToPolymerBundle(srcDir)

	// Warn if patches were modified
	p, err := LoadPatcher(workDir)
	if err != nil {
		return "", err
	}

	reapply, err := p.ShouldReapply()
	if err != nil {
		return "", err
	}
	if reapply && !c.Bool("force") {
		fmt.Println("[ERROR] some patches were modified since they were last applied")
		fmt.Println("[ERROR] consider re-applying patches or force build with -f option")
		return "", fmt.Errorf("patches check failed")
	}

	// Apply grd modding (this doesn't do the string replacements called at init)
	err = beaconModGRDAll(filepath.Join(srcDir, "beacon", "overrides"), srcDir, false)
	return srcDir, err
}

func handleBuildDebugCmd(c *cli.Context) error {
	srcDir, err := prepareBuild(c)
	if err != nil {
		return err
	}
	target := c.String("target")
	suffix := ""
	if target != "" {
		suffix += "_" + target
	}

	buildDir := "out/Debug" + suffix

	if err := genBuild(srcDir, DevelopmentChannel, buildDir, target); err != nil {
		return err
	}

	return run("Building Beacon", srcDir,
		"autoninja", "-C", buildDir, "beacon")
}

func handleBuildReleaseCmd(c *cli.Context) error {
	srcDir, err := prepareBuild(c)
	if err != nil {
		return err
	}
	target := c.String("target")

	suffix := ""
	if target != "" {
		suffix += "_" + target
	}
	buildDir := "out/Release" + suffix

	if err := genBuild(srcDir, ReleaseChannel, buildDir, target); err != nil {
		return err
	}

	return run("Building Beacon", srcDir,
		"autoninja", "-C", buildDir, "beacon")
}

func updatePatchesCmd(c *cli.Context) error {
	workDir := mustBeInRootAndInitialized()
	p, err := LoadPatcher(workDir)
	if err != nil {
		return err
	}

	return p.Update()
}

func applyPatchedCmd(c *cli.Context) error {
	workDir := mustBeInRootAndInitialized()
	p, err := LoadPatcher(workDir)
	if err != nil {
		return err
	}
	return p.Apply()
}

func resetStringsCmd(c *cli.Context) error {
	workDir := mustBeInRootAndInitialized()

	srcDir := filepath.Join(workDir, "src")
	mustExist(srcDir)

	check(run("Resetting GRD Files", srcDir,
		"git", "checkout", "--", "*.grd"))

	check(run("Resetting GRDP Files", srcDir,
		"git", "checkout", "--", "*.grdp"))

	return nil
}

func applyStringsCmd(c *cli.Context) error {
	workDir := mustBeInRootAndInitialized()

	srcDir := filepath.Join(workDir, "src")
	mustExist(srcDir)

	beaconifyChromium(srcDir)
	return beaconModGRDAll(filepath.Join(srcDir, "beacon", "overrides"), srcDir, false)
}

func applyResourceOverridesCmd(c *cli.Context) error {
	rootDir := mustBeInRootAndInitialized()

	srcDir := filepath.Join(rootDir, "src")
	mustExist(srcDir)

	beaconDir := filepath.Join(srcDir, "beacon")
	mustExist(beaconDir)

	return copyBeaconResources(beaconDir, srcDir, c.Bool("dry"))
}

func rebaseStringsCmd(c *cli.Context) error {
	if err := resetStringsCmd(c); err != nil {
		return err
	}

	if err := applyStringsCmd(c); err != nil {
		return err
	}

	return nil
}

// buildTools finds all go projects in tools/src/* directory
// and builds them. Output goes to tools/bin/*
func buildToolsCmd(c *cli.Context) error {
	workDir := mustBeInRootAndInitialized()
	force := c.Bool("force")

	toolsDir := filepath.Join(workDir, "src", "beacon", "tools")
	toolsSrc := filepath.Join(toolsDir, "src")
	mustExist(toolsSrc)

	toolsBin := filepath.Join(toolsDir, "bin")
	os.MkdirAll(toolsBin, 0755)
	mustExist(toolsBin)

	// find all projects in tools/src
	projects, err := ioutil.ReadDir(toolsSrc)
	if err != nil {
		log.Fatalf("failed listing tools dir `%s`: %v", toolsDir, err)
	}

	for _, proj := range projects {
		toolSrc := filepath.Join(toolsSrc, proj.Name())
		toolBin := filepath.Join(toolsBin, proj.Name())
		if runtime.GOOS == "windows" {
			toolBin += ".exe"
		}
		if fileExists(toolBin) && !force {
			continue
		}

		check(run("Building tools/"+proj.Name(), toolSrc,
			"go", "build", "-o", toolBin, toolSrc))

	}
	return nil
}

func main() {
	app := &cli.App{
		Name:  "butil",
		Usage: "Beacon browser development utility",
	}
	app.Commands = []*cli.Command{
		{
			Name:  "init",
			Usage: "Initializes Beacon by fetching repos, applying patches, preparing build tools ... etc.",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "beacon-repo",
					Aliases: []string{"b"},
					Usage:   "Beacon browser git repo to be cloned at src/beacon",
					Value:   "https://github.com/imperviousinc/beacon",
				},
			},
			Action: handleInitCmd,
		},
		{
			Name:  "build",
			Usage: "Builds Beacon",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "target",
					Usage: "Specify build target",
					Value: "",
				},
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "Force build ignoring errors",
					Value:   false,
				},
			},
			Subcommands: []*cli.Command{
				{
					Name:   "debug",
					Usage:  "Creates a debug/component build",
					Action: handleBuildDebugCmd,
				},
				{
					Name:   "release",
					Usage:  "Creates a release build",
					Action: handleBuildReleaseCmd,
				},
			},
		},
		{
			Name:   "build-tools",
			Usage:  "Builds all tools in beacon/tools/ directory",
			Action: buildToolsCmd,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "Force re-compiling all build tools",
					Value:   false,
				},
			},
		},
		{
			Name:  "sync",
			Usage: "Updates Chromium and re-applies patches",
		},
		{
			Name:  "patches",
			Usage: "Manage patches to Chromium repo",
			Subcommands: []*cli.Command{
				{
					Name:   "apply",
					Usage:  "Applies patches",
					Action: applyPatchedCmd,
				},
				{
					Name:   "update",
					Usage:  "Reads all modifications and updates patches",
					Action: updatePatchesCmd,
				},
			},
		},
		{
			Name:  "strings",
			Usage: "Manage modifications to string resources",
			Subcommands: []*cli.Command{
				{
					Name:   "reset",
					Usage:  "Clears all strings replacements",
					Action: resetStringsCmd,
				},
				{
					Name:   "apply",
					Usage:  "Applies strings replacements",
					Action: applyStringsCmd,
				},
				{
					Name:   "rebase",
					Usage:  "Resets and re-applies strings replacements",
					Action: rebaseStringsCmd,
				},
			},
		},
		{
			Name:  "overrides",
			Usage: "Manage Beacon files copied to Chromium",
			Subcommands: []*cli.Command{
				{
					Name:   "apply",
					Usage:  "Applies resource overrides",
					Action: applyResourceOverridesCmd,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:    "dry",
							Aliases: []string{"d"},
							Usage:   "A dry run only showing files that will be overridden",
							Value:   false,
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
