package dufu

import (
	"testing"
)

func TestPCIBusScan(t *testing.T) {
	for _, bus := range buses {
		bus.Scan()
	}
	for _, bus := range buses {
		bus.Probe()
	}
}
