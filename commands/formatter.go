package commands

import (
	"fmt"

	"github.com/metrue/fx/api"
)

//DownMessage format a message for the function down command
func DownMessage(downs []*api.DownMsgMeta) (msg string) {
	msgPrefix := `------------------------------------------------------
ID			ServiceStatus		ResourceStatus		Error`

	msg += msgPrefix
	for _, down := range downs {
		msg += fmt.Sprintf("\n%s\t\t%s\t\t\t%s\t\t\t\t%s", down.ContainerId, down.ContainerStatus, down.ImageStatus, down.Error)
	}
	msgSuffix := "\n------------------------------------------------------\n"
	msg += msgSuffix
	return msg
}

//UpMessage format a message for the function up command
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

//ListMessage format a message for the function list command
func ListMessage(containers []*api.ListItem) (msg string) {

	format := "%-15s\t%-10s\t%s\n"
	msg = fmt.Sprintf(format, "Function ID", "State", "Service URL")
	msg += "-----------------------------------------------------------------\n"
	for _, container := range containers {
		msg += fmt.Sprintf(format, container.FunctionID, container.State, container.ServiceURL)
	}

	return msg
}
