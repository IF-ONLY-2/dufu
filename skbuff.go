package dufu

import (
	"fmt"
)

// SkBuff is a buffer used for socket, it can be prepended from the head.
type SkBuff struct {
	buf   []byte
	index int
}

func NewSkBuff(size int) *SkBuff {
	return &SkBuff{
		buf:   make([]byte, size),
		index: size,
	}
}

func (skb *SkBuff) TrimFront(size int) {
	fmt.Println("call TrimFron, index", skb.index)
	skb.index += size
}

func (skb *SkBuff) Data() []byte {
	fmt.Println("call Data, index", skb.index)
	return skb.buf[skb.index:]
}

func (skb *SkBuff) Prepend(size uint8) {
	fmt.Println("call Prepend, index", skb.index)
	skb.index -= int(size)
}
