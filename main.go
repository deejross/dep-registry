package main

import (
	"log"
	"net/http"

	"github.com/deejross/dep-registry/web"
)

func main() {
	http.HandleFunc("/", web.IndexHandler)
	log.Println(http.ListenAndServe(":8080", nil))
}
