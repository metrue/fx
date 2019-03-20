package api

import (
	"testing"

	"github.com/metrue/fx/types"
)

func TestDockerHTTP(t *testing.T) {
	dockerRemoteAPI := "http://127.0.0.1:1234"
	version, err := Version(dockerRemoteAPI)
	if err != nil {
		t.Fatal(err)
	}
	api := NewWithDockerRemoteAPI(dockerRemoteAPI, version)

	serviceName := "a-test-service"
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

	if err != nil {
		t.Fatal(err)
	}
	if service.Status != types.ServiceStatusInit {
		t.Fatalf("should get %d but got %d", types.ServiceStatusInit, service.Status)
	}
	if service.Name != serviceName {
		t.Fatalf("should get %s but got %s", serviceName, service.Name)
	}

	if err := api.Run(&service); err != nil {
		t.Fatal(err)
	}

	services, err := api.list(serviceName)
	if err != nil {
		t.Fatal(err)
	}
	if len(services) != 1 {
		t.Fatal("service number should be 1")
	}

	if err := api.Stop(serviceName); err != nil {
		t.Fatal(err)
	}
}
