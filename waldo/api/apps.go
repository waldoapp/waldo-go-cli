package api

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

//-----------------------------------------------------------------------------

type AppInfo struct {
	AppID    string
	AppName  string
	Platform lib.Platform
}

type FetchAppsResponse struct {
	Items []*FetchedAppItem `json:"items"`
}

type FetchedAppItem struct {
	AppID string `json:"id"`
	Name  string `json:"name,omitempty"`
	Type  string `json:"type"`
}

//-----------------------------------------------------------------------------

func FetchApps(apiToken string, platform lib.Platform, verbose bool, ios *lib.IOStreams) ([]*AppInfo, error) {
	var far *FetchAppsResponse

	client := &http.Client{}

	req, err := http.NewRequest("GET", makeURL(platform), nil)

	if err != nil {
		return nil, fmt.Errorf("Unable to fetch apps, error: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token %v", apiToken))
	req.Header.Add("User-Agent", data.FullVersion())

	if verbose {
		lib.DumpRequest(ios, req, true)
	}

	rsp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Unable to fetch apps, error: %v", err)
	}

	defer rsp.Body.Close()

	if verbose {
		lib.DumpResponse(ios, rsp, true)
	}

	status := rsp.StatusCode

	if status < 200 || status > 299 {
		return nil, fmt.Errorf("Unable to fetch apps, error: %v", rsp.Status)
	}

	far, err = parseFetchAppsResponse(rsp)

	if err != nil {
		return nil, fmt.Errorf("Unable to fetch apps, error: %v", err)
	}

	apps := lib.CompactMap(far.Items, func(item *FetchedAppItem) (*AppInfo, bool) {
		if len(item.AppID) > 0 && len(item.Name) > 0 {
			return &AppInfo{
				AppID:    item.AppID,
				AppName:  item.Name,
				Platform: lib.ParsePlatform(item.Type)}, true
		}

		return nil, false
	})

	slices.SortStableFunc(apps, func(a, b *AppInfo) int {
		return cmp.Compare(strings.ToLower(a.AppName), strings.ToLower(b.AppName))
	})

	return apps, nil
}

//-----------------------------------------------------------------------------

func makeURL(platform lib.Platform) string {
	endpoint := getFetchAppsEndpoint()

	if ptype := platformType(platform); len(ptype) > 0 {
		return fmt.Sprintf("%v?type=%v", endpoint, ptype)
	}

	return endpoint
}

func parseFetchAppsResponse(rsp *http.Response) (*FetchAppsResponse, error) {
	data, err := io.ReadAll(rsp.Body)

	if err != nil {
		return nil, err
	}

	far := &FetchAppsResponse{}

	if err = json.Unmarshal(data, &far.Items); err != nil {
		return nil, err
	}

	return far, nil
}

func platformType(platform lib.Platform) string {
	switch platform {
	case lib.PlatformAndroid:
		return "android"

	case lib.PlatformIos:
		return "ios"

	default:
		return ""
	}
}
