package api

import (
	"github.com/apex/log"
	"github.com/metrue/fx/types"
)

// UpOptions options for up
type UpOptions struct {
	Body []byte
	Lang string
	Name string
	Port int
}

// Up up a source code of function to be a service
func (api *API) Up(opt UpOptions) error {
	fn := types.ServiceFunctionSource{
		Language: opt.Lang,
		Source:   string(opt.Body),
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
