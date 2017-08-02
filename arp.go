package dufu

import(
	"fmt"
)

const (
	// ARPProtocolNumber is the ARP network protocol number.
	// TODO: make a protocol interface.
	ARPProtocolNumber = 0x0806

	// ARPSize is the size of an IPv4-over-Ethernet ARP packet.
	ARPSize = 2 + 2 + 1 + 1 + 2 + 2*6 + 2*4

	// Typical ARP opcodes defined in RFC 826.
	ARPRequest uint16 = 1
	ARPReply   uint16 = 2
)

var ARPCache map[[6]byte][4]byte

// ARP represents the arp packet.
type ARP []byte

func (a ARP) InitIPv4OverEthernetARPPacket(op uint16) {
	a[0], a[1] = 0, 1       // hardware type of ethernet
	a[2], a[3] = 0x08, 0x00 // IPv4 protocol number
	a[4] = 6                // MAC address size
	a[5] = 4                // IPv4 address size
	a[6] = uint8(op >> 8)
	a[7] = uint8(op)
}

// HardwareType is the l2 hardware type, often ethernet (1).
func (a ARP) HardwareType() uint16 { return uint16(a[0])<<8 | uint16(a[1]) }

// ProtocolType is the type for l3 network type, often IPv4 (0x0800).
func (a ARP) ProtocolType() uint16 { return uint16(a[2])<<8 | uint16(a[3]) }

// HardwareAddressSize specifies the size of l2 hardware address.
func (a ARP) HardwareAddressSize() int { return int(a[4]) }

// ProtocolAddressSize specifies the size of l3 network address.
func (a ARP) ProtocolAddressSize() int { return int(a[5]) }

// Op is the ARP opcode.
func (a ARP) Op() uint16 { return uint16(a[6])<<8 | uint16(a[7]) }

// SenderHardwareAddress is the link address of the sender.
// It is a view on to the ARP packet so it can be used to set the value.
func (a ARP) SenderHardwareAddress() []byte { const s = 8; return a[s : s+6] }

// SenderProtocolAddress is the protocol address of the sender.
// It is a view on to the ARP packet so it can be used to set the value.
func (a ARP) SenderProtocolAddress() []byte { const s = 8 + 6; return a[s : s+4] }

// TargetHardwareAddress is the link address of the target.
// It is a view on to the ARP packet so it can be used to set the value.
func (a ARP) TargetHardwareAddress() []byte { const s = 8 + 6 + 4; return a[s : s+6] }

// TargetProtocolAddress is the protocol address of the target.
// It is a view on to the ARP packet so it can be used to set the value.
func (a ARP) TargetProtocolAddress() []byte { const s = 8 + 6 + 4 + 6; return a[s : s+4] }

func ARPHandle(l2l *L2Layer, packet []byte) {
	request := ARP(packet)
	if request.Op() == ARPRequest {
		buf := make([]byte, 14+ARPSize)
		reply := ARP(buf[14:])
		reply.InitIPv4OverEthernetARPPacket(ARPReply)
		copy(reply.SenderHardwareAddress(), l2l.HardwareAddr[:])
		copy(reply.SenderProtocolAddress(), request.TargetProtocolAddress())
		copy(reply.TargetHardwareAddress(), request.SenderHardwareAddress())
		copy(reply.TargetProtocolAddress(), request.SenderProtocolAddress())

		frame := Frame(buf)
		copy(frame.Destination(), request.SenderHardwareAddress()[:])
		copy(frame.Source(), l2l.HardwareAddr[:])
		copy(frame.EtherType(), []byte{0x08, 0x06})
		fmt.Println(buf)
		for _,b:=range buf{
			fmt.Printf("%.2x ",b)
		}
		fmt.Println("")
		go l2l.Send(Frame(buf))
	}

	if request.Op() == ARPReply {
		// TODO: address caching system
		// ARPCache[net.HardwareAddress(a.SenderHardwareAddress())] = net.IP(a.SenderProtocolAddress())
	}
}
