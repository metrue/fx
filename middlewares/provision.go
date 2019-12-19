package middlewares

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	dockerHTTP "github.com/metrue/fx/container_runtimes/docker/http"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/infra"
	dockerInfra "github.com/metrue/fx/infra/docker"
	k8sInfra "github.com/metrue/fx/infra/k8s"
	"github.com/metrue/fx/types"
	"github.com/pkg/errors"
)

// Provision make sure infrastructure is healthy
func Provision(ctx context.Contexter) (err error) {
	fxConfig := ctx.Get("config").(*config.Config)
	meta, err := fxConfig.GetCurrentCloud()
	if err != nil {
		return err
	}
	cloudType, err := fxConfig.GetCurrentCloudType()
	if err != nil {
		return err
	}

	ctx.Set("cloud_type", cloudType)
	var cloud infra.Clouder
	switch cloudType {
	case types.CloudTypeK8S:
		cloud, err = k8sInfra.Load(meta)
	case types.CloudTypeDocker:
		cloud, err = dockerInfra.Load(meta)
	}
	if err != nil {
		return err
	}

	ok, err := cloud.IsHealth()
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("infrastrure is not health, please try to run create infrastructure use 'fx infra create ...' command")
	}

	ctx.Set("cloud", cloud)

	conf, err := cloud.GetConfig()
	if err != nil {
		return err
	}
	var deployer infra.Deployer
	if os.Getenv("KUBECONFIG") != "" {
		cloudType = types.CloudTypeK8S
		conf = os.Getenv("KUBECONFIG")
		ctx.Set("cloud_type", types.CloudTypeK8S)
	}

	if cloudType == types.CloudTypeDocker {
		var meta map[string]string
		if err := json.Unmarshal([]byte(conf), &meta); err != nil {
			return err
		}
		docker, err := dockerHTTP.Create(meta["ip"], constants.AgentPort)
		if err != nil {
			return errors.Wrapf(err, "please make sure docker is installed and running on your host")
		}

		// TODO should clean up, but it needed in middlewares.Build
		ctx.Set("docker", docker)
		deployer, err = dockerInfra.CreateDeployer(docker)
		if err != nil {
			return err
		}
	} else if cloudType == types.CloudTypeK8S {
		deployer, err = k8sInfra.CreateDeployer(conf)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupport cloud type %s, please make sure you config is correct", cloud.GetType())
	}

	ctx.Set("deployer", deployer)

	return nil
}
