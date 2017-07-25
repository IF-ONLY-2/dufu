package main

import (
	"testing"
)

func TestPCIRead(t *testing.T) {
	parseResource("/sys/class/uio/uio0/device/resource")
}
