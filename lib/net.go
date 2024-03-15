package lib

import (
	"net/http"
	"net/http/httputil"
)

func DumpRequest(ios *IOStreams, req *http.Request, body bool) {
	dump, err := httputil.DumpRequestOut(req, body)

	if err != nil {
		return
	}

	ios.Printf("\n--- Request ---\n%s\n", dump)

}

func DumpResponse(ios *IOStreams, rsp *http.Response, body bool) {
	dump, err := httputil.DumpResponse(rsp, body)

	if err != nil {
		return
	}

	ios.Printf("\n--- Response ---\n%s\n", dump)
}

func ShouldRetry(rsp *http.Response) bool {
	switch rsp.StatusCode {
	case 408, 429, 500, 502, 503, 504:
		return true

	default:
		return false
	}
}
