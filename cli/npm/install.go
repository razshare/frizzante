package npm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func InstallNpmPackages(packages []string, bun string, appDir string) error {
	if len(packages) == 0 {
		return nil
	}

	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		return fmt.Errorf("%s directory does not exist", appDir)
	}

	packageJsonPath := filepath.Join(appDir, "package.json")
	if _, err := os.Stat(packageJsonPath); os.IsNotExist(err) {
		return fmt.Errorf("package.json not found in %s directory", appDir)
	}

	successCount := 0
	for _, pkgName := range packages {
		s := spinner.New(fmt.Sprintf("installing %s to %s/node_modules", pkgName, appDir))
		go spinner.Start(s)

		// Install the package using bun add
		cmd := exec.Command(bun, "add", pkgName)
		cmd.Dir = appDir
		cmd.Env = os.Environ()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		spinner.Stop(s)

		if err != nil {
			messages.Error(fmt.Sprintf("Failed to install %s: %v", pkgName, err))
			continue
		}

		messages.Success(fmt.Sprintf("Installed %s to %s/node_modules", pkgName, appDir))
		successCount++
	}

	if successCount > 0 {
		messages.Success(fmt.Sprintf("Successfully installed %d package(s) to %s/node_modules", successCount, appDir))
	}

	return nil
}