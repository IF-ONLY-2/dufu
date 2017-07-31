package dufu

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
type Frame struct {
	FrameHeader []byte
	Packet      []byte
}

func (l2l *L2Layer) Read() (Frame, error) {
	var (
		buf   [1526]byte
		frame Frame
	)
	n, err := l2l.Read(buf[:])
}
