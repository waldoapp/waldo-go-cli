package tool

import (
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

type BuildDetector struct {
	ioStreams *lib.IOStreams
	verbose   bool
}

//-----------------------------------------------------------------------------

type FoundBuildPath struct {
	AbsPath   string
	RelPath   string
	BuildTool BuildTool
}

//-----------------------------------------------------------------------------

func NewBuildDetector(verbose bool, ioStreams *lib.IOStreams) *BuildDetector {
	return &BuildDetector{
		ioStreams: ioStreams,
		verbose:   verbose}
}

//-----------------------------------------------------------------------------

func NewFoundBuildPath(buildTool BuildTool, path string) *FoundBuildPath {
	absPath, err := filepath.Abs(path)

	if err != nil {
		absPath = path
	}

	return &FoundBuildPath{
		AbsPath:   absPath,
		RelPath:   lib.MakeRelativeToCWD(absPath),
		BuildTool: buildTool}
}

//-----------------------------------------------------------------------------

func (bd *BuildDetector) Detect(rootPath string) ([]*FoundBuildPath, error) {
	bd.ioStreams.Printf("\nSearching for possible build paths…\n")

	var results []*FoundBuildPath

	err := filepath.WalkDir(rootPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if bd.shouldSkip(filepath.Base(path), entry) {
			return filepath.SkipDir
		}

		if entry.IsDir() {
			skipChildren := false

			if expo.IsPossibleExpoContainer(path) {
				found := NewFoundBuildPath(BuildToolExpo, path)

				if bd.verbose {
					bd.ioStreams.Printf("\nFound possible Expo container: %q\n", path)
				}

				results = append(results, found)

				skipChildren = true
			}

			if flutter.IsPossibleFlutterContainer(path) {
				found := NewFoundBuildPath(BuildToolFlutter, path)

				if bd.verbose {
					bd.ioStreams.Printf("\nFound possible Flutter container: %q\n", path)
				}

				results = append(results, found)

				skipChildren = true
			}

			if gradle.IsPossibleGradleContainer(path) {
				found := NewFoundBuildPath(BuildToolGradle, path)

				if bd.verbose {
					bd.ioStreams.Printf("\nFound possible Gradle container: %q\n", path)
				}

				results = append(results, found)

				skipChildren = true
			}

			if ionic.IsPossibleIonicContainer(path) {
				found := NewFoundBuildPath(BuildToolIonic, path)

				if bd.verbose {
					bd.ioStreams.Printf("\nFound possible Ionic container: %q\n", path)
				}

				results = append(results, found)

				skipChildren = true
			}

			if reactnative.IsPossibleReactNativeContainer(path) {
				found := NewFoundBuildPath(BuildToolReactNative, path)

				if bd.verbose {
					bd.ioStreams.Printf("\nFound possible React Native container: %q\n")
				}

				results = append(results, found)

				skipChildren = true
			}

			if xcode.IsPossibleXcodeContainer(path) {
				found := NewFoundBuildPath(BuildToolXcode, path)

				if bd.verbose {
					bd.ioStreams.Printf("\nFound possible Xcode container: %q\n", path)
				}

				results = append(results, found)

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

	return results, nil
}

//-----------------------------------------------------------------------------

func (bd *BuildDetector) shouldSkip(name string, entry fs.DirEntry) bool {
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
