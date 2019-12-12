package infra

import (
	"fmt"
)

// TODO upgrade to latest when k3s fix the tls scan issue
// https://github.com/rancher/k3s/issues/556
const k3sVersion = "v0.9.1"

// Scripts to provision host
var Scripts = map[string]interface{}{
	"docker_version":   "docker version",
	"install_docker":   "curl -fsSL https://download.docker.com/linux/static/stable/x86_64/docker-18.06.3-ce.tgz -o docker.tgz && tar zxvf docker.tgz && mv docker/* /usr/bin && rm -rf docker docker.tgz",
	"start_dockerd":    "dockerd >/dev/null 2>&1 & sleep 2",
	"check_fx_agent":   "docker inspect fx-agent",
	"start_fx_agent":   "docker run -d --name=fx-agent --rm -v /var/run/docker.sock:/var/run/docker.sock -p 0.0.0.0:8866:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock",
	"check_k3s_server": "ps aux | grep 'k3s server --docker'",
	"setup_k3s_master": func(ip string) string {
		return fmt.Sprintf("curl -sLS https://get.k3s.io | INSTALL_K3S_EXEC='server --docker --tls-san %s' INSTALL_K3S_VERSION='%s' sh -", ip, k3sVersion)
	},
	"check_k3s_agent": "ps aux | grep 'k3s agent --docker'",
	"setup_k3s_agent": func(masterURL string, tok string) string {
		return fmt.Sprintf("curl -fL https://get.k3s.io/ | K3S_URL='%s' K3S_TOKEN='%s' INSTALL_K3S_VERSION='%s' sh -s - --docker", masterURL, tok, k3sVersion)
	},
	"get_k3s_token":      "cat /var/lib/rancher/k3s/server/node-token",
	"get_k3s_kubeconfig": "cat /etc/rancher/k3s/k3s.yaml",
}
