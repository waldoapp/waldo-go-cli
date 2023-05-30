package tool

type BuildTool int

const (
	BuildToolCustom BuildTool = iota
	BuildToolExpo
	BuildToolFlutter
	BuildToolGradle
	BuildToolReactNative
	BuildToolXcode
)

func (bt BuildTool) String() string {
	return [...]string{
		"Custom",
		"Expo",
		"Flutter",
		"Gradle",
		"React Native",
		"Xcode"}[bt]
}
