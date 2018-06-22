package server

import (
	"flag"
	"fmt"
	"log"

	"github.com/metrue/fx/api/service"
	"github.com/metrue/fx/common"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/env"
)

func pullBaseImage(verbose bool) {
	err := env.Init(verbose)
	if err != nil {
		common.HandleEnvError(err)
	}
}

// Start parses input and launches the fx server in a blocking process
func Start(verbose bool) error {
	flag.Parse()
	log.SetFlags(0)

	go pullBaseImage(verbose)

	go func() {
		err := service.Start(config.GrpcEndpoint)
		if err != nil {
			log.Fatal(err)
		}
	}()

	addr := fmt.Sprintf("%s:%s", config.Server["host"], config.Server["port"])
	log.Printf("fx serves on %s", addr)
	// log.Fatal(http.ListenAndServe(addr, nil))
	return Run(config.GrpcEndpoint, addr)
}
