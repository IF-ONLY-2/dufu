package device

var devices = map[string]Device{}

type Device interface{
	Name() string
}

// Register registers a device to global device list
func Register(device Device) {
	devices[device.Name()] = device
}
