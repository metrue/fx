package handlers

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	api "github.com/metrue/fx/container_runtimes/docker/http"
	"github.com/metrue/fx/packer"
	"github.com/metrue/fx/provision"
	"github.com/metrue/fx/utils"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// BuildImage build image
func BuildImage(cfg config.Configer) HandleFunc {
	return func(ctx *cli.Context) error {
		funcFile := ctx.Args().First()
		tag := ctx.String("tag")
		if tag == "" {
			tag = uuid.New().String()
		}

		body, err := ioutil.ReadFile(funcFile)
		if err != nil {
			log.Fatalf("function code load failed: %v", err)
			return err
		}
		log.Infof("function code loaded: %v", constants.CheckedSymbol)
		lang := utils.GetLangFromFileName(funcFile)
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("could not get current work directory: %v", err)
			return err
		}
		tarFile := fmt.Sprintf("%s.%s.tar", pwd, tag)
		defer os.RemoveAll(tarFile)

		if err := packer.PackIntoTar(lang, string(body), tarFile); err != nil {
			log.Fatalf("could not pack function: %v", err)
			return err
		}
		log.Infof("function packed: %v", constants.CheckedSymbol)

		hosts, err := cfg.ListActiveMachines()
		if err != nil {
			log.Fatalf("could not list active machine: %v", err)
			return errors.Wrap(err, "list active machines failed")
		}

		if len(hosts) == 0 {
			log.Warnf("no active machines")
			return nil
		}
		for n, host := range hosts {
			if !host.Provisioned {
				provisionor := provision.New(host)
				if err := provisionor.Start(); err != nil {
					return errors.Wrapf(err, "could not provision %s", n)
				}
				log.Infof("provision machine %v: %s", n, constants.CheckedSymbol)
				if err := cfg.UpdateProvisionedStatus(n, true); err != nil {
					return errors.Wrap(err, "update machine provision status failed")
				}
			}

			if err := api.MustCreate(host.Host, constants.AgentPort).
				BuildImage(tarFile, tag, map[string]string{}); err != nil {
				return err
			}
			log.Infof("image built on machine %s: %v", n, constants.CheckedSymbol)
		}
		return nil
	}
}

// ExportImage export service's code into a directory
func ExportImage() HandleFunc {
	return func(ctx *cli.Context) (err error) {
		funcFile := ctx.Args().First()
		outputDir := ctx.String("output")
		if outputDir == "" {
			log.Fatalf("output directory required")
			return nil
		}

		body, err := ioutil.ReadFile(funcFile)
		if err != nil {
			return errors.Wrap(err, "read source failed")
		}
		lang := utils.GetLangFromFileName(funcFile)

		if err := packer.PackIntoDir(lang, string(body), outputDir); err != nil {
			log.Fatalf("write source code to file failed: %v", constants.UncheckedSymbol)
			return err
		}
		log.Infof("exported to %v: %v", outputDir, constants.CheckedSymbol)
		return nil
	}
}
