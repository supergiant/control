package main

//generate

//go:generate go-bindata -pkg ui -o bindata/ui/bindata.go cmd/ui/assets/dist/...

import (
	"fmt"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/supergiant/supergiant/bindata/ui"
)

func main() {
	assetDir := &assetfs.AssetFS{Asset: ui.Asset, AssetDir: ui.AssetDir, AssetInfo: ui.AssetInfo, Prefix: "cmd/ui/assets/dist/"}
	http.Handle("/", http.FileServer(assetDir))

	fmt.Println("Serving on port 3001")
	http.ListenAndServe(":3001", nil)
}
