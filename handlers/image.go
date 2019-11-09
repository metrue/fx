package handlers

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/metrue/fx/constants"
	dockerHTTP "github.com/metrue/fx/container_runtimes/docker/http"
	dockerSDK "github.com/metrue/fx/container_runtimes/docker/sdk"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/packer"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
	"github.com/pkg/errors"
)

// BuildImage build image
func BuildImage() HandleFunc {
	return func(ctx *context.Context) error {
		cli := ctx.GetCliContext()
		funcFile := cli.Args().First()
		tag := cli.String("tag")
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

		if err := packer.PackIntoTar(types.Func{Language: lang, Source: string(body)}, tarFile); err != nil {
			log.Fatalf("could not pack function: %v", err)
			return err
		}
		log.Infof("function packed: %v", constants.CheckedSymbol)

		dockerAPI, ok := ctx.Get("docker_http").(*dockerHTTP.API)
		if ok {
			if err := dockerAPI.BuildImage(tarFile, tag, map[string]string{}); err != nil {
				return err
			}
			log.Infof("image built: %v", constants.CheckedSymbol)
			return nil
		}
		dockerCli, ok := ctx.Get("docker_sdk").(*dockerSDK.Docker)
		if ok {
			if err := dockerCli.BuildImage(ctx.Context, pwd, tag); err != nil {
				return err
			}
			log.Infof("image built: %v", constants.CheckedSymbol)
			return nil
		}
		return fmt.Errorf("no available docker cli")
	}
}

// ExportImage export service's code into a directory
func ExportImage() HandleFunc {
	return func(ctx *context.Context) (err error) {
		cli := ctx.GetCliContext()
		funcFile := cli.Args().First()
		outputDir := cli.String("output")
		if outputDir == "" {
			log.Fatalf("output directory required")
			return nil
		}

		body, err := ioutil.ReadFile(funcFile)
		if err != nil {
			return errors.Wrap(err, "read source failed")
		}
		lang := utils.GetLangFromFileName(funcFile)

		if err := packer.PackIntoDir(types.Func{Language: lang, Source: string(body)}, outputDir); err != nil {
			log.Fatalf("write source code to file failed: %v", constants.UncheckedSymbol)
			return err
		}
		log.Infof("exported to %v: %v", outputDir, constants.CheckedSymbol)
		return nil
	}
}
