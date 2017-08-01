package dufu

import (
	"testing"
	"time"
)

func TestTap(t *testing.T) {
	f, err := NewTAP("test")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	time.Sleep(time.Duration(100) * time.Second)
}
