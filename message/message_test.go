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
