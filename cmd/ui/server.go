package main

//generate

//go:generate go-bindata -pkg ui -o bindata/ui/bindata.go cmd/ui/assets/dist/...

import (
	"fmt"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/supergiant/supergiant/bindata/ui"
)

type fsWithDefault struct {
	underlying http.FileSystem
	defaultDoc string // Filename of the 404 file to serve when there's an error serving original file.
}

func (fs fsWithDefault) Open(name string) (http.File, error) {
	f, err := fs.underlying.Open(name)
	if err != nil {
		// If there's an error (perhaps worth checking that the error is "file doesn't exist", up to you),
		// then serve your actual "404.html" file or handle it any way you wish.
		return fs.underlying.Open(fs.defaultDoc)
	}
	return f, err
}

func main() {
	fs := fsWithDefault{
		underlying: &assetfs.AssetFS{Asset: ui.Asset, AssetDir: ui.AssetDir, AssetInfo: ui.AssetInfo, Prefix: "cmd/ui/assets/dist/"},
		defaultDoc: "index.html",
	}

	http.Handle("/", http.FileServer(fs))

	fmt.Println("Serving on port 3001")
	http.ListenAndServe(":3001", nil)
}
