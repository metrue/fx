package constants

// DockerRemoteAPIEndpoint docker remote api
const DockerRemoteAPIEndpoint = "13.124.202.227:1234"

// AgentPort fx server agent port
const AgentPort = 8866

// BaseImages base images to build fx functions
var BaseImages = []string{
	"metrue/fx-java-base",
	"metrue/fx-julia-base",
	"metrue/fx-python-base",
	"metrue/fx-node-base",
	"metrue/fx-d-base",
	"metrue/fx-go-base",
}
