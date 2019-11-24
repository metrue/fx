package handlers

import (
	"fmt"
	"strings"

	"github.com/metrue/fx/config"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/infra/docker"
	"github.com/metrue/fx/infra/k3s"
	"github.com/metrue/fx/pkg/spinner"
)

func setupK3S(masterInfo string, agentsInfo string) ([]byte, error) {
	info := strings.Split(masterInfo, "@")
	if len(info) != 2 {
		return nil, fmt.Errorf("incorrect master info, should be <user>@<ip> format")
	}
	master := k3s.MasterNode{
		User: info[0],
		IP:   info[1],
	}
	agents := []k3s.AgentNode{}
	agentsInfoList := strings.Split(agentsInfo, ",")
	for _, agent := range agentsInfoList {
		info := strings.Split(agent, "@")
		if len(info) != 2 {
			return nil, fmt.Errorf("incorrect agent info, should be <user>@<ip> format")
		}
		agents = append(agents, k3s.AgentNode{
			User: info[0],
			IP:   info[1],
		})
	}
	k3sOperator := k3s.New(master, agents)
	return k3sOperator.Provision()
}

func setupDocker(hostInfo string) ([]byte, error) {
	info := strings.Split(hostInfo, "@")
	if len(info) != 2 {
		return nil, fmt.Errorf("incorrect master info, should be <user>@<ip> format")
	}
	user := info[1]
	host := info[0]
	dockr := docker.New(user, host)
	return dockr.Provision()
}

// Setup infra
func Setup(ctx *context.Context) (err error) {
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
	} else if typ == "k3s" {
		if cli.String("master") == "" {
			return fmt.Errorf("master required, eg. 'root@123.1.2.12'")
		}
	} else if typ == "k8s" {
	} else {
		return fmt.Errorf("invalid type, 'docker', 'k3s' and 'k8s' support")
	}

	fxConfig := ctx.Get("config").(*config.Config)

	switch strings.ToLower(typ) {
	case "k3s":
		kubeconf, err := setupK3S(cli.String("master"), cli.String("agents"))
		if err != nil {
			return err
		}
		return fxConfig.AddK8SCloud(name, kubeconf)
	case "docker":
		config, err := setupDocker(cli.String("host"))
		if err != nil {
			return err
		}
		return fxConfig.AddDockerCloud(name, config)
	case "k8s":
		fmt.Println("WIP")
	}
	return nil
}
