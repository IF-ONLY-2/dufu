package dufu

const (
	UDPHeaderLength uint8 = 8 // udp header length

	IPHeaderLength uint8 = 20 // normallly ip header length is 20
)

// UDP
type UDP []byte

func SendUDP(skb *SkBuff) {
	skb.Prepend(IPHeaderLength)
}
