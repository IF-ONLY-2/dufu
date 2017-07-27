package bus

type Device struct {
}

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
}
func (p PCIDevice)Name()string{
return p.name
}
