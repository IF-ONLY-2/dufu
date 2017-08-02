package dufu

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	MaxFrameSize      = 1526 // max frame size
	EthMACAddressSize = 6
	EthEtherTypeSize  = 2
	EthHeaderSize     = EthEtherTypeSize + 2*EthMACAddressSize
)

type L2Layer struct {
	*TapDevice
}

// Link layer frame
type Frame []byte

//  Destination is the mac destination address.
func (f Frame) Destination() []byte { return f[0:6] }

// Source is the mac source address.
func (f Frame) Source() []byte { return f[6:12] }

//  EtherType is the ethernet type.
func (f Frame) EtherType() []byte { return f[12:14] }

// Packet returns frame payload.
func (f Frame) Packet() []byte { return f[14:] }

func (l2l *L2Layer) Read() (Frame, error) {
	var (
		buf   [MaxFrameSize]byte
		frame Frame
	)

	n, err := l2l.TapDevice.Read(buf[:])
	if err != nil {
		return frame, err
	}
	frame = Frame(buf[:n])
	return frame, nil
}

func (l2l *L2Layer) Loop() {
	for {
		frame, err := l2l.Read()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("source", net.HardwareAddr(frame.Source()), "dest", net.HardwareAddr(frame.Destination()))
		fmt.Printf("0x%.4x\n", binary.BigEndian.Uint16(frame.EtherType()))
		if binary.BigEndian.Uint16(frame.EtherType()) == ARPProtocolNumber {
			ARPHandle(l2l, frame.Packet())
		}
		if binary.BigEndian.Uint16(frame.EtherType()) == 4 {
			IPHandle(l2l, frame.Packet())
		}
	}
}

func (l2l *L2Layer) Send(frame Frame) {
	l2l.TapDevice.Write([]byte(frame))
}
