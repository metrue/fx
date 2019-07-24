package command

import "testing"

func TestLocalCommand(t *testing.T) {
	localRunner := &LocalRunner{}

	cases := []struct {
		name   string
		script string
	}{
		{
			name:   "get who am i",
			script: "whoami",
		},
		{
			name:   "check docker version",
			script: "docker version",
		},
	}

	for _, c := range cases {
		cmd := New(c.name, c.script, localRunner)
		out, err := cmd.Exec()
		if err != nil {
			t.Fatal(err)
		}
		if len(out) == 0 {
			t.Fatal("whoami should not get empty")
		}
	}
}

func TestRemoteCommand(t *testing.T) {
	// TODO
}
