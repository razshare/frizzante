package platforms

import (
	"runtime"
	"testing"
)

func TestPlatformLinuxAmd64(t *testing.T) {
	if runtime.GOOS+"/"+runtime.GOARCH == "linux/amd64" && Detect() != LinuxAmd64 {
		t.Fatal("platform should be linux amd64")
	}
}

func TestPlatformLinuxArm64(t *testing.T) {
	if runtime.GOOS+"/"+runtime.GOARCH == "linux/arm64" && Detect() != LinuxArm64 {
		t.Fatal("platform should be linux arm64")
	}
}

func TestPlatformDarwinAmd64(t *testing.T) {
	if runtime.GOOS+"/"+runtime.GOARCH == "darwin/amd64" && Detect() != DarwinAmd64 {
		t.Fatal("platform should be darwin amd64")
	}
}

func TestPlatformDarwinArm64(t *testing.T) {
	if runtime.GOOS+"/"+runtime.GOARCH == "darwin/arm64" && Detect() != DarwinArm64 {
		t.Fatal("platform should be darwin arm64")
	}
}

func TestPlatformWindowsAmd64(t *testing.T) {
	if runtime.GOOS+"/"+runtime.GOARCH == "windows/amd64" && Detect() != WindowsAmd64 {
		t.Fatal("platform should be windows amd64")
	}
}

func TestPlatformWindowsArm64(t *testing.T) {
	if runtime.GOOS+"/"+runtime.GOARCH == "windows/arm64" && Detect() != WindowsArm64 {
		t.Fatal("platform should be windows arm64")
	}
}
