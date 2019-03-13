package api

import "github.com/apex/log"

// Down destroy services
func (api *API) Down(names []string) error {
	for _, name := range names {
		if err := api.Stop(name); err != nil {
			log.Fatalf("Down Service %s: %v", name, err)
		} else {
			log.Infof("Down Service %s: \u2713", name)
		}
	}
	return nil
}
