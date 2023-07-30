package httputils

import (
	"net/http"
	"time"
)

const (
	HeaderXForwardedProto = "X-Forwarded-Proto"
	HeaderExpires         = "Expires"
	HeaderCacheControl    = "Cache-Control"
	HeaderPragma          = "Pragma"
)

func IsStatusCode2xx(code int) bool {
	return code >= http.StatusOK && code < http.StatusMultipleChoices
}

func WriteNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set(HeaderExpires, time.Unix(0, 0).Format(time.RFC1123))
	w.Header().Set(HeaderCacheControl, "no-cache, private, max-age=0")
	w.Header().Set(HeaderPragma, "no-cache")
}
