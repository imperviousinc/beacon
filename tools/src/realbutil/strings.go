package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/beevik/etree"
)

// Utilities for applying string replacements similar to
// Brave replacement scripts. This also supports some more advanced
// GRD replacements using beacon-add/beacon-override tags.
var (
	// paths for applying Beacon branding
	beaconifyStringPaths = []string{
		"chrome/app/chromium_strings.grd",
		"components/components_chromium_strings.grd",
		"chrome/browser/ui/android/strings/android_chrome_strings.grd",
		"components/components_strings.grd",
		"chrome/app/generated_resources.grd",
	}
)

// Some of those strings technically need new translation,
// but we don't properly support localization, yet
// they'll fall back to English.
var (
	beaconReplacers = []*strings.Replacer{
		strings.NewReplacer(
			"The Chromium Authors. All rights reserved.", "The Beacon Authors. All rights reserved.",
			"Google LLC. All rights reserved.", "The Beacon Authors. All rights reserved.",
			"The Chromium Authors", "Impervious Inc",
			"Google Chrome", "Beacon",
			"Chromium", "Beacon",
			"Chrome", "Beacon",
			"People", "Profiles",
			"You and Google", "General",
		),
		strings.NewReplacer(
			"Beacon Web Store", "Web Store",
			"Beacon Docs", "Google Docs",
			"Beacon Drive", "Google Drive",
			"Beacon Safe Browsing", "Google Safe Browsing",
			"Safe Browsing (protects you and your device from dangerous sites)",
			"Google Safe Browsing (protects you and your device from dangerous sites)",
			"Sends URLs of some pages you visit to Beacon", "Sends URLs of some pages you visit to Google",
			"Google Google", "Google",
			"BeaconOS", "ChromeOS",
			"Beacon OS", "Chrome OS",
			`Beacon is made possible by the Beacon`, `Beacon is made possible by the Chromium`,
		),
	}

	beaconRegexReplacements = map[*regexp.Regexp]string{}
)

func beaconifyText(text string) string {
	for _, beaconReplacer := range beaconReplacers {
		text = beaconReplacer.Replace(text)
	}

	for pattern, repl := range beaconRegexReplacements {
		text = pattern.ReplaceAllLiteralString(text, repl)
	}
	return text
}

func beaconifyNode(ele *etree.Element) {
	txt := ele.Text()
	if txt != "" {
		ele.SetText(beaconifyText(ele.Text()))
	}

	tail := ele.Tail()
	if tail != "" {
		ele.SetTail(beaconifyText(ele.Tail()))
	}

	desc := ele.SelectAttr("desc")
	if desc != nil && desc.Value != "" {
		desc.Value = beaconifyText(desc.Value)
	}

	for _, child := range ele.ChildElements() {
		beaconifyNode(child)
	}
}

func beaconifyGRD(grdFilePath string) (string, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(grdFilePath); err != nil {
		return "", fmt.Errorf("failed reading grd file `%s`:%v", grdFilePath, err)
	}
	for _, ele := range doc.FindElements("//message") {
		beaconifyNode(ele)
	}
	for _, ele := range doc.FindElements("//comment()") {
		beaconifyNode(ele)
	}

	out, err := doc.WriteToString()
	if err != nil {
		return "", fmt.Errorf("failed making brand replacements for `%s`:%v", grdFilePath, err)
	}

	return out, nil
}

func attrsHaveKey(arr []etree.Attr, attr etree.Attr) bool {
	for _, cur := range arr {
		if cur.Key == attr.Key {
			return true
		}
	}
	return false
}

func elementToPathQuery(e *etree.Element, includeAttrs []etree.Attr) string {
	path := []string{}
	for seg := e; seg != nil; seg = seg.Parent() {
		if seg.Tag != "" {
			pathSeg := ""
			pathSeg += seg.Tag

			if len(seg.Attr) > 0 {
				pathSeg += "["
				pathSeg += fmt.Sprintf(`@%s='%s'`, seg.Attr[0].Key, seg.Attr[0].Value)
				if len(includeAttrs) > 0 {
					for _, a := range seg.Attr[1:] {
						if attrsHaveKey(includeAttrs, a) {
							pathSeg += " and "
							pathSeg += fmt.Sprintf(`@%s='%s'`, a.Key, a.Value)
						}

					}
				}
				pathSeg += "]"
			}
			path = append(path, pathSeg)
		}
	}

	// Reverse the path.
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return "/" + strings.Join(path, "/")
}

func attrsContain(original []etree.Attr, override []etree.Attr) bool {
	if len(override) > len(original) {
		return false
	}

	for _, a := range original {
		for _, b := range override {
			if b.Key == a.Key && b.Value != a.Value {
				return false
			}
		}
	}

	return true
}

