package middlewares

import (
	"fmt"
	"strings"

	"github.com/metrue/fx/context"
	dockerInfra "github.com/metrue/fx/infra/docker"
)

// Provision make sure infrastructure is healthy
func Provision(ctx context.Contexter) (err error) {
	host := ctx.Get("host").(string)
	port := ctx.Get("ssh_port").(string)
	keyfile := ctx.Get("ssh_key").(string)
	kubeconf := ctx.Get("kubeconf").(string)
	if host == "" && kubeconf == "" {
		return fmt.Errorf("at least host or kubeconf provided")
	}

	if host != "" {
		addr := strings.Split(host, "@")
		if len(addr) != 2 {
			return fmt.Errorf("invalid host information, should be format of <user>@<ip>")
		}
		user := addr[0]
		ip := addr[1]

		cloud, err := dockerInfra.Create(ip, user, port, keyfile)
		if err != nil {
			return err
		}
		if err := cloud.Provision(); err != nil {
			return err
		}

		ok, err := cloud.IsHealth()
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("target docker host is not healthy")
		}
	}

	return nil
}
