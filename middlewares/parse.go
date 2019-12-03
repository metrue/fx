package middlewares

import (
	"io/ioutil"

	"github.com/metrue/fx/context"
	"github.com/metrue/fx/pkg/spinner"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
	"github.com/pkg/errors"
)

// Parse parse input
func Parse(ctx context.Contexter) (err error) {
	const task = "parsing"
	spinner.Start(task)
	defer func() {
		spinner.Stop(task, err)
	}()

	cli := ctx.GetCliContext()
	funcFile := cli.Args().First()
	lang := utils.GetLangFromFileName(funcFile)
	body, err := ioutil.ReadFile(funcFile)
	if err != nil {
		return errors.Wrap(err, "read source failed")
	}
	fn := types.Func{
		Language: lang,
		Source:   string(body),
	}
	ctx.Set("fn", fn)

	name := cli.String("name")
	ctx.Set("name", name)
	port := cli.Int("port")
	ctx.Set("port", port)

	return nil
}
