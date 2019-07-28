package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/golang/mock/gomock"
	"github.com/metrue/fx/config"
	mockConfig "github.com/metrue/fx/config/mocks"
	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/types"
	gock "gopkg.in/h2non/gock.v1"
)

func TestMakeTar(t *testing.T) {
	serviceName := "mock-service-abc"
	project := types.Project{
		Name:     serviceName,
		Language: "node",
		Files: []types.ProjectSourceFile{
			types.ProjectSourceFile{
				Path: "Dockerfile",
				Body: `
FROM metrue/fx-node-base

COPY . .
EXPOSE 3000
CMD ["node", "app.js"]`,
				IsHandler: false,
			},
			types.ProjectSourceFile{
				Path: "app.js",
				Body: `
const Koa = require('koa');
const bodyParser = require('koa-bodyparser');
const func = require('./fx');

const app = new Koa();
app.use(bodyParser());
app.use(ctx => {
  const msg = func(ctx.request.body);
  ctx.body = msg;
});

app.listen(3000);`,
				IsHandler: false,
			},
			types.ProjectSourceFile{
				Path: "fx.js",
				Body: `
module.exports = (input) => {
    return input.a + input.b
}
					`,
				IsHandler: true,
			},
		},
	}
	tarDir, err := ioutil.TempDir("/tmp", "fx-tar")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tarDir)

	tarFilePath := filepath.Join(tarDir, fmt.Sprintf("%s.tar", serviceName))
	if err := makeTar(project, tarFilePath); err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(tarFilePath)
	if err != nil {
		t.Fatal(err)
	}
	stat, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if stat.Name() != serviceName+".tar" {
		t.Fatalf("should get %s but got %s", serviceName+".tar", stat.Name())
	}
	if stat.Size() <= 0 {
		t.Fatalf("tarfile invalid: size %d", stat.Size())
	}
}

func TestBuild(t *testing.T) {
	defer gock.Off()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	host := config.Host{Host: "127.0.0.1"}
	cfg := mockConfig.NewMockConfiger(ctrl)
	cfg.EXPECT().GetDefaultHost().Return(host, nil)
	box := packr.NewBox("./images")
	api := New(cfg, box)
	if err := api.Init(); err != nil {
		t.Fatal(err)
	}

	url := "http://" + host.Host + ":" + constants.AgentPort
	gock.New(url).
		Post("/v" + api.version + "/build").
		AddMatcher(func(req *http.Request, ereq *gock.Request) (bool, error) {
			if strings.Contains(req.URL.String(), "/v"+api.version+"/build") {
				return true, nil
			}
			return false, nil
		}).
		Reply(200).
		JSON(map[string]string{
			"stream": "Step 1/5...",
		})

	serviceName := "mock-service-abc"
	project := types.Project{
		Name:     serviceName,
		Language: "node",
		Files: []types.ProjectSourceFile{
			types.ProjectSourceFile{
				Path: "Dockerfile",
				Body: `
FROM metrue/fx-node-base

COPY . .
EXPOSE 3000
CMD ["node", "app.js"]`,
				IsHandler: false,
			},
			types.ProjectSourceFile{
				Path: "app.js",
				Body: `
const Koa = require('koa');
const bodyParser = require('koa-bodyparser');
const func = require('./fx');

const app = new Koa();
app.use(bodyParser());
app.use(ctx => {
  const msg = func(ctx.request.body);
  ctx.body = msg;
});

app.listen(3000);`,
				IsHandler: false,
			},
			types.ProjectSourceFile{
				Path: "fx.js",
				Body: `
module.exports = (input) => {
    return input.a + input.b
}
					`,
				IsHandler: true,
			},
		},
	}

	service, err := api.Build(project)
	if err != nil {
		t.Fatal(err)
	}
	if service.Name != serviceName {
		t.Fatalf("should get %s but got %s", serviceName, service.Name)
	}
	if service.Image == "" {
		t.Fatal("service image should not be empty")
	}
}
