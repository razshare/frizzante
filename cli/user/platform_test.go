package user

import (
	"testing"

	"github.com/razshare/frizzante/cli/apps"
	"github.com/razshare/frizzante/platform"
)

func TestPlatformLinuxAmd64(t *testing.T) {
	platStr := "linux/amd64"

	plat, err := Platform(&apps.App{Platform: &platStr})
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.LinuxAmd64 {
		t.Fatal("platform should be linux amd64")
	}
}

func TestPlatformLinuxArm64(t *testing.T) {
	platStr := "linux/arm64"

	plat, err := Platform(&apps.App{Platform: &platStr})
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.LinuxArm64 {
		t.Fatal("platform should be linux arm64")
	}
}

func TestPlatformDarwinAmd64(t *testing.T) {
	platStr := "darwin/amd64"

	plat, err := Platform(&apps.App{Platform: &platStr})
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.DarwinAmd64 {
		t.Fatal("platform should be darwin amd64")
	}
}

func TestPlatformDarwinArm64(t *testing.T) {
	platStr := "darwin/arm64"

	plat, err := Platform(&apps.App{Platform: &platStr})
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.DarwinArm64 {
		t.Fatal("platform should be darwin arm64")
	}
}

func TestPlatformWindowsAmd64(t *testing.T) {
	platStr := "windows/amd64"

	plat, err := Platform(&apps.App{Platform: &platStr})
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.WindowsAmd64 {
		t.Fatal("platform should be windows amd64")
	}
}

func TestPlatformWindowsArm64(t *testing.T) {
	platStr := "windows/arm64"

	plat, err := Platform(&apps.App{Platform: &platStr})
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.WindowsArm64 {
		t.Fatal("platform should be windows arm64")
	}
}
