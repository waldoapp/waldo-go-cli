package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	cliName    = "Waldo CLI"
	cliVersion = "2.0.4"

	cliAssetBaseURL = "https://github.com/waldoapp/waldo-go-agent/releases"
)

var (
	cliAgentPath   string
	cliAssetURL    string
	cliWorkingPath string

	cliArch         = detectArch()
	cliAssetVersion = detectAssetVersion()
	cliPlatform     = detectPlatform()
	cliVerbose      = detectVerbose()
)

func cleanupTarget() {
	os.RemoveAll(cliWorkingPath)
}

func detectArch() string {
	arch := runtime.GOARCH

	switch arch {
	case "amd64":
		return "x86_64"

	default:
		return arch
	}
}

func detectAssetVersion() string {
	if version := os.Getenv("WALDO_CLI_ASSET_VERSION"); len(version) > 0 {
		return version
	}

	return "latest"
}

func detectPlatform() string {
	platform := runtime.GOOS

	switch platform {
	case "darwin":
		return "macOS"

	default:
		return strings.Title(platform)
	}
}

func detectVerbose() bool {
	if verbose := os.Getenv("WALDO_CLI_VERBOSE"); verbose == "1" {
		return true
	}

	return false
}

func determineAgentPath() string {
	agentName := "waldo-agent"

	if cliPlatform == "windows" {
		agentName += ".exe"
	}

	return filepath.Join(cliWorkingPath, "waldo-agent")
}

func determineAssetURL() string {
	assetName := fmt.Sprintf("waldo-agent-%s-%s", cliPlatform, cliArch)

	if cliPlatform == "windows" {
		assetName += ".exe"
	}

	assetBaseURL := cliAssetBaseURL

	if cliAssetVersion != "latest" {
		assetBaseURL += "/download/" + cliAssetVersion
	} else {
		assetBaseURL += "/latest/download"
	}

	return assetBaseURL + "/" + assetName
}

func determineWorkingPath() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("WaldoGoCLI-%d", os.Getpid()))
}

func displayVersion() {
	if cliVerbose {
		fmt.Printf("%s %s (%s/%s)\n", cliName, cliVersion, cliPlatform, cliArch)
	}
}

func downloadAgent() {
	fmt.Printf("\nDownloading latest Waldo Agent…\n\n")

	client := &http.Client{}

	req, err := http.NewRequest("GET", cliAssetURL, nil)

	var resp *http.Response

	if err == nil {
		dumpRequest(req, false)

		resp, err = client.Do(req)
	}

	if err == nil {
		defer resp.Body.Close()

		dumpResponse(resp, false)

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			fail(fmt.Errorf("Unable to download Waldo Agent, HTTP status: %s", resp.Status))
		}
	}

	var file *os.File = nil

	if err == nil {
		file, err = os.OpenFile(cliAgentPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0775)
	}

	if err == nil {
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
	}

	if err != nil {
		fail(fmt.Errorf("Unable to download Waldo Agent, error: %v, url: %s", err, cliAssetURL))
	}
}

func dumpRequest(req *http.Request, body bool) {
	if cliVerbose {
		dump, err := httputil.DumpRequestOut(req, body)

		if err == nil {
			fmt.Printf("\n--- Request ---\n%s\n", dump)
		}
	}
}

func dumpResponse(resp *http.Response, body bool) {
	if cliVerbose {
		dump, err := httputil.DumpResponse(resp, body)

		if err == nil {
			fmt.Printf("\n--- Response ---\n%s\n", dump)
		}
	}
}

func enrichEnvironment() []string {
	env := os.Environ()

	//
	// If _both_ wrapper overrides are already set, do _not_ replace them with
	// the CLI name/version:
	//
	wrapperName := os.Getenv("WALDO_WRAPPER_NAME_OVERRIDE")
	wrapperVersion := os.Getenv("WALDO_WRAPPER_VERSION_OVERRIDE")

	if len(wrapperName) == 0 || len(wrapperVersion) == 0 {
		setEnvironVar(&env, "WALDO_WRAPPER_NAME_OVERRIDE", cliName)
		setEnvironVar(&env, "WALDO_WRAPPER_VERSION_OVERRIDE", cliVersion)
	}

	return env
}

func execAgent() {
	args := os.Args[1:]

	cmd := exec.Command(cliAgentPath, args...)

	cmd.Env = enrichEnvironment()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()

	if ee, ok := err.(*exec.ExitError); ok {
		os.Exit(ee.ExitCode())
	}
}

func fail(err error) {
	fmt.Printf("\n") // flush stdout

	os.Stderr.WriteString(fmt.Sprintf("waldo-cli: %v\n", err))

	os.Exit(1)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fail(fmt.Errorf("Unhandled panic: %v", err))
		}
	}()

	displayVersion()

	prepareSource()
	prepareTarget()

	defer cleanupTarget()

	downloadAgent()
	execAgent()
}

func prepareSource() {
	cliAssetURL = determineAssetURL()
}

func prepareTarget() {
	cliWorkingPath = determineWorkingPath()
	cliAgentPath = determineAgentPath()

	err := os.RemoveAll(cliWorkingPath)

	if err == nil {
		err = os.MkdirAll(cliWorkingPath, 0755)
	}

	if err != nil {
		fail(err)
	}
}

func setEnvironVar(env *[]string, key, value string) {
	for idx := range *env {
		if strings.HasPrefix((*env)[idx], key+"=") {
			(*env)[idx] = key + "=" + value

			return
		}
	}

	*env = append(*env, key+"="+value)
}
