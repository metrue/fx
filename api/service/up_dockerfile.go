package service

import (
	"github.com/metrue/fx/api"
	"github.com/metrue/fx/handlers"
)

//Up deploy and run functions
func UpDockerfile(req *api.UpDockerfileRequest) (*api.UpResponse, error) {
	dockerfiles := req.Dockerfiles;
	count := len(dockerfiles)
	results := make(chan upTask, count)

	for _, dockerfile := range dockerfiles {
		go func(meta *api.DockerfileMeta) {
			results <- newUpTask(handlers.UpDockerfile(*meta))
		}(dockerfile)
	}

	// collect up results
	var ups []*api.UpMsgMeta
	for result := range results {
		upResult := result.Val
		if result.Err != nil {
			upResult = &api.UpMsgMeta{
				Error: result.Err.Error(),
			}
		}
		ups = append(ups, upResult)
		if len(ups) == count {
			close(results)
		}
	}

	return &api.UpResponse{
		Instances: ups,
	}, nil
}
