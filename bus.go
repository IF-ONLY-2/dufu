package dufu

type Bus interface {
	Name() string
	Probe()
	Scan()
}

var buses = map[string]Bus{}

func RegisterBus(bus Bus) {
	buses[bus.Name()] = bus
}

func init() {
	Register(NewPCIBus())
}
