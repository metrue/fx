package message

import "fmt"

const TYPE_DOWN = "TYPE_DOWN"

func NewMessage(msgType string, msgBody string) {
	if msgType == TYPE_DOWN {

	}
}

type DownMsgMeta struct {
	ContainerId     string
	ContainerStatus string
	ImageStatus     string
}

//----------------------------
// ID	Contain 	Image
// id 	Stopped 	Removed
//----------------------------
func CreateDownMessage(downs []DownMsgMeta) (msg string) {
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

type UpMsgMeta struct {
	FunctionSource string
	LocalAddress   string
	RemoteAddress  string
}

func CreateUpMessage(ups []UpMsgMeta) (msg string) {
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
