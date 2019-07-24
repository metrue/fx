package provision

import "testing"

func TestStart(t *testing.T) {
	opt := Options{
		Host:     "13.124.202.227",
		User:     "root",
		Password: "hithit",
	}
	provisionor := New(opt)
	if err := provisionor.Start(); err != nil {
		t.Fatal(err)
	}
}
