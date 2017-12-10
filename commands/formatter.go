package commands

import (
	"fmt"

	"github.com/metrue/fx/api"
)

func DownMessage(downs []*api.DownMsgMeta) (msg string) {
	msgPrefix := `------------------------------------------------------
ID			ServiceStatus		ResourceStatus`

	msg += msgPrefix
	for _, down := range downs {
		msg += fmt.Sprintf("\n%s\t\t%s\t\t\t%s", down.ContainerId, down.ContainerStatus, down.ImageStatus)
	}
	msgSuffix := "\n------------------------------------------------------"
	msg += msgSuffix
	return msg
}

func UpMessage(ups []*api.UpMsgMeta) (msg string) {
	msgPrefix := `-----------------------------------------------------------------
FunctionSource				LocalAddress			RemoteAddress`

	msg += msgPrefix
	for _, up := range ups {
		msg += fmt.Sprintf("\n%s\t\t%s\t\t\t%s", up.FunctionSource, up.LocalAddress, up.RemoteAddress)
	}
	msgSuffix := "\n-----------------------------------------------------------------"
	msg += msgSuffix
	return msg
}

func ListMessage(containers []*api.ListItem) (msg string) {

	format := "%-15s\t%-10s\t%s"
	msg = fmt.Sprintf(format, "Function ID", "State", "Service URL")

	for _, container := range containers {
		msg += fmt.Sprintf(format, container.FunctionID, container.State, container.ServiceURL)
	}

	return msg
}
