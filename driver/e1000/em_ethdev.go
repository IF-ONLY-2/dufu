package e1000

import (
	"github.com/ggaaooppeenngg/dufu"
)

type EmEthDev struct {
}

func (emDev *EmEthDev) String() string {
	return "e1000_driver"
}

func (emDev *EmEthDev) Start(dev *dufu.PCIDevice) {
	return
}

func (emDev *EmEthDev) Stop(dev *dufu.PCIDevice) {
	return
}
