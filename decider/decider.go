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

func NewDecider() {
	return &Decider()
}

func (decider *Decider) GetUpTargetType(string arg) int {
	if utils.IsValidFile(arg) {
		if path.Base(arg) == "Dockerfile" || path.Base(arg) == "docerfile" {
			return UP_TARGET_IS_DOCKER_FILE
		} else {
			return UP_TARGET_IS_FUNCTION
		}
	} else if utils.IsValidDir(arg) {
		if utils.IsValidFile(path.join(arg, "Dockerfile")) || utils.IsValidFile(path.join(arg, "dockerfile")) {
			return UP_TARGET_IS_DOCKER_FILE
		}
	} else if utils.isValidContainer(arg) {
		return UP_TARGET_IS_DOCKER_CONTAINER
	}
	return UP_TARGET_IS_UNKNOWN
}
