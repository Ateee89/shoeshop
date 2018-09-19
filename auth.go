package main

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, r *http.Request) {
		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(resp, "Authorization failed", http.StatusUnauthorized)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !validate(pair[0], pair[1]) {
			http.Error(resp, "Authorization failed", http.StatusUnauthorized)
			return
		}

		h(resp, r)
	}
}

func validate(username, password string) bool {
	if username == "admin" && password == "admin" {
		return true
	}
	return false
}
