package main

import (
	"github.com/adrg/xdg"
	"github.com/samcofer/tam-offline-download-email/cmd"
	"os"
	"path/filepath"
	"runtime"
)

var (
	version = "dev"
)

func main() {
	// normalize the config home on osx to linux and get rid of the path spacing
	if runtime.GOOS == "darwin" {
		hd, _ := os.UserHomeDir()
		os.Setenv("XDG_CONFIG_HOME", filepath.Join(hd, ".config"))
		os.Setenv("XDG_DATA_HOME", filepath.Join(hd, ".local", "share"))
		xdg.Reload()
	}

	cmd.Execute(version, os.Args[1:])
}
