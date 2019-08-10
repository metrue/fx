package api

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/docker/docker/api/types/container"
	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/types"
)

// UpOptions options for up
type UpOptions struct {
	Body       []byte
	Lang       string
	Name       string
	Port       int
	HealtCheck bool
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
	log.Infof("Pack Service: %s", constants.CheckedSymbol)

	service, err := api.Build(project)
	if err != nil {
		log.Fatalf("Build Service: %v", err)
		return err
	}
	log.Infof("Build Service: %s", constants.CheckedSymbol)

	if err := api.Run(opt.Port, &service); err != nil {
		log.Fatalf("Run Service: %v", err)
		return err
	}
	log.Infof("Run Service: %s", constants.CheckedSymbol)
	log.Infof("Service (%s) is running on: %s:%d", service.Name, service.Host, service.Port)

	if opt.HealtCheck {
		go func() {
			resultC, errC := api.ContainerWait(
				context.Background(),
				service.ID,
				container.WaitConditionNextExit,
				20*time.Second,
			)
			for {
				select {
				case res := <-resultC:
					var msg string
					if res.Error != nil {
						msg = res.Error.Message
					}
					log.Warnf("container exited: Code(%d) %s %s", res.StatusCode, msg, constants.UncheckedSymbol)
				case err := <-errC:
					log.Fatalf("wait container status exit: %s, %v", constants.UncheckedSymbol, err)
				}
			}
		}()

		trys := 0
		for {
			if trys > 2 {
				break
			}
			info, err := api.inspect(service.ID)
			if err != nil {
				log.Fatalf("healt checking failed: %v", err)
			}
			if info.State.Running {
				log.Info("service is running")
			} else {
				log.Warnf("service is %s", info.State.Status)
			}
			time.Sleep(1 * time.Second)
			trys++
		}
	}

	return nil
}
