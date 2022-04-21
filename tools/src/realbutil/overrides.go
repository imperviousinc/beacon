package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var beaconResourcesOverrides = []string{
	// We still need to override some product logos since these paths are defined in
	// chrome_unscaled_resources.grd with expressions like <if expr="not _google_chrome">
	// it's easier to override than to modify the grd file.
	"app/theme/beacon", "chrome/app/theme/chromium",
	"app/theme/beacon", "chrome/app/theme/beacon",

	// There's more defined in theme_resources.grd
	"app/theme/default_100_percent/beacon", "chrome/app/theme/default_100_percent/chromium",
	"app/theme/default_200_percent/beacon", "chrome/app/theme/default_200_percent/chromium",
	"app/theme/default_100_percent/beacon", "chrome/app/theme/default_100_percent/beacon",
	"app/theme/default_200_percent/beacon", "chrome/app/theme/default_200_percent/beacon",

	// And some in `common` directory
	"app/theme/default_100_percent/common", "chrome/app/theme/default_100_percent/common",
	"app/theme/default_200_percent/common", "chrome/app/theme/default_200_percent/common",

	// We can't just override chromium_behaviors.cc using chromium_src
	// since chrome looks for ${branding_path_component}_behaviors.cc
	"chromium_src/chrome/installer/setup/beacon_behaviors.cc", "chrome/installer/setup/beacon_behaviors.cc",

	// anything placed in this dir will get copied
	// to chromium
	"overrides", "",
}

func beaconModGRDAll(modSrcDir, chromiumSrcDir string, dry bool) error {
	walkFunc := func(grdOverridePath string, overrideInfo fs.FileInfo, err error) error {
		if err != nil || overrideInfo.IsDir() {
			return nil
		}
		if !strings.HasSuffix(grdOverridePath, "__override.grd") &&
			!strings.HasSuffix(grdOverridePath, "__override.grdp") {
			return nil
		}

		relPath, err := filepath.Rel(modSrcDir, grdOverridePath)
		if err != nil {
			return fmt.Errorf("failed getting relative path from beacon src: %v", err)
		}

		relPath = strings.Replace(relPath, "__override.grd", ".grd", 1)
		chromeGRDOutputPath := filepath.Join(chromiumSrcDir, relPath)
		chromeGRDInfo, err := os.Stat(chromeGRDOutputPath)
		if err != nil {
			return fmt.Errorf("failed reading chromium's grd file `%s`: %v", grdOverridePath, err)
		}
		// if mod times are equal skip
		if overrideInfo.ModTime().Equal(chromeGRDInfo.ModTime()) {
			return nil
		}

		out, err := beaconModGRD(chromeGRDOutputPath, grdOverridePath, chromiumSrcDir)
		if err != nil {
			return fmt.Errorf("failed modding grd file")
		}

		if dry {
			log.Printf("Will mod GRD file `%s`", chromeGRDOutputPath)
			return nil
		}

		if err := writeFile(chromeGRDOutputPath, []byte(out)); err != nil {
			return fmt.Errorf("failed updating grd file `%s`:%v", chromeGRDOutputPath, err)
		}

		if err = compareFixOverrideTime(chromeGRDOutputPath, grdOverridePath,
			chromeGRDInfo.ModTime(), overrideInfo.ModTime()); err != nil {
			return err
		}

		return nil
	}

	return filepath.Walk(modSrcDir, walkFunc)
}

// copyResource copies a src file or directory recursively to dst
// it will override any files that exist in dst, and will create
// any paths that don't exist in dst.
func copyResource(src, dst string, dry bool) error {
	finfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed reading source file: %w", err)
	}
	walkFunc := func(path string, overrideInfo fs.FileInfo, err error) error {
		if err != nil || overrideInfo.IsDir() {
			return nil
		}

		// skip grd overrides
		if strings.HasSuffix(path, "__override.grd") || strings.HasSuffix(path, "__override.grdp") {
			return nil
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("failed getting relative path from beacon src: %v", err)
		}

		chromePath := filepath.Join(dst, relPath)
		originalInfo, _ := os.Stat(chromePath)
		// if src and dst files have same mod time, skip.
		if originalInfo != nil && originalInfo.ModTime().Equal(overrideInfo.ModTime()) {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed reading file from beacon path `%s`: %v", path, err)
		}

		if dry {
			if originalInfo != nil {
				log.Printf("Will override to `%s`", chromePath)
			} else {
				log.Printf("Will copy to `%s`", chromePath)
			}
			return nil
		}

		if err := writeFile(chromePath, data); err != nil {
			return fmt.Errorf("failed copying beacon src to dst file `%s`:%v", chromePath, err)
		}

		if originalInfo != nil {
			// writeFile changes mod time so we need to fix it
			return fixOverrideTime(chromePath, path)
		}
		return nil
	}

	if !finfo.IsDir() {
		return walkFunc(src, finfo, nil)
	}

	return filepath.Walk(src, walkFunc)
}

// copyBeaconResources src and dst can both either be a directory or
// or a file. Copies and overrides  all files found in src to dst
func copyBeaconResources(beaconRootDir, chromeRootDir string, dry bool) error {
	if len(beaconResourcesOverrides)%2 != 0 {
		panic("beaconResourcesOverrides must have pairs of src -> dst")
	}

	for i := 0; i < len(beaconResourcesOverrides); i += 2 {
		src := beaconResourcesOverrides[i]
		dst := beaconResourcesOverrides[i+1]
		beaconSrc := filepath.Join(beaconRootDir, src)
		chromeDst := filepath.Join(chromeRootDir, dst)

		if err := copyResource(beaconSrc, chromeDst, dry); err != nil {
			return fmt.Errorf("override failed: %v", err)
		}
	}

	return nil
}
