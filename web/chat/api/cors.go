package api

import (
	"net/http"
	"slices"
	"strings"

	"example/chat/constants"
)

func isPreflight(r *http.Request) bool {
	return r.Method == "OPTIONS" &&
		r.Header.Get("Origin") != "" &&
		r.Header.Get("Access-Control-Request-Method") != ""
}

func EnableCors(w *http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if isPreflight(r) {
		method := r.Header.Get("Access-Control-Request-Method")
		if slices.Contains(constants.OriginAllowlist, origin) && slices.Contains(constants.MethodAllowList, method) {
			// Access-Control-Allow-Origin indicates whether the response can be shared
			(*w).Header().Set("Access-Control-Allow-Origin", origin)
			// Access-Control-Allow-Methods - Indicates which methods are supported by the response's URL
			(*w).Header().Set("Access-Control-Allow-Methods", strings.Join(constants.MethodAllowList, ", "))

		}
	} else {
		if slices.Contains(constants.OriginAllowlist, origin) {
			(*w).Header().Set("Access-Control-Allow-Origin", origin)
		}
	}
	// Avoid problems with caching proxies between the server and the client
	// see also: https://fetch.spec.whatwg.org/#cors-protocol-and-http-caches
	(*w).Header().Set("Vary", "Origin")
	// Access-Control-Allow-Headers - Indicates which headers are supported
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
