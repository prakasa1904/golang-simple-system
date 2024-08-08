package helper

import "net/http"

func IsHTMXRequest(header http.Header) bool {
	return header.Get("Hx-Request") == "true"
}
