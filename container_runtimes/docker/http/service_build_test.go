package api

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"testing"
//
// 	"github.com/metrue/fx/types"
// )
//
// func TestMakeTar(t *testing.T) {
// 	serviceName := "mock-service-abc"
// 	project := types.Project{
// 		Name:     serviceName,
// 		Language: "node",
// 		Files: []types.ProjectSourceFile{
// 			types.ProjectSourceFile{
// 				Path: "Dockerfile",
// 				Body: `
// FROM metrue/fx-node-base
//
// COPY . .
// EXPOSE 3000
// CMD ["node", "app.js"]`,
// 				IsHandler: false,
// 			},
// 			types.ProjectSourceFile{
// 				Path: "app.js",
// 				Body: `
// const Koa = require('koa');
// const bodyParser = require('koa-bodyparser');
// const func = require('./fx');
//
// const app = new Koa();
// app.use(bodyParser());
// app.use(ctx => {
//   const msg = func(ctx.request.body);
//   ctx.body = msg;
// });
//
// app.listen(3000);`,
// 				IsHandler: false,
// 			},
// 			types.ProjectSourceFile{
// 				Path: "fx.js",
// 				Body: `
// module.exports = (input) => {
//     return input.a + input.b
// }
// 					`,
// 				IsHandler: true,
// 			},
// 		},
// 	}
// 	tarDir, err := ioutil.TempDir("/tmp", "fx-tar")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.RemoveAll(tarDir)
//
// 	tarFilePath := filepath.Join(tarDir, fmt.Sprintf("%s.tar", serviceName))
// 	if err := makeTar(project, tarFilePath); err != nil {
// 		t.Fatal(err)
// 	}
//
// 	file, err := os.Open(tarFilePath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	stat, err := file.Stat()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if stat.Name() != serviceName+".tar" {
// 		t.Fatalf("should get %s but got %s", serviceName+".tar", stat.Name())
// 	}
// 	if stat.Size() <= 0 {
// 		t.Fatalf("tarfile invalid: size %d", stat.Size())
// 	}
// }
