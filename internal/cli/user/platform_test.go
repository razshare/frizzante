package user

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/internal/cli/app"
	"github.com/razshare/frizzante/internal/platform"
)

func TestPlatformLinuxAmd64(t *testing.T) {
	PlatformMutex.Lock()
	defer PlatformMutex.Unlock()

	cache, err := FrizzanteCache()
	if err != nil {
		t.Fatal(err)
	}

	_ = os.Remove(filepath.Join(cache, "platform.txt"))
	defer func() { _ = os.Remove(filepath.Join(cache, "platform.txt")) }()

	platStr := "linux/amd64"
	a := &app.App{Platform: &platStr}

	plat, err := Platform(a)
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.LinuxAmd64 {
		t.Fatal("platform should be linux amd64")
	}

	if !files.IsFile(filepath.Join(cache, "platform.txt")) {
		t.Fatal("~/platform.txt should be a file")
	}

	d, err := os.ReadFile(filepath.Join(cache, "platform.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != "linux/amd64" {
		t.Fatal("~/platform.txt should contain linux/amd64")
	}
}

func TestPlatformLinuxArm64(t *testing.T) {
	PlatformMutex.Lock()
	defer PlatformMutex.Unlock()

	cache, err := FrizzanteCache()
	if err != nil {
		t.Fatal(err)
	}

	_ = os.Remove(filepath.Join(cache, "platform.txt"))
	defer func() { _ = os.Remove(filepath.Join(cache, "platform.txt")) }()

	platStr := "linux/arm64"
	a := &app.App{Platform: &platStr}

	plat, err := Platform(a)
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.LinuxArm64 {
		t.Fatal("platform should be linux arm64")
	}

	if !files.IsFile(filepath.Join(cache, "platform.txt")) {
		t.Fatal("~/platform.txt should be a file")
	}

	d, err := os.ReadFile(filepath.Join(cache, "platform.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != "linux/arm64" {
		t.Fatal("~/platform.txt should contain linux/arm64")
	}
}

func TestPlatformDarwinAmd64(t *testing.T) {
	PlatformMutex.Lock()
	defer PlatformMutex.Unlock()

	cache, err := FrizzanteCache()
	if err != nil {
		t.Fatal(err)
	}

	_ = os.Remove(filepath.Join(cache, "platform.txt"))
	defer func() { _ = os.Remove(filepath.Join(cache, "platform.txt")) }()

	platStr := "darwin/amd64"
	a := &app.App{Platform: &platStr}

	plat, err := Platform(a)
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.DarwinAmd64 {
		t.Fatal("platform should be darwin amd64")
	}

	if !files.IsFile(filepath.Join(cache, "platform.txt")) {
		t.Fatal("~/platform.txt should be a file")
	}

	d, err := os.ReadFile(filepath.Join(cache, "platform.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != "darwin/amd64" {
		t.Fatal("~/platform.txt should contain darwin/amd64")
	}
}

func TestPlatformDarwinArm64(t *testing.T) {
	PlatformMutex.Lock()
	defer PlatformMutex.Unlock()

	cache, err := FrizzanteCache()
	if err != nil {
		t.Fatal(err)
	}

	_ = os.Remove(filepath.Join(cache, "platform.txt"))
	defer func() { _ = os.Remove(filepath.Join(cache, "platform.txt")) }()

	platStr := "darwin/arm64"
	a := &app.App{Platform: &platStr}

	plat, err := Platform(a)
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.DarwinArm64 {
		t.Fatal("platform should be darwin arm64")
	}

	if !files.IsFile(filepath.Join(cache, "platform.txt")) {
		t.Fatal("~/platform.txt should be a file")
	}

	d, err := os.ReadFile(filepath.Join(cache, "platform.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != "darwin/arm64" {
		t.Fatal("~/platform.txt should contain darwin/arm64")
	}
}

func TestPlatformWindowsAmd64(t *testing.T) {
	PlatformMutex.Lock()
	defer PlatformMutex.Unlock()

	cache, err := FrizzanteCache()
	if err != nil {
		t.Fatal(err)
	}

	_ = os.Remove(filepath.Join(cache, "platform.txt"))
	defer func() { _ = os.Remove(filepath.Join(cache, "platform.txt")) }()

	platStr := "windows/amd64"
	a := &app.App{Platform: &platStr}

	plat, err := Platform(a)
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.WindowsAmd64 {
		t.Fatal("platform should be windows amd64")
	}

	if !files.IsFile(filepath.Join(cache, "platform.txt")) {
		t.Fatal("~/platform.txt should be a file")
	}

	d, err := os.ReadFile(filepath.Join(cache, "platform.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != "windows/amd64" {
		t.Fatal("~/platform.txt should contain windows/amd64")
	}
}

func TestPlatformWindowsArm64(t *testing.T) {
	PlatformMutex.Lock()
	defer PlatformMutex.Unlock()

	cache, err := FrizzanteCache()
	if err != nil {
		t.Fatal(err)
	}

	_ = os.Remove(filepath.Join(cache, "platform.txt"))
	defer func() { _ = os.Remove(filepath.Join(cache, "platform.txt")) }()

	platStr := "windows/arm64"
	a := &app.App{Platform: &platStr}

	plat, err := Platform(a)
	if err != nil {
		t.Fatal(err)
	}

	if plat != platform.WindowsArm64 {
		t.Fatal("platform should be windows arm64")
	}

	if !files.IsFile(filepath.Join(cache, "platform.txt")) {
		t.Fatal("~/platform.txt should be a file")
	}

	data, err := os.ReadFile(filepath.Join(cache, "platform.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "windows/arm64" {
		t.Fatal("~/platform.txt should contain windows/arm64")
	}
}

func TestTestPlatformFresh(t *testing.T) {
	PlatformMutex.Lock()
	defer PlatformMutex.Unlock()

	cache, err := FrizzanteCache()
	if err != nil {
		t.Fatal(err)
	}

	if err = os.RemoveAll(cache); err != nil {
		t.Fatal()
		return
	}

	platStr := "windows/arm64"
	a := &app.App{Platform: &platStr}
	_, err = Platform(a)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPlatformCached(t *testing.T) {
	PlatformMutex.Lock()
	defer PlatformMutex.Unlock()

	cache, err := FrizzanteCache()
	if err != nil {
		t.Fatal(err)
	}

	if err = os.WriteFile(filepath.Join(cache, "platform.txt"), []byte("windows/arm64"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	plat, err := Platform(&app.App{})
	if err != nil {
		_ = os.Remove(filepath.Join(cache, "platform.txt"))
		t.Fatal(err)
	}
	_ = os.Remove(filepath.Join(cache, "platform.txt"))

	if plat != platform.WindowsArm64 {
		t.Fatal("platform should be windows/arm64")
	}
}
