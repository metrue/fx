package handlers

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/apex/log"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/pkg/command"
	"github.com/metrue/fx/pkg/spinner"
	"github.com/metrue/fx/provision"
	"github.com/pkg/errors"
)

func setup(host string, user string, passwd string, script string) (err error) {
	task := fmt.Sprintf("setup %s@%s", user, host)
	spinner.Start(task)
	defer func() {
		spinner.Stop(task, err)
	}()

	runner := command.NewLocalRunner()
	cmd := command.New(task,
		script, // TODO install k3sup first
		runner,
	)
	if output, err := cmd.Exec(); err != nil {
		msg := fmt.Sprintf("%s: %s", task, string(output))
		return errors.Wrapf(err, msg)
	}
	return nil
}

// Init start fx-agent
func Init() HandleFunc {
	return func(ctx *context.Context) (err error) {
		const task = "init"
		spinner.Start(task)
		defer func() {
			spinner.Stop(task, err)
		}()

		cli := ctx.GetCliContext()
		master := cli.String("master")
		// TODO support different user on master and node
		acount := cli.String("user")
		passwd := cli.String("password")
		agents := strings.Split(cli.String("agents"), ",")
		if acount == "" {
			acount = "root"
		}
		if master != "" {
			// TODO install k3sup first
			script := fmt.Sprintf("k3sup install --ip %s --user %s", master, acount)
			if err := setup(master, acount, passwd, script); err != nil {
				return err
			}
		}
		if master != "" && len(agents) > 0 {
			var wg sync.WaitGroup
			for _, agent := range agents {
				wg.Add(1)
				go func(host string) error {
					script := fmt.Sprintf("k3sup join --ip %s  --server-ip %s  --user %s", host, master, acount)
					err := setup(host, acount, passwd, script)
					wg.Done()
					return err
				}(agent)
			}
			wg.Wait()
		}

		host := os.Getenv("DOCKER_REMOTE_HOST_ADDR")
		user := os.Getenv("DOCKER_REMOTE_HOST_USER")
		passord := os.Getenv("DOCKER_REMOTE_HOST_PASSWORD")
		if host == "" {
			host = "127.0.0.1"
		}
		provisioner := provision.NewWithHost(host, user, passord)
		if !provisioner.IsFxAgentRunning() {
			if err := provisioner.StartFxAgent(); err != nil {
				log.Fatalf("could not start fx agent on host: %s", err)
				return err
			}
			log.Info("fx agent started")
		}
		log.Info("fx agent already started")
		return nil
	}
}
