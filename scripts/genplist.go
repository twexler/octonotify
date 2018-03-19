package main

import (
	"os"
	"path/filepath"

	"howett.net/plist"
)

const plistPath = "build/octonotify.app/Contents/Info.plist"

var version string

func main() {
	path, err := filepath.Abs(plistPath)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	enc := plist.NewEncoder(f)
	enc.Encode(map[string]string{
		"CFBundleName":        "octonotify",
		"CFBundleDisplayName": "octonotify",
		"CFBundleIdentifier":  "com.twexler.octonotify",
		"CFBundleVersion":     version,
		"CFBundlePackageType": "APPL",
		"CFBundleSignature":   "octo",
		"CFBundleExecutable":  "octonotify",
		"CFBundleIconFile":    "octonotify",
		"LSUIElement":         "1",
	})
}
