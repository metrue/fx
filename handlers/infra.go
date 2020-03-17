package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/metrue/fx/context"
	k8sInfra "github.com/metrue/fx/infra/k8s"
	"github.com/metrue/fx/pkg/spinner"
	"github.com/metrue/fx/utils"
	"github.com/mitchellh/go-homedir"
)

func configDir() (string, error) {
	const defaultFxConfig = "~/.fx/config.yml"
	configFile, err := homedir.Expand(defaultFxConfig)
	if err != nil {
		return "", err
	}
	if os.Getenv("FX_CONFIG") != "" {
		configFile = os.Getenv("FX_CONFIG")
	}
	p, err := filepath.Abs(configFile)
	if err != nil {
		return "", err
	}
	return path.Dir(p), nil
}

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
	if typ == "k8s" {
		if cli.String("master") == "" {
			return fmt.Errorf("master required, eg. 'root@123.1.2.12'")
		}
	} else {
		return fmt.Errorf("invalid type, 'docker' and 'k8s' support")
	}

	switch strings.ToLower(typ) {
	case "k8s":
		dir, err := configDir()
		if err != nil {
			return err
		}
		kubeconf, err := setupK8S(dir, name, cli.String("master"), cli.String("agents"))
		if err != nil {
			return err
		}
		// TODO just write to ~/.fx/<name>.kubeconf

		kubeconfigPath := path.Join(dir, name+".kubeconf")
		if err := utils.EnsureFile(kubeconfigPath); err != nil {
			return err
		}
		if err := ioutil.WriteFile(kubeconfigPath, kubeconf, 0644); err != nil {
			return err
		}
	}
	return nil
}
