package action

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/cli/npm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/npmselect"
)

func Npm(o NpmOptions) error {
	selectedPackages, err := npmselect.Send()
	if err != nil {
		if err.Error() == "cancelled" {
			messages.Info("npm search cancelled")
			return nil
		}
		return fmt.Errorf("failed to run npm search: %w", err)
	}

	if len(selectedPackages) == 0 {
		messages.Info("No packages selected")
		return nil
	}

	bunPath := filepath.Join(".gen", "bun", "bun")
	if _, err := os.Stat(bunPath); os.IsNotExist(err) {
		bunPath = "bun"
	}

	return npm.InstallNpmPackages(selectedPackages, bunPath)
}