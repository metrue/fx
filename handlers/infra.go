package handlers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/metrue/fx/config"
	"github.com/metrue/fx/context"
	dockerInfra "github.com/metrue/fx/infra/docker"
	k8sInfra "github.com/metrue/fx/infra/k8s"
	"github.com/metrue/fx/pkg/spinner"
)

func setupK8S(configDir string, name, masterInfo string, agentsInfo string) ([]byte, error) {
	info := strings.Split(masterInfo, "@")
	if len(info) != 2 {
		return nil, fmt.Errorf("incorrect master info, should be <user>@<ip> format")
	}
	master, err := k8sInfra.CreateNode(info[1], info[0], "k3s_master", "master")
	if err != nil {
		return nil, err
	}
	nodes := []k8sInfra.Noder{master}
	if agentsInfo != "" {
		agentsInfoList := strings.Split(agentsInfo, ",")
		for idx, agent := range agentsInfoList {
			info := strings.Split(agent, "@")
			if len(info) != 2 {
				return nil, fmt.Errorf("incorrect agent info, should be <user>@<ip> format")
			}
			node, err := k8sInfra.CreateNode(info[1], info[0], "k3s_agent", fmt.Sprintf("agent-%d", idx))
			if err != nil {
				return nil, err
			}

			nodes = append(nodes, node)
		}
	}
	kubeconfigPath := filepath.Join(configDir, name+".kubeconfig")
	cloud := k8sInfra.NewCloud(kubeconfigPath, nodes...)
	if err := cloud.Provision(); err != nil {
		return nil, err
	}
	return cloud.Dump()
}

func setupDocker(hostInfo string, name string) ([]byte, error) {
	info := strings.Split(hostInfo, "@")
	if len(info) != 2 {
		return nil, fmt.Errorf("incorrect master info, should be <user>@<ip> format")
	}
	user := info[0]
	host := info[1]

	cloud, err := dockerInfra.Create(host, user, name)
	if err != nil {
		return nil, err
	}
	if err := cloud.Provision(); err != nil {
		return nil, err
	}
	return cloud.Dump()
}

// Setup infra
func Setup(ctx context.Contexter) (err error) {
	const task = "setup infra"
	spinner.Start(task)
	defer func() {
		spinner.Stop(task, err)
	}()

	cli := ctx.GetCliContext()
	typ := cli.String("type")
	name := cli.String("name")
	if name == "" {
		return fmt.Errorf("name required")
	}
	if typ == "docker" {
		if cli.String("host") == "" {
			return fmt.Errorf("host required, eg. 'root@123.1.2.12'")
		}
	} else if typ == "k8s" {
		if cli.String("master") == "" {
			return fmt.Errorf("master required, eg. 'root@123.1.2.12'")
		}
	} else {
		return fmt.Errorf("invalid type, 'docker' and 'k8s' support")
	}

	fxConfig := ctx.Get("config").(*config.Config)

	switch strings.ToLower(typ) {
	case "k8s":
		dir, err := fxConfig.Dir()
		if err != nil {
			return err
		}
		kubeconf, err := setupK8S(dir, name, cli.String("master"), cli.String("agents"))
		if err != nil {
			return err
		}
		return fxConfig.AddCloud(name, kubeconf)
	case "docker":
		config, err := setupDocker(cli.String("host"), name)
		if err != nil {
			return err
		}
		return fxConfig.AddCloud(name, config)
	}
	return nil
}
