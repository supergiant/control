package ui

import (
	"net/http"
	"strings"
	"net/url"
	"github.com/sirupsen/logrus"
)

func trimPrefix(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This code path is for static resources
		if len(strings.Split(r.URL.Path, ".")) == 1 {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = "/"
			logrus.Debugf("Change path URL %s to %s",
				r.URL.Path, r2.URL.Path)
			h.ServeHTTP(w, r2)
		} else {
			// This codepath is for URL paths from browser line
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = r2.URL.Path[1:]
			logrus.Debugf("Change asset URL %s to %s",
				r.URL.Path, r2.URL.Path)
			h.ServeHTTP(w, r2)
		}
	})
}
