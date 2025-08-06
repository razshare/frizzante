package cli

import (
	"embed"
)

type Platform uint

const PlatformLinuxAmd64 Platform = 0
const PlatformLinuxArm64 Platform = 1
const PlatformDarwinAmd64 Platform = 2
const PlatformDarwinArm64 Platform = 3
const PlatformWindowsAmd64 Platform = 4
const PlatformWindowsArm64 Platform = 5

type FeatureCopyInstruction struct {
	FeatureName          string
	OriginDirectory      string
	DestinationDirectory string
}

type Cli struct {
	Efs    embed.FS
	Parsed bool
}
