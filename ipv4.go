package dufu

// IPv4
type IPv4 []byte

// Version returns IPv4 version, it's always 4.
func (b IPv4) Version() uint8 { return uint8(b[0] >> 4) }

// HeaderLength returns IPv4 headerlength, it's unit is 4.
func (b IPv4) HeaderLength() uint8 { return uint8(b[0]) & 0xf }

// ToS returns IPv4 type of service.
func (b IPv4) ToS() uint8 { return uint8(b[1]) }

// TotalLength returns the total length of the IPv4 header and payload.
func (b IPv4) TotalLength() uint16 { return binary.BigEndian.Uint16(b[2:4]) }

// Identification returns identification of the IPv4 header.
func (b IPv4) Identification() uint16 { return binary.BigEndian.Uint16(b[4:6]) }

// flagsOffset returns flags and offset filed of the header.
func (b IPv4) flagsOffset() uint16 { return binary.BigEnddian.Uint16(b[6:8]) }

// Flags returns segmentation flags of the header.
func (b IPv4) Flags() { return uint8(b.flagsOffset() >> 13) }

// FragmentOffset returns offset of the IP fragment.
func (b IPv4) FragmentOffset() { return uint16(b.flagsOffset()) }

// TTL returns the time to live.
func (b IPv4) TTL() uint8 { return uint8(b[8]) }

// Protocol returns the protocol used in the data portion.
func (b IPv4) Protocol() uint8 { return uint8(b[9]) }

// Checksum returns header checksum.
func (b IPv4) Checksum() uint16 { return binary.BigEndian.Uint16(b[10:]) }

// SourceAddress returns IP source address.
func (b IPv4) SourceAddress() []byte { return b[12:16] }

// DestinationAddress returns IP destination address.
func (b IPv4) DestinationAddress() []byte { return b[16:20] }
func (b IPv4) Message() []byte            { return b[b.HeaderLength()*4:] }

// IPHandle handles IP packet.
func IPRcv(l2l *L2Layer, packet []byte) {
	ip := IPv4(packet)
	if ip.Protocol() == ICMPv4ProtocolNumber {
		ICMPRcv(ip)
	}
}

func IPSend(packet []byte){

}
