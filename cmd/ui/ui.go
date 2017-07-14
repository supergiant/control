package main

import (
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./client/dist")))
	http.ListenAndServe(":3000", nil)
}
