package message

import "testing"

func TestCreateDownMessage(t *testing.T) {
	a := DownMsgMeta{
		containerId:     "id1",
		containerStatus: "stopped",
		imageStatus:     "removed",
	}
	b := DownMsgMeta{
		containerId:     "id2",
		containerStatus: "stopped",
		imageStatus:     "removed",
	}

	downs := []DownMsgMeta{a, b}
	msg := CreateDownMessage(downs)
	expectMsg := `------------------------------------------------------
ID		ServiceStatus		ResourceStatus
id1		stopped			removed
id2		stopped			removed
------------------------------------------------------`
	if msg != expectMsg {
		t.Errorf("expect: \n%s, \ngot: \n%s", expectMsg, msg)
	}
}
