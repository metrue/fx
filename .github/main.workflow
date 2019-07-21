workflow "build and push to dockerhub" {
  on = "push"
  resolves = [
    "lint",
    "login",
    "build-fx-node-image",
    "push-fx-node-image",
    "build-fx-rust-image",
    "push-fx-rust-image",
    "build-fx-go-image",
    "push-fx-go-image",
    "notify"
  ]
}

action "login" {
  uses = "actions/docker/login@8cdf801b322af5f369e00d85e9cf3a7122f49108"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "build-fx-node-image" {
  uses = "actions/docker/cli@master"
  args = "build -t metrue/fx-node-base:latest -f api/asserts/dockerfiles/base/node/Dockerfile api/asserts/dockerfiles/base/node"
}

action "push-fx-node-image" {
  needs = ["build-fx-node-image", "login"]
  uses = "actions/docker/cli@master"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
  args = "push metrue/fx-node-base:latest"
}

action "build-fx-go-image" {
  uses = "actions/docker/cli@master"
  args = "build -t metrue/fx-go-base:latest -f api/asserts/dockerfiles/base/go/Dockerfile api/asserts/dockerfiles/base/go"
}

action "push-fx-go-image" {
  needs = ["build-fx-go-image", "login"]
  uses = "actions/docker/cli@master"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
  args = "push metrue/fx-go-base:latest"
}

action "build-fx-rust-image" {
  uses = "actions/docker/cli@master"
  args = "build -t metrue/fx-rust-base:latest -f api/asserts/dockerfiles/base/rust/Dockerfile api/asserts/dockerfiles/base/rust"
}

action "push-fx-rust-image" {
  needs = ["build-fx-rust-image", "login"]
  uses = "actions/docker/cli@master"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
  args = "push metrue/fx-rust-base:latest"
}

action "lint" {
  uses = "actions-contrib/golangci-lint@master"
  args = "run"
}

action "notify" {
  needs = [
    "push-fx-node-image",
    "push-fx-rust-image",
    "push-fx-go-image"
  ]
  uses = "metrue/noticeme-github-action@master"
  secrets = ["NOTICE_ME_TOKEN"]
  args = ["BuildFxGitHubActionDone"]
}
