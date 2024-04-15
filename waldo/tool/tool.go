package tool

type BuildTool int

const (
	BuildToolUnknown BuildTool = iota
	BuildToolExpo
	BuildToolFlutter
	BuildToolGradle
	BuildToolIonic
	BuildToolReactNative
	BuildToolXcode
)

func (bt BuildTool) CanSupportAndroid() bool {
	switch bt {
	case BuildToolXcode, BuildToolUnknown:
		return false

	default:
		return true
	}
}

func (bt BuildTool) CanSupportIos() bool {
	switch bt {
	case BuildToolGradle, BuildToolUnknown:
		return false

	default:
		return true
	}
}

func (bt BuildTool) String() string {
	return [...]string{
		"Unknown",
		"Expo",
		"Flutter",
		"Gradle",
		"Ionic",
		"React Native",
		"Xcode"}[bt]
}
