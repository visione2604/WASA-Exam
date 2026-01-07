//go:build webui

package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/visione2604/WASA-Exam/webui"
)

// registerWebUI registers the embedded WebUI files to be served.
// This function is only compiled when the 'webui' build tag is specified.
func registerWebUI(hdl http.Handler) (http.Handler, error) {
	distDirectory, err := fs.Sub(webui.Dist, "dist")
	if err != nil {
		return nil, fmt.Errorf("error embedding WebUI dist/ directory: %w", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/dashboard/") {
			http.StripPrefix("/dashboard/", http.FileServer(http.FS(distDirectory))).ServeHTTP(w, r)
			return
		}
		hdl.ServeHTTP(w, r)
	}), nil
}
