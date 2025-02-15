package main

import (
	"net/http"
	"slices"
	"strings"
)

func EnableCors(w *http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if isPreflight(r) {
		method := r.Header.Get("Access-Control-Request-Method")
		if slices.Contains(OriginAllowlist, origin) && slices.Contains(MethodAllowList, method) {
			// Access-Control-Allow-Origin indicates whether the response can be shared
			(*w).Header().Set("Access-Control-Allow-Origin", origin)
			// Access-Control-Allow-Methods - Indicates which methods are supported by the response's URL
			(*w).Header().Set("Access-Control-Allow-Methods", strings.Join(MethodAllowList, ", "))

		}
	} else {
		if slices.Contains(OriginAllowlist, origin) {
			(*w).Header().Set("Access-Control-Allow-Origin", origin)
		}
	}
	// Avoid problems with caching proxies between the server and the client
	// see also: https://fetch.spec.whatwg.org/#cors-protocol-and-http-caches
	(*w).Header().Set("Vary", "Origin")
	// Access-Control-Allow-Headers - Indicates which headers are supported
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
