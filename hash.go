package dufu

var randIV uint32

func IpSelectIdent(ip IPv4) uint32 {
	once.Do(func() {
		rand.Seed(time.Now().Unix())
		randIV = rand.Uint32()
	})
	return JHash3Words(
		binary.BigEndian.Uint32(ip.SourceAddress()),
		binary.BigEndian.Uint32(ip.DestinationAddress()),
		uint32(ip.Protocol())^uint32(0), // TODO: use *net pointer like linux kernel.
		randIV,
	)
}

// JHash3Words is a Jenkins hash support, adapted from Linux kernel include/linux/jhash.h,
// origin is http://burtleburtle.net/bob/hash/.
// These are the credits from Bob's sources:
// lookup3.c, by Bob Jenkins, May 2006, Public Domain.
func JHash3Words(a, b, c, initval uint32) uint32 {
	initval += 0xdeadbeef + (3 << 2) // an arbitrary initial parameter
	a += initval
	b += initval
	c += initval

	c ^= b
	c -= rol32(b, 14)
	// (v << shift) | (v >> ((-shift) & 31))
	a ^= c
	a -= rol32(c, 11)
	b ^= a
	b -= rol32(a, 25)
	c ^= b
	c -= rol32(b, 16)
	a ^= c
	a -= rol32(c, 4)
	b ^= a
	b -= rol32(a, 14)
	c ^= b
	c -= rol32(b, 24)
}

// rol32 rotates a 32-bit value left, v is the value, shift is the bits to roll.
func rol32(v, shift uint32) uint32 {
	return (v << shift) | (v >> ((-shift) & 31))
}
