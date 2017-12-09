package message

import "testing"

func TestCreateDownMessage(t *testing.T) {
	a := DownMsgMeta{
		ContainerId:     "id1",
		ContainerStatus: "stopped",
		ImageStatus:     "removed",
	}
	b := DownMsgMeta{
		ContainerId:     "id2",
		ContainerStatus: "stopped",
		ImageStatus:     "removed",
	}

	downs := []DownMsgMeta{a, b}
	msg := CreateDownMessage(downs)
	expectMsg := `------------------------------------------------------
ID			ServiceStatus		ResourceStatus
id1		stopped			removed
id2		stopped			removed
------------------------------------------------------`
	if msg != expectMsg {
		t.Errorf("expect: \n%s, \ngot: \n%s", expectMsg, msg)
	}
}

func TestCreateUpMessage(t *testing.T) {
	a := UpMsgMeta{
		FunctionSource: "a.js",
		LocalAddress:   "127.0.0.1:9090",
		RemoteAddress:  "156.23.23.1:9090",
	}
	b := UpMsgMeta{
		FunctionSource: "b.js",
		LocalAddress:   "127.0.0.1:9090",
		RemoteAddress:  "156.23.23.1:9090",
	}

	ups := []UpMsgMeta{a, b}
	msg := CreateUpMessage(ups)
	expectMsg := `-----------------------------------------------------------------
FunctionSource				LocalAddress			RemoteAddress
a.js		127.0.0.1:9090			156.23.23.1:9090
b.js		127.0.0.1:9090			156.23.23.1:9090
-----------------------------------------------------------------`
	if msg != expectMsg {
		t.Errorf("expect: \n%s, \ngot: \n%s", expectMsg, msg)
	}
}
