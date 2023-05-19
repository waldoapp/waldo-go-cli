package data

import "strings"

type BuildFlavor string

const (
	BuildFlavorAndroid = "Android"
	BuildFlavorIos     = "iOS"
	BuildFlavorUnknown = "Unknown"
)

//-----------------------------------------------------------------------------

func ParseBuildFlavor(value string) BuildFlavor {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "android":
		return BuildFlavorAndroid

	case "ios":
		return BuildFlavorIos

	default:
		return BuildFlavorUnknown
	}
}
