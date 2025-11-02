package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ensureCacheDir() (string, error) {
	cacheDir, _ := os.UserHomeDir()
	cacheDir = filepath.Join(cacheDir, ".ecac", "plugins")
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return "", err
	}
	return cacheDir, nil
}

func FetchAndBuild(source, version string) (string, error) {
	cacheDir, err := ensureCacheDir()
	if err != nil {
		return "", err
	}

	repoName := filepath.Base(source)
	versionDir := filepath.Join(cacheDir, repoName, version)
	binPath := filepath.Join(versionDir, "plugin-bin")

	// Skip if already built
	if _, err := os.Stat(binPath); err == nil {
		return binPath, nil
	}

	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("plugin-%s", repoName))
	_ = os.RemoveAll(tmp)

	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", version, "https://"+source+".git", tmp)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git clone failed: %w", err)
	}

	_ = os.MkdirAll(versionDir, 0755)
	build := exec.Command("go", "build", "-o", binPath, ".")
	build.Dir = tmp
	if out, err := build.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build failed: %v\n%s", err, string(out))
	}

	return binPath, nil
}
