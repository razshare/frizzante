package cli

import "embed"

type FeatureCopyInstruction struct {
	FeatureName          string
	OriginDirectory      string
	DestinationDirectory string
}

type Cli struct {
	Efs    embed.FS
	Parsed bool
}
