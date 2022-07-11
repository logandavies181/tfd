package git

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// GetRootGetRootOfRepo determines the path to the root of the git repo that "path" is under and returns it and the
// relative path from the root to "path"
func GetRootOfRepo(path string) (string, string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	// TODO: sanity check that this is a dir
	cmd.Dir = path

	pathToRootBytes, err := cmd.Output()
	if err != nil {
		return "", "", err
	}
	pathToRoot := strings.TrimSpace(string(pathToRootBytes))

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", "", err
	}

	relDir, err := filepath.Rel(pathToRoot, absPath)
	if err != nil {
		return "", "", err
	}

	return pathToRoot, relDir, nil
}
