package common

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/metrue/fx/api"
	"github.com/olekukonko/tablewriter"
)

func HandleDownResult(downs []*api.DownMsgMeta) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Container Status", "Image Status", "Error"})
	for _, down := range downs {
		table.Append([]string{
			down.ContainerId,
			down.ContainerStatus,
			down.ImageStatus,
			down.Error})
	}
	table.Render()
}

func HandleUpResult(ups []*api.UpMsgMeta) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"local url", "external url"})
	for _, up := range ups {
		table.Append([]string{
			up.LocalAddress,
			up.RemoteAddress})
	}
	table.Render()
}

func HandleListResult(containers []*api.ListItem) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"id", "state", "url"})
	for _, container := range containers {
		table.Append([]string{
			container.FunctionID,
			container.State,
			container.ServiceURL})
	}
	table.Render()
}

func HandleCallResult(res interface{}) {
	msg, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("Could not make %v to json string", err)
	}
	fmt.Println(string(msg))
}
