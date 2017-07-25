package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const (
	PCI_COMMAND              uint16 = 0x04  /* 16 bits */
	PCI_COMMAND_IO                  = 0x1   /* Enable response in I/O space */
	PCI_COMMAND_MEMORY              = 0x2   /* Enable response in Memory space */
	PCI_COMMAND_MASTER              = 0x4   /* Enable bus mastering */
	PCI_COMMAND_SPECIAL             = 0x8   /* Enable response to special cycles */
	PCI_COMMAND_INVALIDATE          = 0x10  /* Use memory write and invalidate */
	PCI_COMMAND_VGA_PALETTE         = 0x20  /* Enable palette snooping */
	PCI_COMMAND_PARITY              = 0x40  /* Enable parity checking */
	PCI_COMMAND_WAIT                = 0x80  /* Enable address/data stepping */
	PCI_COMMAND_SERR                = 0x100 /* Enable SERR */
	PCI_COMMAND_FAST_BACK           = 0x200 /* Enable back-to-back writes */
	PCI_COMMAND_INTX_DISABLE uint16 = 0x400 /* INTx Emulation Disable */
)

const (
	// At most there 6 BARs from 0 to 5, but most devices just use 2 or 3 BARs.
	PCI_MAX_RESOURCE = 6
	// IO resource type: //
	IORESOURCE_IO  = 0x00000100
	IORESOURCE_MEM = 0x00000200
)

// I/O+ Mem+ BusMaster+ SpecCycle- MemWINV- VGASnoop- ParErr- Stepping- SERR- FastB2B- DisINTx+

func commandString(i uint16) string {
	var s string
	if i&PCI_COMMAND_IO != 0 {
		s += "I/O+ "
	} else {
		s += "I/O- "
	}
	if i&PCI_COMMAND_MEMORY != 0 {
		s += "Mem+ "
	} else {
		s += "Mem- "
	}
	if i&PCI_COMMAND_MASTER != 0 {
		s += "BusMaster+ "
	} else {
		s += "BusMaster- "
	}
	if i&PCI_COMMAND_SPECIAL != 0 {
		s += "SpecCycle+ "
	} else {
		s += "SpecCycle- "
	}
	if i&PCI_COMMAND_INVALIDATE != 0 {
		s += "MemWINV+ "
	} else {
		s += "MemWINV- "
	}

	if i&PCI_COMMAND_VGA_PALETTE != 0 {
		s += "VGASnoop+ "
	} else {
		s += "VGASnoop- "
	}
	if i&PCI_COMMAND_PARITY != 0 {
		s += "ParErr+ "
	} else {
		s += "ParErr- "
	}
	if i&PCI_COMMAND_WAIT != 0 {
		s += "Stepping+ "
	} else {
		s += "Stepping- "
	}
	if i&PCI_COMMAND_SERR != 0 {
		s += "SERR+ "
	} else {
		s += "SERR- "
	}
	if i&PCI_COMMAND_FAST_BACK != 0 {
		s += "FastB2B+ "
	} else {
		s += "FastB2B- "
	}
	if i&PCI_COMMAND_INTX_DISABLE != 0 {
		s += "DisINTx+ "
	} else {
		s += "DisINTx- "
	}
	s = s[:len(s)-1]
	return s
}

type MemResource struct {
	PhysicalAddress uint64
	Length          uint64
	Address         []byte
}

func parseResource(file string) []MemResource {
	fResource, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fResource.Close()
	var mrs []MemResource
	lineReader := bufio.NewReader(fResource)
	for i := 0; i < PCI_MAX_RESOURCE; i++ {
		line, err := lineReader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		ss := strings.Split(line, " ")
		addr, err := strconv.ParseUint(ss[0][2:len(ss[0])-1], 16, 64)
		if err != nil {
			panic(err)
		}
		fmt.Printf("addr 0x%x\n",addr)
		end, err := strconv.ParseUint(ss[1][2:len(ss[1])-1], 16, 64)
		if err != nil {
			panic(err)
		}
		fmt.Printf("end 0x%x\n",end)
		flags, err := strconv.ParseUint(ss[2][2:len(ss[2])-1], 16, 64)
		if err != nil {
			panic(err)
		}
		fmt.Printf("0x%x\n", flags)
		if flags&IORESOURCE_MEM == 0 {
			continue
		}
		fmt.Println("opening", fmt.Sprintf("%s%d", file, i))
		f, err := os.OpenFile(fmt.Sprintf("%s%d", file, i), os.O_RDWR, 0755)
		if err != nil {
			panic(err)
		}
		address, err := syscall.Mmap(int(f.Fd()), 0, int(end-addr+1), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
		if err != nil {
			panic(err)
		}
		fmt.Printf("CTRL 0x%x\n", uint32(address[0]) + uint32(address[1]) << 8 + uint32(address[2]) << 16 + uint32(address[3]) << 24)
		mrs = append(mrs, MemResource{
			PhysicalAddress: addr,
			Length:          end - addr + 1,
			Address:         address,
		})
		fmt.Println(ss)
		fmt.Println(line)
	}
	fmt.Println(mrs)
	return mrs
}

func main() {
	var command uint16
	var uio string
	flag.StringVar(&uio, "uio", "uio0", "uio device name")

	fUIO, err := os.OpenFile(fmt.Sprintf("/dev/%s", uio), os.O_RDWR, 0755)

	defer fUIO.Close()
	fConfig, err := os.OpenFile(fmt.Sprintf("/sys/class/uio/%s/device/config", uio), os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer fConfig.Close()
	var bufCommand = make([]byte, 2)
	_, err = fConfig.ReadAt(bufCommand, 4)
	if err != nil {
		log.Fatal(err)
	}
	command = binary.LittleEndian.Uint16(bufCommand) & ^PCI_COMMAND_INTX_DISABLE
	fmt.Println(commandString(command))
	binary.LittleEndian.PutUint16(bufCommand, command)
	mrs := parseResource(fmt.Sprintf("/sys/class/uio/%s/device/resource", uio))

	fmt.Println(&(mrs[0].Address[0]))
	var intBuf = make([]byte, 4)
	for i := 0; i < 10; i++ {
		fmt.Println("read count", i)
		_, err = fConfig.WriteAt(bufCommand, 4)
		if err != nil {
			log.Fatal(err)
		}
		_, err = fUIO.Read(intBuf)
		if err != nil {
			log.Fatal(err)
			break
		}
	}
}
