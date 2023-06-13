package tool

type BuildTool int

const (
	BuildToolCustom BuildTool = iota
	BuildToolExpo
	BuildToolFlutter
	BuildToolGradle
	BuildToolIonic
	BuildToolReactNative
	BuildToolXcode
)

func (bt BuildTool) String() string {
	return [...]string{
		"Custom",
		"Expo",
		"Flutter",
		"Gradle",
		"Ionic",
		"React Native",
		"Xcode"}[bt]
}
