package api

// func TestDockerHTTP(t *testing.T) {
// 	const addr = "127.0.0.1"
// 	const user = ""
// 	const passord = ""
// 	provisioner := provision.NewWithHost(addr, user, passord)
// 	if err := utils.RunWithRetry(func() error {
// 		if !provisioner.IsFxAgentRunning() {
// 			if err := provisioner.StartFxAgent(); err != nil {
// 				log.Infof("could not start fx agent on host: %s", err)
// 				return err
// 			}
// 			log.Infof("fx agent started")
// 		} else {
// 			log.Infof("fx agent is running")
// 		}
// 		return nil
// 	}, 2*time.Second, 10); err != nil {
// 		t.Fatal(err)
// 	} else {
// 		defer provisioner.StopFxAgent()
// 	}
//
// 	host := config.Host{Host: "127.0.0.1"}
// 	api, err := Create(host.Host, constants.AgentPort)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	serviceName := "a-test-service"
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
//
// 	service, err := api.Build(project)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if service.Name != serviceName {
// 		t.Fatalf("should get %s but got %s", serviceName, service.Name)
// 	}
//
// 	if err := api.Run(9999, &service); err != nil {
// 		t.Fatal(err)
// 	}
//
// 	services, err := api.ListContainer(serviceName)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(services) != 1 {
// 		t.Fatal("service number should be 1")
// 	}
//
// 	if err := api.Stop(serviceName); err != nil {
// 		t.Fatal(err)
// 	}
//
// 	const network = "fx-net"
// 	if err := api.CreateNetwork(network); err != nil {
// 		t.Fatal(err)
// 	}
//
// 	nws, err := api.GetNetwork(network)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if nws[0].Name != network {
// 		t.Fatalf("should get %s but got %s", network, nws[0].Name)
// 	}
// }
