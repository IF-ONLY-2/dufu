package e1000

// e1000 read register
func readReg(hw []byte, reg int) uint32 {
	return uint32(hw[reg]) + uint32(hw[reg+1]<<8) + uint32(hw[reg+2]<<16) + uint32(hw[reg+3]<<24)
}

// e1000 write register
func writeReg(hw []uint8, reg int, val byte) {
	hw[reg] = byte(val)
	hw[reg+1] = byte(val >> 8)
	hw[reg+2] = byte(val >> 16)
	hw[reg+3] = byte(val >> 24)
}
