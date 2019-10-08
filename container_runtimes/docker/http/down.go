package api

import (
	"sync"

	"github.com/apex/log"
)

// Down destroy services
func (api *API) Down(names []string) error {
	var wg sync.WaitGroup

	for _, name := range names {
		wg.Add(1)
		go func(s string) {
			if err := api.Stop(s); err != nil {
				log.Fatalf("Down Service %s: %v", s, err)
			} else {
				log.Infof("Down Service %s: \u2713", s)
			}
			defer wg.Done()
		}(name)
	}

	wg.Wait()

	return nil
}
