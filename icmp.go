package dufu

var ICMPv4ProtocolNumber uint8 = 1

// ICMPv4 represents an ICMPv4 header stored in a byte array.
type ICMPv4 []byte

// Typical values of ICMPv4Type defined in RFC 792.
const (
	ICMPv4EchoReply      uint8 = 0
	ICMPv4DstUnreachable uint8 = 3
	ICMPv4SrcQuench      uint8 = 4
	ICMPv4Redirect       uint8 = 5
	ICMPv4Echo           uint8 = 8
	ICMPv4TimeExceeded   uint8 = 11
	ICMPv4ParamProblem   uint8 = 12
	ICMPv4Timestamp      uint8 = 13
	ICMPv4TimestampReply uint8 = 14
	ICMPv4InfoRequest    uint8 = 15
	ICMPv4InfoReply      uint8 = 16
)

func (b ICMPv4) Type() uint8        { return uint8(b[0]) }
func (b ICMPv4) SetType(typ uint8)  { b[0] = typ; return }
func (b ICMPv4) Code() uint8        { return uint8(b[1]) }
func (b ICMPv4) SetCode(code uint8) { b[1] = code }
func (b ICMPv4) Checksum() []byte   { return b[2:4] }

// An Echo represents an ICMP echo request or reply message body.
type Echo []byte

func (e Echo) ID() []byte   { return e[:2] }
func (e Echo) Seq() []byte  { return e[2:4] }
func (e Echo) Data() []byte { return e[4:] }

// TODO: handle ICMP message
func ICMPRcv(_ *L3Layer, packet IPv4) {
	// ether + ip + icmp header
	request := ICMPv4(packet.Message())
	if request.Type() == ICMPv4Echo {
		request.SetType(ICMPv4EchoReply)
		// 可以直接用 echoRequest 回，除了 Type 不一样 echo 的 body 是一样的。
		// TODO: 看看协议栈是怎么处理的
		tmp := make([]byte, 4)
		copy(tmp, packet.DestinationAddress())
		copy(packet.DestinationAddress(), packet.SourceAddress())
		copy(packet.SourceAddress(), tmp)
		ICMPSend(request)
	}
	return
}

func ICMPSend(packet ICMPv4) {
	/*
	   	dpdk 中 ICMP checksum 的计算方式
	           icmp_h->icmp_type = IP_ICMP_ECHO_REPLY;
	           cksum = ~icmp_h->icmp_cksum & 0xffff;
	           cksum += ~htons(IP_ICMP_ECHO_REQUEST << 8) & 0xffff;
	           cksum += htons(IP_ICMP_ECHO_REPLY << 8);
	           cksum = (cksum & 0xffff) + (cksum >> 16);
	           cksum = (cksum & 0xffff) + (cksum >> 16);
	           icmp_h->icmp_cksum = ~cksum;
	*/
}
