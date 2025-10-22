package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// GetGitBranch retrieves the git branch of the current directory
func GetGitBranch(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error running git command: %w, stderr: %s", err, errBuf.String())
	}

	branch := strings.TrimSpace(out.String())
	if branch == "HEAD" {
		// means detached HEAD
		return "", fmt.Errorf("detached HEAD or no branch (got: %q)", branch)
	}

	return branch, nil
}
