package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/metrue/fx/constants"
	containerruntimes "github.com/metrue/fx/container_runtimes"
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

		workdir := fmt.Sprintf("/tmp/fx-%d", time.Now().Unix())
		defer os.RemoveAll(workdir)

		body, err := ioutil.ReadFile(funcFile)
		if err != nil {
			log.Fatalf("function code load failed: %v", err)
			return err
		}
		log.Infof("function code loaded: %v", constants.CheckedSymbol)
		lang := utils.GetLangFromFileName(funcFile)

		fn := types.Func{Language: lang, Source: string(body)}

		if err := packer.PackIntoDir(fn, workdir); err != nil {
			log.Fatalf("could not pack function %v: %v", fn, err)
			return err
		}

		docker, ok := ctx.Get("docker").(containerruntimes.ContainerRuntime)
		if ok {
			nameWithTag := tag + ":latest"
			if err := docker.BuildImage(ctx.Context, workdir, nameWithTag); err != nil {
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
