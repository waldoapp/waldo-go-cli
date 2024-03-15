package tool

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/expo"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/flutter"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/gradle"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/ionic"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/reactnative"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/xcode"
)

type BuildPath struct {
	AbsPath   string
	RelPath   string
	BuildTool BuildTool
}

//-----------------------------------------------------------------------------

func DetectBuildPaths(rootPath string, verbose bool, ios *lib.IOStreams) ([]*BuildPath, error) {
	ios.Printf("\nFinding possible build paths…\n")

	var buildPaths []*BuildPath

	err := filepath.WalkDir(rootPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if shouldSkip(filepath.Base(path), entry) {
			return filepath.SkipDir
		}

		if entry.IsDir() {
			skipChildren := false

			if expo.IsPossibleExpoContainer(path) {
				buildPath := newBuildPath(BuildToolExpo, path)

				if verbose {
					ios.Printf("\nFound possible Expo container: %q\n", path)
				}

				buildPaths = append(buildPaths, buildPath)

				skipChildren = true
			}

			if flutter.IsPossibleFlutterContainer(path) {
				buildPath := newBuildPath(BuildToolFlutter, path)

				if verbose {
					ios.Printf("\nFound possible Flutter container: %q\n", path)
				}

				buildPaths = append(buildPaths, buildPath)

				skipChildren = true
			}

			if gradle.IsPossibleGradleContainer(path) {
				buildPath := newBuildPath(BuildToolGradle, path)

				if verbose {
					ios.Printf("\nFound possible Gradle container: %q\n", path)
				}

				buildPaths = append(buildPaths, buildPath)

				skipChildren = true
			}

			if ionic.IsPossibleIonicContainer(path) {
				buildPath := newBuildPath(BuildToolIonic, path)

				if verbose {
					ios.Printf("\nFound possible Ionic container: %q\n", path)
				}

				buildPaths = append(buildPaths, buildPath)

				skipChildren = true
			}

			if reactnative.IsPossibleReactNativeContainer(path) {
				buildPath := newBuildPath(BuildToolReactNative, path)

				if verbose {
					ios.Printf("\nFound possible React Native container: %q\n")
				}

				buildPaths = append(buildPaths, buildPath)

				skipChildren = true
			}

			if xcode.IsPossibleXcodeContainer(path) {
				buildPath := newBuildPath(BuildToolXcode, path)

				if verbose {
					ios.Printf("\nFound possible Xcode container: %q\n", path)
				}

				buildPaths = append(buildPaths, buildPath)

				skipChildren = true
			}

			if skipChildren {
				return filepath.SkipDir
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(buildPaths) == 0 {
		return nil, errors.New("No build paths buildPath")
	}

	return buildPaths, nil
}

//-----------------------------------------------------------------------------

func newBuildPath(buildTool BuildTool, path string) *BuildPath {
	absPath, err := filepath.Abs(path)

	if err != nil {
		absPath = path
	}

	return &BuildPath{
		AbsPath:   absPath,
		RelPath:   lib.MakeRelativeToCWD(absPath),
		BuildTool: buildTool}
}

func shouldSkip(name string, entry fs.DirEntry) bool {
	if !entry.IsDir() {
		return false
	}

	if strings.HasPrefix(name, ".") ||
		strings.HasSuffix(name, ".docset") ||
		strings.HasSuffix(name, ".framework") ||
		strings.HasSuffix(name, ".lproj") ||
		strings.HasSuffix(name, ".xcassets") ||
		strings.HasSuffix(name, ".xcodeproj") ||
		strings.HasSuffix(name, ".xcworkspace") ||
		name == "build" ||
		name == "Carthage" ||
		name == "CordovaLib" ||
		name == "fastlane" ||
		name == "gradle" ||
		name == "node_modules" ||
		name == "Pods" {
		return true
	}

	return false
}
