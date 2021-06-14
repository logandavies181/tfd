package git

import (
	"os/exec"
	"strings"
)

func GetRootOfRepo(path string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	// TODO: sanity check that this is a dir
	cmd.Dir = path

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
