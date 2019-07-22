workflow "Release" {
  on = "push"
  resolves = [
    "goreleaser",
    "notify"
  ]
}

action "tag?" {
  uses = "actions/bin/filter@master"
  args = "tag"
}

action "goreleaser" {
  uses = "docker://goreleaser/goreleaser"
  secrets = [
    "GORELEASER_GITHUB_TOKEN",
    "DOCKER_USERNAME",
    "DOCKER_PASSWORD",
  ]
  args = "release"
  needs = ["tag?"]
}

action "notify" {
  needs = [
    "goreleaser"
  ]
  uses = "metrue/noticeme-github-action@master"
  secrets = ["NOTICE_ME_TOKEN"]
  args = ["fx release ok"]
}
