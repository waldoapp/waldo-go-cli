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
	AbsPath        string
	RelPath        string
	BuildTool      BuildTool
	AndroidSupport bool
	IosSupport     bool
}

//-----------------------------------------------------------------------------

func DetectBuildPaths(rootPath string, verbose bool, ios *lib.IOStreams) ([]*BuildPath, error) {
	ios.Printf("\nFinding possible build paths\n")

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

			if bp := checkContainer(BuildToolExpo, expo.IsPossibleContainer, path, verbose, ios); bp != nil {
				buildPaths = append(buildPaths, bp)

				skipChildren = true
			}

			if bp := checkContainer(BuildToolFlutter, flutter.IsPossibleContainer, path, verbose, ios); bp != nil {
				buildPaths = append(buildPaths, bp)

				skipChildren = true
			}

			if bp := checkContainer(BuildToolGradle, gradle.IsPossibleContainer, path, verbose, ios); bp != nil {
				buildPaths = append(buildPaths, bp)

				skipChildren = true
			}

			if bp := checkContainer(BuildToolIonic, ionic.IsPossibleContainer, path, verbose, ios); bp != nil {
				buildPaths = append(buildPaths, bp)

				skipChildren = true
			}

			if bp := checkContainer(BuildToolReactNative, reactnative.IsPossibleContainer, path, verbose, ios); bp != nil {
				buildPaths = append(buildPaths, bp)

				skipChildren = true
			}

			if bp := checkContainer(BuildToolXcode, xcode.IsPossibleContainer, path, verbose, ios); bp != nil {
				buildPaths = append(buildPaths, bp)

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

type checkFunc func(path string) (bool, bool)

func checkContainer(bt BuildTool, fn checkFunc, path string, verbose bool, ios *lib.IOStreams) *BuildPath {
	hasAndroid, hasIos := fn(path)

	if !hasAndroid && !hasIos {
		return nil
	}

	bp := newBuildPath(bt, path, hasAndroid, hasIos)

	if verbose {
		bp.describe(ios)
	}

	return bp
}

func newBuildPath(bt BuildTool, path string, androidSupport, iosSupport bool) *BuildPath {
	absPath, err := filepath.Abs(path)

	if err != nil {
		absPath = path
	}

	return &BuildPath{
		AbsPath:        absPath,
		RelPath:        lib.MakeRelativeToCWD(absPath),
		BuildTool:      bt,
		AndroidSupport: androidSupport,
		IosSupport:     iosSupport}
}

//-----------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------

func (bp *BuildPath) describe(ios *lib.IOStreams) {
	bt := bp.BuildTool

	ios.Printf("\nFound possible %v container: %q\n", bt, bp.RelPath)

	if bt.CanSupportAndroid() && !bp.AndroidSupport {
		ios.Printf("      (Android support may be missing)\n")
	}

	if bt.CanSupportIos() && !bp.IosSupport {
		ios.Printf("      (iOS support may be missing)\n")
	}
}
