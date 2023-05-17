package lib

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	BinaryContentType = "application/octet-stream"
	JsonContentType   = "application/json"
	ZipContentType    = "application/zip"
)

func AddIfNotEmpty(query *url.Values, key string, value string) {
	if len(key) > 0 && len(value) > 0 {
		query.Add(key, value)
	}
}

func DumpRequest(ioStreams *IOStreams, request *http.Request, body bool) {
	dump, err := httputil.DumpRequestOut(request, body)

	if err == nil {
		ioStreams.Printf("\n--- Request ---\n%s\n", dump)
	}
}

func DumpResponse(ioStreams *IOStreams, response *http.Response, body bool) {
	dump, err := httputil.DumpResponse(response, body)

	if err == nil {
		ioStreams.Printf("\n--- Response ---\n%s\n", dump)
	}
}

func ShouldRetry(rsp *http.Response) bool {
	switch rsp.StatusCode {
	case 408, 429, 500, 502, 503, 504:
		return true

	default:
		return false
	}
}
