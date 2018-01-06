package decider

import (
	"path"

	"github.com/metrue/fx/utils"
)

const (
	UP_TARGET_IS_UNKNOWN          = -1
	UP_TARGET_IS_FUNCTION         = 0
	UP_TARGET_IS_DOCKER_FILE      = 1
	UP_TARGET_IS_DOCKER_CONTAINER = 2
)

type Decider struct {
}

func NewDecider() *Decider {
	return &Decider{}
}

func isDockerfileInDir(dir string) (bool, error) {
	isDir, err := utils.IsValidDir(dir)
	if err == nil && isDir {
		isDockerfile, err := utils.IsValidFile(path.Join(dir, "Dockerfile"))
		if err == nil && isDockerfile {
			return true, err
		}

		isDockerfile, err = utils.IsValidFile(path.Join(dir, "dockerfile"))
		if err == nil && isDockerfile {
			return true, err
		}
	}
	return false, err
}

func (decider *Decider) GetUpTargetType(arg string) int {
	isFile, err := utils.IsValidFile(arg)
	if err == nil && isFile {
		if path.Base(arg) == "Dockerfile" || path.Base(arg) == "docerfile" {
			return UP_TARGET_IS_DOCKER_FILE
		} else {
			return UP_TARGET_IS_FUNCTION
		}
	}

	isDockerfile, err := isDockerfileInDir(arg)
	if err == nil && isDockerfile {
		return UP_TARGET_IS_DOCKER_FILE
	}

	isContainer, err := utils.IsValidDockerContainer(arg)
	if err == nil && isContainer {
		return UP_TARGET_IS_DOCKER_CONTAINER
	}

	return UP_TARGET_IS_UNKNOWN
}
