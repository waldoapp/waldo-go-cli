package tool

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
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
			// if bd.verbose {
			// 	bd.ioStreams.Printf("\nSkipping directory %q\n", path)
			// }

			return filepath.SkipDir
		}

		if entry.IsDir() {
			// if bd.verbose {
			// 	bd.ioStreams.Printf("\nChecking directory %q\n", path)
			// }

			if IsPossibleExpoContainer(path) {
				found := NewFoundBuildPath(BuildToolExpo, path)

				if bd.verbose {
					bd.ioStreams.Printf("\nFound possible Expo container: %q\n", path)
				}

				results = append(results, found)
			}

			if IsPossibleFlutterContainer(path) {
				found := NewFoundBuildPath(BuildToolFlutter, path)

				if bd.verbose {
					bd.ioStreams.Printf("\nFound possible Flutter container: %q\n", path)
				}

				results = append(results, found)
			}

			if IsPossibleGradleContainer(path) {
				found := NewFoundBuildPath(BuildToolGradle, path)

				if bd.verbose {
					bd.ioStreams.Printf("\nFound possible Gradle container: %q\n", path)
				}

				results = append(results, found)
			}

			if IsPossibleReactNativeContainer(path) {
				found := NewFoundBuildPath(BuildToolReactNative, path)

				if bd.verbose {
					bd.ioStreams.Printf("\nFound possible React Native container: %q\n")
				}

				results = append(results, found)
			}

			if IsPossibleXcodeContainer(path) {
				found := NewFoundBuildPath(BuildToolXcode, path)

				if bd.verbose {
					bd.ioStreams.Printf("\nFound possible Xcode container: %q\n", path)
				}

				results = append(results, found)
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
		strings.HasSuffix(name, ".lproj") ||
		strings.HasSuffix(name, ".xcassets") ||
		strings.HasSuffix(name, ".xcodeproj") ||
		strings.HasSuffix(name, ".xcworkspace") ||
		name == "build" ||
		name == "fastlane" ||
		name == "gradle" {
		return true
	}

	return false
}
