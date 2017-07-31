package dufu

type PCIDevice struct {
	PCIAddress
	Vendor          uint16
	Device          uint16
	SubsystemVendor uint16
	SubsystemDevice uint16
	Class           uint32
	MaxVFs          uint16
	NumaNode        uint64

	name string

	Driver string
	Mem    []MemResource

	HwAddress []byte
}

func (p PCIDevice) Name() string {
	return p.name
}
