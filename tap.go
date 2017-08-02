package dufu

import (
	"net"
	"os"
	"syscall"
	"unsafe"
)

const (
	IFF_TUN   = 0x0001
	IFF_TAP   = 0x0002
	IFF_NO_PI = 0x1000
)

type ifReq struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

// NewTAP creates a new tap device with the name
func NewTAP(name string) (*TapDevice, error) {
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	//  req is a helper struct, to call TUNSETIFF ioctl
	req := struct {
		Name  [0x10]byte
		Flags uint16
		pad   [0x28 - 0x10 - 2]byte
	}{}
	req.Flags = IFF_TAP | IFF_NO_PI
	copy(req.Name[:], name)
	// create tap device
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		err = errno
		return nil, err
	}
	// req1 is a helper struct, to call SIOCGIFHWADDR ioctl
	req1 := struct {
		pad1            [0x12]byte
		HardwareAddress [6]byte
		pad2            [0x28 - 0x12 - 6]byte
	}{}
	// get MAC address
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.SIOCGIFHWADDR), uintptr(unsafe.Pointer(&req1)))
	hwa := net.HardwareAddr(req1.HardwareAddress[:])
	return &TapDevice{
		File:         file,
		HardwareAddr: hwa,
	}, err
}

type TapDevice struct {
	net.HardwareAddr
	*os.File
}
