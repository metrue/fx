package k8s

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	mock_infra "github.com/metrue/fx/infra/k8s/mocks"
)

func TestLoad(t *testing.T) {
	t.Run("empty meta", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var createNodeFn = func(n Noder) (Noder, error) {
			return nil, nil
		}

		_, err := Load([]byte{}, createNodeFn)
		if err == nil {
			t.Fatalf("should load with error")
		}
	})

	t.Run("only master node", func(t *testing.T) {
		kubeconfig := "./kubeconfig.yml"
		defer func() {
			if err := os.RemoveAll("./kubeconfig.yml"); err != nil {
				t.Fatal(err)
			}
		}()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		master := mock_infra.NewMockNoder(ctrl)
		var createNodeFn = func(n Noder) (Noder, error) {
			return master, nil
		}
		typ := NodeTypeMaster
		name := "master"
		ip := "127.0.0.1"
		user := "testuser"
		kubeconfContent := "sample-content"
		master.EXPECT().GetName().Return(name)
		master.EXPECT().GetType().Return(typ).Times(2)
		master.EXPECT().GetIP().Return(ip).Times(2)
		master.EXPECT().GetUser().Return(user)
		master.EXPECT().GetConfig().Return(kubeconfContent, nil)

		claud := &Cloud{
			KubeConfig: kubeconfig,
			Type:       "k8s",
			url:        "",
			token:      "",
			Nodes:      map[string]Noder{"master-node": master},
		}

		meta, err := json.Marshal(claud)
		if err != nil {
			t.Fatal(err)
		}
		cloud, err := Load(meta, createNodeFn)
		if err != nil {
			t.Fatal(err)
		}
		if len(cloud.Nodes) != 1 {
			t.Fatalf("should get %d but got %d", 1, len(cloud.Nodes))
		}

		master.EXPECT().Provision(map[string]string{}).Return(nil)
		master.EXPECT().GetToken().Return("tok-1", nil)
		if err := cloud.Provision(); err != nil {
			t.Fatal(err)
		}

		content, err := ioutil.ReadFile(claud.KubeConfig)
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != kubeconfContent {
			t.Fatalf("should get %s but got %s", kubeconfContent, content)
		}
	})

	t.Run("one master node and one agent", func(t *testing.T) {
		kubeconfig := "./kubeconfig.yml"
		defer func() {
			if err := os.RemoveAll("./kubeconfig.yml"); err != nil {
				t.Fatal(err)
			}
		}()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		master := mock_infra.NewMockNoder(ctrl)
		node := mock_infra.NewMockNoder(ctrl)
		var createNodeFn = func(n Noder) (Noder, error) {
			if n.GetType() == NodeTypeMaster {
				return master, nil
			}
			return node, nil
		}
		typ := NodeTypeMaster
		name := "master"
		ip := "127.0.0.1"
		user := "testuser"
		kubeconfContent := "sample-config"
		master.EXPECT().GetName().Return(name)
		master.EXPECT().GetType().Return(typ).Times(2)
		master.EXPECT().GetIP().Return(ip).Times(3)
		master.EXPECT().GetConfig().Return(kubeconfContent, nil)
		master.EXPECT().GetUser().Return(user)

		nodeType := NodeTypeAgent
		nodeName := "agent_name"
		nodeIP := "12.12.12.12"
		nodeUser := "testuser"
		node.EXPECT().GetName().Return(nodeName)
		node.EXPECT().GetType().Return(nodeType).Times(2)
		node.EXPECT().GetIP().Return(nodeIP)
		node.EXPECT().GetUser().Return(nodeUser)

		url := fmt.Sprintf("https://%s:6443", master.GetIP())
		tok := "tok-1"
		claud := &Cloud{
			KubeConfig: kubeconfig,
			url:        url,
			token:      tok,
			Type:       "k8s",
			Nodes:      map[string]Noder{"master-node": master, "agent-node": node},
		}
		meta, err := json.Marshal(claud)
		if err != nil {
			t.Fatal(err)
		}

		cloud, err := Load(meta, createNodeFn)
		if err != nil {
			t.Fatal(err)
		}
		if len(cloud.Nodes) != 2 {
			t.Fatalf("should get %d but got %d", 2, len(cloud.Nodes))
		}

		master.EXPECT().Provision(map[string]string{}).Return(nil)
		master.EXPECT().GetToken().Return(tok, nil)
		node.EXPECT().Provision(map[string]string{
			"url":   cloud.url,
			"token": cloud.token,
		}).Return(nil)
		if err := cloud.Provision(); err != nil {
			t.Fatal(err)
		}
		content, err := ioutil.ReadFile(claud.KubeConfig)
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != kubeconfContent {
			t.Fatalf("should get %s but got %s", kubeconfContent, content)
		}
	})
}

func TestProvision(t *testing.T) {}
