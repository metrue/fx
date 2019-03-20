package api

import (
	"io/ioutil"

	"github.com/apex/log"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
)

// Up up a source code of function to be a service
func (api *API) Up(name string, file string) error {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Read Source: %v", err)
		return err
	}
	log.Info("Read Source: \u2713")

	lang := utils.GetLangFromFileName(file)
	fn := types.ServiceFunctionSource{
		Language: lang,
		Source:   string(src),
	}

	project, err := api.Pack(name, fn)
	if err != nil {
		log.Fatalf("Pack Service: %v", err)
		return err
	}
	log.Info("Pack Service: \u2713")

	service, err := api.Build(project)
	if err != nil {
		log.Fatalf("Build Service: %v", err)
		return err
	}
	log.Info("Build Service: \u2713")

	if err := api.Run(&service); err != nil {
		log.Fatalf("Run Service: %v", err)
		return err
	}
	log.Info("Run Service: \u2713")

	return nil
}
