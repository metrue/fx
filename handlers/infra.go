package handlers

import (
	"fmt"
	"strings"

	"github.com/metrue/fx/context"
	"github.com/metrue/fx/infra/k3s"
	"github.com/metrue/fx/pkg/spinner"
)

func setupK3S(masterInfo string, agentsInfo string) error {
	info := strings.Split(masterInfo, "@")
	if len(info) != 2 {
		return fmt.Errorf("incorrect master info, should be <user>@<ip> format")
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
			return fmt.Errorf("incorrect agent info, should be <user>@<ip> format")
		}
		agents = append(agents, k3s.AgentNode{
			User: info[0],
			IP:   info[1],
		})
	}
	k3sOperator := k3s.New(master, agents)
	if err := k3sOperator.SetupMaster(); err != nil {
		return err
	}
	if err := k3sOperator.SetupAgent(); err != nil {
		return err
	}
	return nil
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
	switch strings.ToLower(typ) {
	case "k3s":
		return setupK3S(cli.String("master"), cli.String("agents"))
	case "k8s":
		fmt.Println("WIP")
	case "docker":
		fmt.Println("docker")
	}
	return nil
}
