package npm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func InstallNpmPackages(packages []string, bun string) error {
	if len(packages) == 0 {
		return nil
	}

	appDir := "app"
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		return fmt.Errorf("app directory does not exist")
	}

	packageJsonPath := filepath.Join(appDir, "package.json")
	if _, err := os.Stat(packageJsonPath); os.IsNotExist(err) {
		return fmt.Errorf("package.json not found in app directory")
	}

	successCount := 0
	for _, pkgName := range packages {
		s := spinner.New(fmt.Sprintf("installing %s to app/node_modules", pkgName))
		go spinner.Start(s)

		// Install the package using bun add
		cmd := exec.Command(bun, "add", pkgName)
		cmd.Dir = appDir
		cmd.Env = os.Environ()

		output, err := cmd.CombinedOutput()
		spinner.Stop(s)

		if err != nil {
			messages.Error(fmt.Sprintf("Failed to install %s: %v\nOutput: %s", pkgName, err, string(output)))
			continue
		}

		messages.Success(fmt.Sprintf("Installed %s to app/node_modules", pkgName))
		successCount++
	}

	if successCount > 0 {
		messages.Success(fmt.Sprintf("Successfully installed %d package(s) to app/node_modules", successCount))
	}

	return nil
}