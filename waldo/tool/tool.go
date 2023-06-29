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
