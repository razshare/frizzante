package cli

type PlatformType uint

const PlatformTypeLinuxAmd64 PlatformType = 0
const PlatformTypeLinuxArm64 PlatformType = 1
const PlatformTypeDarwinAmd64 PlatformType = 2
const PlatformTypeDarwinArm64 PlatformType = 3
const PlatformTypeWindowsAmd64 PlatformType = 4
const PlatformTypeWindowsArm64 PlatformType = 5

type FeatureCopyInstruction struct {
	FeatureName          string
	OriginDirectory      string
	DestinationDirectory string
}
