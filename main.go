package main

import (
	"log"
	"net/http"
	"os"

	"github.com/deejross/dep-registry/auth"
	"github.com/deejross/dep-registry/binstore"
	"github.com/deejross/dep-registry/config"
	"github.com/deejross/dep-registry/gate"
	"github.com/deejross/dep-registry/metastore"
	"github.com/deejross/dep-registry/storemanager"
	"github.com/deejross/dep-registry/web"
)

func main() {
	cfg := &config.Config{}
	if len(os.Args) > 1 {
		name := os.Args[1]
		cfg2, err := config.FromFile(cfg, name)
		if err != nil {
			log.Println("While loading config:", err, ", using default values")
		}
		cfg = cfg2
	}

	cfg = config.FromEnvironment(cfg)
	if err := cfg.Validate(); err != nil {
		log.Fatalln(err)
	}

	bs, err := binstore.Resolve(cfg.BinStorePath)
	if err != nil {
		log.Fatalln("While creating binstore:", err)
	}

	ms, err := metastore.Resolve(cfg.MetaStorePath)
	if err != nil {
		log.Fatalln("While creating metastore:", err)
	}

	sm := storemanager.NewStoreManager(bs, ms)
	tm := auth.NewTokenManager([]byte(cfg.SigningKey))
	gate := gate.NewGate(nil, sm, tm)
	router := web.NewRouter(gate)
	log.Println(http.ListenAndServe(":"+cfg.Port, router))
}
