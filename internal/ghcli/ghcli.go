// copied from, and should be updated in sync with, gh cli:
// https://github.com/cli/cli/blob/0607ce56c5eb7edd1a50872fb364a11b52930c39/internal/config/config_file.go
package ghcli

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	GH_CONFIG_DIR   = "GH_CONFIG_DIR"
	XDG_CONFIG_HOME = "XDG_CONFIG_HOME"
	XDG_STATE_HOME  = "XDG_STATE_HOME"
	XDG_DATA_HOME   = "XDG_DATA_HOME"
	APP_DATA        = "AppData"
	LOCAL_APP_DATA  = "LocalAppData"
)

// Config path precedence
// 1. GH_CONFIG_DIR
// 2. XDG_CONFIG_HOME
// 3. AppData (windows only)
// 4. HOME
func GetGHCliConfigPath() string {
	var path string
	if a := os.Getenv(GH_CONFIG_DIR); a != "" {
		path = a
	} else if b := os.Getenv(XDG_CONFIG_HOME); b != "" {
		path = filepath.Join(b, "gh")
	} else if c := os.Getenv(APP_DATA); runtime.GOOS == "windows" && c != "" {
		path = filepath.Join(c, "GitHub CLI")
	} else {
		d, _ := os.UserHomeDir()
		path = filepath.Join(d, ".config", "gh")
	}

	return path
}
