package user

import (
	"github.com/razshare/frizzante/cli/app"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/platform"
	"os"
	"path/filepath"
	"testing"
)

func TestPlatformLinuxAmd64(t *testing.T) {
	home, err := FrizzanteHome()
	if err != nil {
		t.Fatal(err)
	}

	_ = os.Remove(filepath.Join(home, "platform.txt"))
	defer func() { _ = os.Remove(filepath.Join(home, "platform.txt")) }()

	plats := "linux/amd64"
	a := &app.App{Platform: &plats}

	plat, err := Platform(a)
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.LinuxAmd64 {
		t.Fatal("platform should be linux amd64")
	}

	if !files.IsFile(filepath.Join(home, "platform.txt")) {
		t.Fatal("~/platform.txt should be a file")
	}

	d, err := os.ReadFile(filepath.Join(home, "platform.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != "linux/amd64" {
		t.Fatal("~/platform.txt should container linux/amd64")
	}
}

func TestPlatformLinuxArm64(t *testing.T) {
	home, err := FrizzanteHome()
	if err != nil {
		t.Fatal(err)
	}

	_ = os.Remove(filepath.Join(home, "platform.txt"))
	defer func() { _ = os.Remove(filepath.Join(home, "platform.txt")) }()

	plats := "linux/arm64"
	a := &app.App{Platform: &plats}

	plat, err := Platform(a)
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.LinuxArm64 {
		t.Fatal("platform should be linux arm64")
	}

	if !files.IsFile(filepath.Join(home, "platform.txt")) {
		t.Fatal("~/platform.txt should be a file")
	}

	d, err := os.ReadFile(filepath.Join(home, "platform.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != "linux/arm64" {
		t.Fatal("~/platform.txt should container linux/arm64")
	}
}

func TestPlatformDarwinAmd64(t *testing.T) {
	home, err := FrizzanteHome()
	if err != nil {
		t.Fatal(err)
	}

	_ = os.Remove(filepath.Join(home, "platform.txt"))
	defer func() { _ = os.Remove(filepath.Join(home, "platform.txt")) }()

	plats := "darwin/amd64"
	a := &app.App{Platform: &plats}

	plat, err := Platform(a)
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.DarwinAmd64 {
		t.Fatal("platform should be darwin amd64")
	}

	if !files.IsFile(filepath.Join(home, "platform.txt")) {
		t.Fatal("~/platform.txt should be a file")
	}

	d, err := os.ReadFile(filepath.Join(home, "platform.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != "darwin/amd64" {
		t.Fatal("~/platform.txt should container darwin/amd64")
	}
}

func TestPlatformDarwinArm64(t *testing.T) {
	home, err := FrizzanteHome()
	if err != nil {
		t.Fatal(err)
	}

	_ = os.Remove(filepath.Join(home, "platform.txt"))
	defer func() { _ = os.Remove(filepath.Join(home, "platform.txt")) }()

	plats := "darwin/arm64"
	a := &app.App{Platform: &plats}

	plat, err := Platform(a)
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.DarwinArm64 {
		t.Fatal("platform should be darwin arm64")
	}

	if !files.IsFile(filepath.Join(home, "platform.txt")) {
		t.Fatal("~/platform.txt should be a file")
	}

	d, err := os.ReadFile(filepath.Join(home, "platform.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != "darwin/arm64" {
		t.Fatal("~/platform.txt should container darwin/arm64")
	}
}

func TestPlatformWindowsAmd64(t *testing.T) {
	home, err := FrizzanteHome()
	if err != nil {
		t.Fatal(err)
	}

	_ = os.Remove(filepath.Join(home, "platform.txt"))
	defer func() { _ = os.Remove(filepath.Join(home, "platform.txt")) }()

	plats := "windows/amd64"
	a := &app.App{Platform: &plats}

	plat, err := Platform(a)
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.WindowsAmd64 {
		t.Fatal("platform should be windows amd64")
	}

	if !files.IsFile(filepath.Join(home, "platform.txt")) {
		t.Fatal("~/platform.txt should be a file")
	}

	d, err := os.ReadFile(filepath.Join(home, "platform.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != "windows/amd64" {
		t.Fatal("~/platform.txt should container windows/amd64")
	}
}

func TestPlatformWindowsArm64(t *testing.T) {
	home, err := FrizzanteHome()
	if err != nil {
		t.Fatal(err)
	}

	_ = os.Remove(filepath.Join(home, "platform.txt"))
	defer func() { _ = os.Remove(filepath.Join(home, "platform.txt")) }()

	plats := "windows/arm64"
	a := &app.App{Platform: &plats}

	plat, err := Platform(a)
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.WindowsArm64 {
		t.Fatal("platform should be windows arm64")
	}

	if !files.IsFile(filepath.Join(home, "platform.txt")) {
		t.Fatal("~/platform.txt should be a file")
	}

	d, err := os.ReadFile(filepath.Join(home, "platform.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != "windows/arm64" {
		t.Fatal("~/platform.txt should container windows/arm64")
	}
}
