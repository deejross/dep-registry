package main

import (
	"log"
	"net/http"

	"github.com/deejross/dep-registry/auth"
	"github.com/deejross/dep-registry/binstore"
	"github.com/deejross/dep-registry/gate"
	"github.com/deejross/dep-registry/metastore"
	"github.com/deejross/dep-registry/storemanager"
	"github.com/deejross/dep-registry/web"
)

func main() {
	// TODO: create Config to use here
	bs, err := binstore.NewBoltBinStore("binstore.bolt")
	if err != nil {
		log.Fatalln("While creating binstore:", err)
	}

	ms, err := metastore.NewBoltMetaStore("metastore.bolt")
	if err != nil {
		log.Fatalln("While creating metastore:", err)
	}

	sm := storemanager.NewStoreManager(bs, ms)
	tm := auth.NewTokenManager([]byte("super-secret-key"))
	gate := gate.NewGate(nil, sm, tm)
	router := web.NewRouter(gate)
	log.Println(http.ListenAndServe(":8080", router))
}
