package k8s

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/metrue/fx/infra"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
)

// Cloud define a cloud
type Cloud struct {
	//  Define where is the location of kubeconf would be saved to
	KubeConfig string           `json:"config"`
	Type       string           `json:"type"`
	Nodes      map[string]Noder `json:"nodes"`

	token string
	url   string
}

// Load a cloud from config
func Load(meta []byte, messup ...func(n Noder) (Noder, error)) (*Cloud, error) {
	var cloud Cloud
	if err := json.Unmarshal(meta, &cloud); err != nil {
		return nil, err
	}

	for name, n := range cloud.Nodes {
		// NOTE messup function is just for unit testing
		// we use it to replace the real not with mock node
		if len(messup) > 0 {
			node, err := messup[0](n)
			if err != nil {
				return nil, err
			}
			cloud.Nodes[name] = node
		}
	}
	return &cloud, nil
}

// NewCloud new a cloud
func NewCloud(kubeconf string, node ...Noder) *Cloud {
	nodes := map[string]Noder{}
	for _, n := range node {
		nodes[n.GetName()] = n
	}

	return &Cloud{
		KubeConfig: kubeconf,
		Type:       types.CloudTypeK8S,
		Nodes:      nodes,
	}
}

// Provision provision cloud
func (c *Cloud) Provision() error {
	var master Noder
	agents := []Noder{}
	for _, n := range c.Nodes {
		if n.GetType() == NodeTypeMaster {
			master = n
		} else {
			agents = append(agents, n)
		}
	}

	// when it's k3s cluster
	if master != nil {
		c.url = fmt.Sprintf("https://%s:6443", master.GetIP())
		if err := master.Provision(map[string]string{}); err != nil {
			return err
		}

		tok, err := master.GetToken()
		if err != nil {
			return err
		}
		c.token = tok

		config, err := master.GetConfig()
		if err != nil {
			return err
		}

		if err := utils.EnsureFile(c.KubeConfig); err != nil {
			return err
		}
		if err := ioutil.WriteFile(c.KubeConfig, []byte(config), 0666); err != nil {
			return err
		}
	}

	if len(agents) > 0 {
		errCh := make(chan error, len(agents))
		defer close(errCh)

		for _, agent := range agents {
			go func(node Noder) {
				errCh <- node.Provision(map[string]string{
					"url":   c.url,
					"token": c.token,
				})
			}(agent)
		}

		for range agents {
			err := <-errCh
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// AddNode a node
func (c *Cloud) AddNode(n Noder, skipProvision bool) error {
	if !skipProvision {
		if err := n.Provision(map[string]string{
			"url":   c.url,
			"token": c.token,
		}); err != nil {
			return err
		}
	}

	c.Nodes[n.GetName()] = n
	return nil
}

// DeleteNode a node
func (c *Cloud) DeleteNode(name string) error {
	node, ok := c.Nodes[name]
	if ok {
		delete(c.Nodes, name)
	}
	if node.GetType() == NodeTypeMaster && len(c.Nodes) > 0 {
		return fmt.Errorf("could not delete master node since there is still agent node running")
	}
	return nil
}

// State get cloud state
func (c *Cloud) State() {}

// UnmarshalJSON unmarsha json
func (c *Cloud) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	c.Nodes = make(map[string]Noder)

	for k, v := range m {
		if k == "nodes" {
			nodes, ok := v.(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid nodes data")
			}
			for name, n := range nodes {
				node, ok := n.(map[string]interface{})
				if !ok {
					return fmt.Errorf("invalid node data")
				}
				n, err := CreateNode(node["ip"].(string), node["user"].(string), node["type"].(string), node["name"].(string))
				if err != nil {
					return err
				}
				c.Nodes[name] = n
			}
		} else if k == "token" {
			tok, ok := v.(string)
			if ok {
				c.token = tok
			} else {
				c.token = ""
			}
		} else if k == "config" {
			config, ok := v.(string)
			if ok {
				c.KubeConfig = config
			} else {
				c.KubeConfig = ""
			}
		} else if k == "type" {
			typ, ok := v.(string)
			if ok {
				c.Type = typ
			} else {
				c.Type = ""
			}
		} else if k == "url" {
			url, ok := v.(string)
			if ok {
				c.url = url
			} else {
				c.url = ""
			}
		}
	}

	return nil
}

// MarshalJSON cloud information
func (c *Cloud) MarshalJSON() ([]byte, error) {
	nodes := map[string]Node{}
	for name, node := range c.Nodes {
		nodes[name] = Node{
			IP:   node.GetIP(),
			Type: node.GetType(),
			User: node.GetUser(),
			Name: node.GetName(),
		}
	}

	body, err := json.Marshal(struct {
		URL        string          `json:"url"`
		KubeConfig string          `json:"config"`
		Type       string          `json:"type"`
		Token      string          `json:"token"`
		Nodes      map[string]Node `json:"nodes"`
	}{
		KubeConfig: c.KubeConfig,
		Type:       c.Type,
		Token:      c.token,
		URL:        c.url,
		Nodes:      nodes,
	})
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetType get type of cloud
func (c *Cloud) GetType() string {
	return c.Type
}

// Dump cloud data
func (c *Cloud) Dump() ([]byte, error) {
	return json.Marshal(c)
}

// GetConfig get config
func (c *Cloud) GetConfig() (string, error) {
	if c.KubeConfig != "" {
		return c.KubeConfig, nil
	}
	if err := c.Provision(); err != nil {
		return "", err
	}
	return c.KubeConfig, nil
}

// IsHealth check if cloud is in health
func (c *Cloud) IsHealth() (bool, error) {
	return true, nil
}

var (
	_ infra.Clouder = &Cloud{}
)
