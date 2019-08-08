package api

import (
	"io/ioutil"

	"github.com/apex/log"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
)

// UpOptions options for up
type UpOptions struct {
	Name string
	Port int
}

// Up up a source code of function to be a service
func (api *API) Up(file string, opt UpOptions) error {
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

	project, err := api.Pack(opt.Name, fn)
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

	if err := api.Run(opt.Port, &service); err != nil {
		log.Fatalf("Run Service: %v", err)
		return err
	}
	log.Info("Run Service: \u2713")
	log.Infof("Service (%s) is running on: %s:%d", service.Name, service.Host, service.Port)

	return nil
}
