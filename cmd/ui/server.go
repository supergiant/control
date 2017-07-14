package main

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/supergiant/supergiant/bindata/ui"
)

func main() {
	assetDir := &assetfs.AssetFS{Asset: ui.Asset, AssetDir: ui.AssetDir, AssetInfo: ui.AssetInfo, Prefix: "assets"}}
	http.Handle("/", http.FileServer(assetDir))
	http.ListenAndServe(":3001", nil)
}
