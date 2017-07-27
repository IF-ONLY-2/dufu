package bus

import (
	"testing"
)

func TestPCIBusScan(t *testing.T) {
	for _, bus := range buses {
		bus.Scan()
	}
}