func verifyTreePaths(og *etree.Element, ov *etree.Element) bool {
	for original, override := og, ov; original != nil && override != nil; original, override = original.Parent(), override.Parent() {
		if original.Tag != override.Tag {
			return false
		}
		if !attrsContain(original.Attr, override.Attr) {
			return false
		}
	}

	return true
}

func beaconModGRD(originalGRDFile, overrideGRDFile, srcDir string) (string, error) {
	// reset first
	check(run("Reset grd file "+originalGRDFile, srcDir,
		"git", "checkout", "--", originalGRDFile))

	// apply string replacements
	beaconifyFile(originalGRDFile)

	// apply mods
	overrideGRD := etree.NewDocument()
	if err := overrideGRD.ReadFromFile(overrideGRDFile); err != nil {
		return "", err
	}
	originalGRD := etree.NewDocument()
	if err := originalGRD.ReadFromFile(originalGRDFile); err != nil {
		return "", err
	}

	addedNodes := overrideGRD.FindElements("[namespace-prefix()='beacon-add']")

	for _, addedNode := range addedNodes {
		parent := addedNode.Parent()
		if parent == nil {
			continue
		}

		pathQuery := elementToPathQuery(parent, nil)
		nodes := originalGRD.FindElements(pathQuery)
		added := false
		for _, node := range nodes {
			if !verifyTreePaths(node, parent) {
				continue
			}

			nodeCopy := addedNode.Copy()
			nodeCopy.Space = ""
			node.AddChild(nodeCopy)
			added = true
		}

		if !added {
			return "", fmt.Errorf("unable to add node no matching path for `%s`", pathQuery)
		}
	}

	replacementNodes := overrideGRD.FindElements("//[namespace-prefix()='beacon-override']")

	for _, replacementNode := range replacementNodes {
		pathQuery := elementToPathQuery(replacementNode, nil)
		nodes := originalGRD.FindElements(pathQuery)
		overridden := false
		for _, node := range nodes {
			if !verifyTreePaths(node, replacementNode) {
				continue
			}

			idx := node.Index()
			parent := node.Parent()
			parent.RemoveChildAt(idx)
			replacement := replacementNode.Copy()
			replacement.Space = ""
			parent.InsertChildAt(idx, replacement)
			overridden = true
		}
		if !overridden {
			return "", fmt.Errorf("unable to override node no matching path for `%s`", pathQuery)
		}
	}

	return originalGRD.WriteToString()
}

func beaconifyFile(grdFilePath string) {
	out, err := beaconifyGRD(grdFilePath)
	if err != nil {
		log.Printf("[WARNING] %v", err)
		return
	}

	if err := ioutil.WriteFile(grdFilePath, []byte(out), 0644); err != nil {
		log.Printf("[WARNING] failed storing brand replacements for `%s`:%v", grdFilePath, err)
	}
}

func beaconifyChromium(chromiumSrcDir string) {
	var beaconifyPaths []string
	for _, grdPath := range beaconifyStringPaths {
		grdPath = filepath.Join(chromiumSrcDir, grdPath)
		if !fileExists(grdPath) {
			log.Printf("[WARNING] grd path `%s` does not exist", grdPath)
			continue
		}

		beaconifyPaths = append(beaconifyPaths, recurseGrdNoMapping(grdPath, []string{})...)
	}

	for _, bp := range beaconifyPaths {
		beaconifyFile(bp)
	}
}

func getGRDPartsFromGRDP(path string) ([]string, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		return nil, fmt.Errorf("failed reading grd file `%s`:%v", path, err)
	}

	var parts []string
	paths := doc.FindElements("//part")
	for _, p := range paths {
		file := p.SelectAttrValue("file", "")
		if file != "" {
			parts = append(parts, file)
		}
	}

	return parts, nil
}

func recurseGrdNoMapping(chromiumPath string, exclude []string) []string {
	mustExist(chromiumPath)
	var files []string

	grdps, err := getGRDPartsFromGRDP(chromiumPath)
	if err != nil {
		log.Fatal(err)
	}

	chromiumDir := filepath.Dir(chromiumPath)

	files = append(files, chromiumPath)
	if len(grdps) == 0 {
		return files
	}

	for _, grdp := range grdps {
		if hasString(exclude, grdp) {
			continue
		}

		chromiumGRDPPath := filepath.Join(chromiumDir, grdp)

		files2 := recurseGrdNoMapping(chromiumGRDPPath, exclude)
		files = append(files, files2...)
	}

	return files
}
