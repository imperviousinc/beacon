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

const ChromiumTag = "101.0.4951.54"

type Channel int

const (
	DevelopmentChannel Channel = 0
	ReleaseChannel     Channel = 1
)

type BuildArgs struct {
	args map[string]interface{}
}

func initBuildDir(workingDir string) {
	if err := generateGClientFile(workingDir); err != nil {
		check(fmt.Errorf("failed creating .glicent file: %v", err))
	}

	// pull chromium's code
	check(run("Running gclient sync", workingDir,
		"gclient", "sync", "--revision", "src@refs/tags/"+ChromiumTag))

	// ensure correct tag
	c := exec.Command("git", "describe", "--tags")
	c.Dir = filepath.Join(workingDir, "src")
	out, err := c.Output()
	if err != nil {
		check(fmt.Errorf("failed checking tag: %v", err))
	}

	actualTag := strings.TrimSpace(string(out))
	if actualTag != ChromiumTag {
		check(fmt.Errorf("got chromium tag = %s, want tag = %s", actualTag, ChromiumTag))
	}
}

func generateGClientFile(initDir string) error {
	gclientFile := `solutions = [
  {
    "name": "src",
    "url": "https://chromium.googlesource.com/chromium/src.git",
    "managed": False,
    "custom_deps": {},
    "custom_vars": {
		"checkout_pgo_profiles": True,
	},
  },
]`

	return writeFile(filepath.Join(initDir, ".gclient"), []byte(gclientFile))
}

func createBuildArgs(srcDir string, ch Channel) map[string]interface{} {

	args := make(map[string]interface{})
	args["branding_path_component"] = "beacon"

	credirect := "credirect"
	if runtime.GOOS == "windows" {
		credirect += ".exe"
	}

	ccWrapper := filepath.Join(srcDir, "beacon", "tools", "bin", credirect)
	mustExist(ccWrapper)
	args["cc_wrapper"] = ccWrapper
	args["root_extra_deps"] = []string{"//beacon"}

	// warn if no cc wrapper set to use ccache
	// for faster builds
	getEnvOrWarn("BEACON_CC_WRAPPER")

	args["chrome_pgo_phase"] = 0
	args["is_component_build"] = true
	args["is_debug"] = true
	args["enable_nacl"] = false

	if ch == ReleaseChannel {
		delete(args, "is_component_build")
		delete(args, "chrome_pgo_phase")
		args["is_official_build"] = true
		args["enable_widevine"] = true
		args["proprietary_codecs"] = true
		args["is_debug"] = false
		args["ffmpeg_branding"] = "Chrome"
		args["disable_fieldtrial_testing_config"] = true
		args["enable_hangout_services_extension"] = true
		args["enable_pseudolocales"] = false
		args["safe_browsing_mode"] = 1
		args["google_api_key"] = getEnvOrWarn("BEACON_GOOGLE_API_KEY")
		args["google_default_client_id"] = getEnvOrWarn("BEACON_GOOGLE_DEFAULT_CLIENT_ID")
		args["google_default_client_secret"] = getEnvOrWarn("BEACON_GOOGLE_DEFAULT_CLIENT_SECRET")
	}

	return args
}

func getEnvOrWarn(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("[WARNING] environment variable `%s` is not set.", key)
	}
	return val
}

func genBuild(srcDir string, ch Channel, buildDir, target string) error {
	argsPath := filepath.Join(srcDir, buildDir, "args.gn")

	if fileExists(argsPath) {
		return nil
	}

	args := createBuildArgs(srcDir, ch)
	if target != "" {
		args["target_cpu"] = target
	}

	var b strings.Builder
	for key, value := range args {
		switch value.(type) {
		case int:
			b.WriteString(fmt.Sprintf("%s=%v\n", key, value))
		case bool:
			b.WriteString(fmt.Sprintf("%s=%v\n", key, value))
		case []string:
			value := value.([]string)
			b.WriteString(fmt.Sprintf("%s=[ ", key))
			b.WriteString(fmt.Sprintf("%q", value[0]))
			for _, item := range value[1:] {
				b.WriteRune(',')
				b.WriteString(fmt.Sprintf("%q, ", item))
			}
			b.WriteString(" ]\n")
		default:
			b.WriteString(fmt.Sprintf("%s=%q\n", key, value))
		}
	}

	if err := writeFile(argsPath, []byte(b.String())); err != nil {
		return err
	}

	return run("Setting up the build", srcDir,
		"gn", "gen", buildDir)
}
