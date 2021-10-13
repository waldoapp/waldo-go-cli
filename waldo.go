package main

import (
    "archive/zip"
	"fmt"
	"io"
	"io/fs"
	"net/http"
    "net/http/httputil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const waldoGoCLIVersion = "0.1.0"

var waldoBuildFlavor string
var waldoBuildPath string
var waldoBuildPayloadPath string
var waldoBuildSuffix string
var waldoBuildUploadID string
var waldoHistory string
var waldoHistoryError string
var waldoIncludeSymbols bool
var waldoPlatform string
var waldoSymbolsPath string
var waldoSymbolsPayloadPath string
var waldoSymbolsSuffix string
var waldoUploadToken string
var waldoVariantName string
var waldoVerbose bool
var waldoWorkingPath string

func checkBuildPath() {
    if len(waldoBuildPath) == 0 {
        failUsage(fmt.Errorf("Missing required argument: ‘path’"))
    }

    var err error

    waldoBuildPath, err = filepath.Abs(waldoBuildPath)

    if err != nil {
        fail(err)
    }

    waldoBuildSuffix = filepath.Ext(waldoBuildPath)[1:]

    switch waldoBuildSuffix {
    case "apk":
        waldoBuildFlavor = "Android"

    case "app", "ipa":
        waldoBuildFlavor = "iOS"

    default:
        fail(fmt.Errorf("File extension of build at ‘%s’ is not recognized", waldoBuildPath))
    }
}

func checkBuildStatus(resp *http.Response) {
    body, err := io.ReadAll(resp.Body)

    if err != nil {
        fail(err)
    }

    bodyString := string(body)

    statusRegex := regexp.MustCompile(`"status":([0-9]+)`)
    statusMatches := statusRegex.FindStringSubmatch(bodyString)

    if len(statusMatches) > 0 {     // status is numeric _only_ on failure
        var status = 0

        status, err = strconv.Atoi(statusMatches[1])

        if err != nil {
            fail(err)
        }

        if status == 401 {
            fail(fmt.Errorf("Upload token is invalid or missing!"))
        }

        if status < 200 || status > 299 {
            fail(fmt.Errorf("Unable to upload build to Waldo, HTTP status: %d", status))
        }
    }

    idRegex := regexp.MustCompile(`"id":"(appv-[0-9a-f]+)"`)
    idMatches := idRegex.FindStringSubmatch(bodyString)

    if len(idMatches) > 0 {
        waldoBuildUploadID = idMatches[1]
    }
}

func checkHistory() {
}

func checkPlatform() {
}

func checkSymbolsPath() {
}

func checkSymbolsStatus(resp *http.Response) {
}

func checkUploadToken() {
    if len(waldoUploadToken) == 0 {
        waldoUploadToken = os.Getenv("WALDO_UPLOAD_TOKEN")
    }

    if len(waldoUploadToken) == 0 {
        failUsage(fmt.Errorf("Missing required option: ‘--upload_token’"))
    }
}

func checkVariantName() {
    if len(waldoVariantName) == 0 {
        waldoVariantName = os.Getenv("WALDO_VARIANT_NAME")
    }
}

func createBuildPayload() {
    parentPath := filepath.Dir(waldoBuildPath)
    buildName := filepath.Base(waldoBuildPath)

    fi, err := os.Stat(waldoBuildPath)

    if err != nil {
        fail(err)
    }

    mode := fi.Mode()

    switch waldoBuildSuffix {
        case "app":
            if !mode.IsDir() {
                fail(fmt.Errorf("Unable to read build at ‘%s’", waldoBuildPath))
            }

            waldoBuildPayloadPath = filepath.Join(waldoWorkingPath, buildName + ".zip")

            err = zipDir(waldoBuildPayloadPath, parentPath, buildName)

            if err != nil {
                fail(err)
            }

        default:
            if !mode.IsRegular() {
                fail(fmt.Errorf("Unable to read build at ‘%s’", waldoBuildPath))
            }

            waldoBuildPayloadPath = waldoBuildPath
    }
}

func createSymbolsPayload() {
}

func createWorkingPath() {
    var err error

    waldoWorkingPath, err = os.MkdirTemp("", "WaldoCLI-*")

    if err != nil {
        fail(err)
    }

    os.RemoveAll(waldoWorkingPath)
    os.MkdirAll(waldoWorkingPath, 0755)
}

func deleteWorkingPath() {
    if len(waldoWorkingPath) > 0 {
        os.RemoveAll(waldoWorkingPath)
    }
}

func displaySummary() {
    fmt.Printf("\n")
    fmt.Printf("Build path:   %s\n", summarize(waldoBuildPath))
    fmt.Printf("Symbols path: %s\n", summarize(waldoSymbolsPath))
    fmt.Printf("Variant name: %s\n", summarize(waldoVariantName))
    fmt.Printf("Upload token: %s\n", summarizeSecure(waldoUploadToken))
    fmt.Printf("\n")

    if waldoVerbose {
        fmt.Printf("Build payload path:   %s\n", summarize(waldoBuildPayloadPath))
        fmt.Printf("Symbols payload path: %s\n", summarize(waldoSymbolsPayloadPath))
        fmt.Printf("\n")
    }
}

func displayUsage() {
    fmt.Printf(`
OVERVIEW: Upload build to Waldo

USAGE: waldo [options] <build-path> [<symbols-path>]

OPTIONS:

  --help                  Display available options
  --include_symbols       Include symbols with the build upload
  --upload_token <value>  Waldo upload token (overrides WALDO_UPLOAD_TOKEN)
  --variant_name <value>  Waldo variant name (overrides WALDO_VARIANT_NAME)
  --verbose               Display extra verbiage
`)
}

func displayVersion() {
    waldoPlatform=getPlatform()

    fmt.Printf("Waldo Go CLI %s (%s)\n", waldoGoCLIVersion, waldoPlatform)
}

func fail(err error) {
    message := fmt.Sprintf("waldo: %v", err)

    if len(waldoUploadToken) > 0 {
        if uploadError(err) {
            message += " -- Waldo team has been informed"
        }
    }

    fmt.Printf("\n")    // flush stdout

    os.Stderr.WriteString(fmt.Sprintf("%s\n", message))

    os.Exit(1)
}

func failUsage(err error) {
    if len(waldoUploadToken) > 0 {
        uploadError(err)
    }

    fmt.Printf("\n")    // flush stdout

    os.Stderr.WriteString(fmt.Sprintf("waldo: %v\n", err))

    displayUsage()

    os.Exit(1)
}

func findSymbolsPath() string {
    return ""
}

func getAuthorization() string {
    return fmt.Sprintf("Upload-Token %s", waldoUploadToken)
}

func getBuildContentType() string {
    switch waldoBuildSuffix {
        case "app":
            return "application/zip"

        default:
            return "application/octet-stream"
    }
}

func getCI() string {
    if len(os.Getenv("APPCENTER_BUILD_ID")) > 0 {
        return "App Center"
    }

    if os.Getenv("BITRISE_IO") == "true" {
        return "Bitrise"
    }

    if len(os.Getenv("BUDDYBUILD_BUILD_ID")) > 0 {
        return "buddybuild"
    }

    if os.Getenv("CIRCLECI") == "true" {
        return "CircleCI"
    }

    if os.Getenv("GITHUB_ACTIONS") == "true" {
        return "GitHub Actions"
    }

    if os.Getenv("TRAVIS") == "true" {
        return "Travis CI"
    }

    return "Go CLI"
}

func getErrorContentType() string {
    return "application/json"
}

func getPlatform() string {
    osName := runtime.GOOS

    switch osName {
        case "darwin":
            return "macOS"

        default:
            return strings.Title(osName)
    }
}

func getSymbolsContentType() string {
    return "application/zip"
}

func getUserAgent() string {
    return fmt.Sprintf("Waldo %s/%s v%s", getCI(), waldoBuildFlavor, waldoGoCLIVersion)
}

func main() {
    displayVersion()

    parseArgs()

    checkPlatform()
    checkBuildPath()
    checkSymbolsPath()
    checkHistory()
    checkUploadToken()
    checkVariantName()

    createWorkingPath()

    defer deleteWorkingPath()

    createBuildPayload()
    createSymbolsPayload()

    displaySummary()

    uploadBuild()
    uploadSymbols()
}

func makeBuildURL() string {
    query := ""

    if len(waldoHistory) > 0 {
        query += "&history="
        query += waldoHistory
    }

    if len(waldoHistoryError) > 0 {
        query += "&historyError="
        query += waldoHistoryError
    }

    if len(waldoVariantName) > 0 {
        query += "&waldoVariantName="
        query += waldoVariantName
    }

    url := "https://api.waldo.io/versions"

    if len(query) > 0 {
        url += "?"
        url += query[1:]
    }

    return url
}

func makeErrorURL() string {
    return "https://api.waldo.io/uploadError"
}

func makeSymbolsURL() string {
    return fmt.Sprintf("https://api.waldo.io/versions/%s/symbols", waldoBuildUploadID)
}

func parseArgs() {
    args := os.Args[1:]

    if len(args) == 0 {
        displayUsage()

        os.Exit(0)
    }

    for len(args) > 0 {
        arg := args[0]
        args = args[1:]

        switch arg {
            case "--help":
                displayUsage()

                os.Exit(0)

            case "--include_symbols":
                waldoIncludeSymbols = true

            case "--upload_token":
                if len(args) == 0 || len(args[0]) == 0 || strings.HasPrefix(args[0], "-") {
                    failUsage(fmt.Errorf("Missing required value for option: ‘%s’", arg))
                }

                waldoUploadToken = args[0]

                args = args[1:]


            case "--variant_name":
                if len(args) == 0 || len(args[0]) == 0 || strings.HasPrefix(args[0], "-") {
                    failUsage(fmt.Errorf("Missing required value for option: ‘%s’", arg))
                }

                waldoVariantName = args[0]

                args = args[1:]

            case "--verbose":
                waldoVerbose = true

            default:
                if strings.HasPrefix(arg, "-") {
                    failUsage(fmt.Errorf("Unknown option: ‘%s’", arg))
                }

                if len(waldoBuildPath) == 0 {
                    waldoBuildPath = arg
                } else if len(waldoSymbolsPath) == 0 {
                    waldoSymbolsPath = arg
                } else {
                    failUsage(fmt.Errorf("Unknown argument: ‘%s’", arg))
                }
        }
    }
}

func summarize(value string) string {
    if len(value) > 0 {
        return fmt.Sprintf("‘%s’", value)
    } else {
        return "(none)"
    }
}

func summarizeSecure(value string) string {
    if !waldoVerbose {
        prefix := value[0:6]
        suffixLen := len(value) - len(prefix)
        secure := "********************************"

        return prefix + secure[0:suffixLen]
    }

    return value
}

func uploadBuild() {
    fmt.Printf("Uploading build to Waldo\n")

    buildName := filepath.Base(waldoBuildPath)

    url := makeBuildURL()

    file, err := os.Open(waldoBuildPayloadPath)

    if err != nil {
        fail(fmt.Errorf("Unable to upload build to Waldo, error: %v, url: %s", err, url))
    }

    defer file.Close()

    client := &http.Client{}

    req, err := http.NewRequest("POST", url, file)

    if err != nil {
        fail(fmt.Errorf("Unable to upload build to Waldo, error: %v, url: %s", err, url))
    }

    req.Header.Add("Authorization", getAuthorization())
    req.Header.Add("Content-Type", getBuildContentType())
    req.Header.Add("User-Agent", getUserAgent())

    if waldoVerbose {
        dump, err := httputil.DumpRequestOut(req, false)

        if err == nil {
            fmt.Printf("\n--- Request ---\n%s\n", dump)
        }
    }

    resp, err := client.Do(req)

    if err != nil {
        fail(fmt.Errorf("Unable to upload build to Waldo, error: %v, url: %s", err, url))
    }

    if waldoVerbose {
        dump, err := httputil.DumpResponse(resp, true)

        if err == nil {
            fmt.Printf("\n--- Response ---\n%s\n", dump)
        }
    }

    checkBuildStatus(resp)

    defer resp.Body.Close()

    fmt.Printf("\nBuild ‘%s’ successfully uploaded to Waldo!\n", buildName)
}

func uploadError(err error) bool {
    body := fmt.Sprintf("{\"message\":\"%s\",\"ci\":\"%s\"}", err.Error(), getCI())

    url := makeErrorURL()

    client := &http.Client{}

    req, err := http.NewRequest("POST", url, strings.NewReader(body))

    if err != nil {
        return false
    }

    req.Header.Add("Authorization", getAuthorization())
    req.Header.Add("Content-Type", getErrorContentType())
    req.Header.Add("User-Agent", getUserAgent())

//    if waldoVerbose {
//        dump, err := httputil.DumpRequestOut(req, true)
//
//        if err == nil {
//            fmt.Printf("\n--- Request ---\n%s\n", dump)
//        }
//    }

    resp, err := client.Do(req)

    if err != nil {
        return false
    }

    defer resp.Body.Close()

//    if waldoVerbose {
//        dump, err := httputil.DumpResponse(resp, true)
//
//        if err == nil {
//            fmt.Printf("\n--- Response ---\n%s\n", dump)
//        }
//    }

    return true
}

func uploadSymbols() {
}

func zipDir(zipPath string, dirPath string, basePath string) error {
    err := os.Chdir(dirPath)

    if err != nil {
        return err
    }

    zipFile, err := os.Create(zipPath)

    if err != nil {
        return err
    }

    defer zipFile.Close()

    zipWriter := zip.NewWriter(zipFile)

    walker := func(path string, entry fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if entry.IsDir() {
            return nil
        }

        file, err := os.Open(path)

        if err != nil {
            return err
        }

        defer file.Close()

        zipEntry, err := zipWriter.Create(path)

        if err != nil {
            return err
        }

        _, err = io.Copy(zipEntry, file)

        return err
    }

    err = filepath.WalkDir(basePath, walker)

    err2 := zipWriter.Close()

    if err != nil {
        return err
    }

    return err2
}
