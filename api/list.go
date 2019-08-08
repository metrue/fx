package api

import (
	"github.com/apex/log"
	"github.com/metrue/fx/utils"
)

// List services
func (api *API) List(name string) error {
	services, err := api.list(name)
	if err != nil {
		log.Fatalf("List Services: %v", err)
		return err
	}

	for _, service := range services {
		if err := utils.OutputJSON(service); err != nil {
			return err
		}
	}
	return nil
}
