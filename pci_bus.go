package dufu

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	PATH_SYSFS_PCI_DEVICES = "/sys/bus/pci/devices"

	PCI_MAX_RESOURCE = 6 // At most there 6 BARs from 0 to 5, but most devices just use 2 or 3 BARs.
	// IO resource type
	IORESOURCE_IO  = 0x00000100
	IORESOURCE_MEM = 0x00000200
)

type PCIBus struct {
	name    string
	devices map[string]PCIDevice
	// drivers map[string]Driver
}

func (p PCIBus) Name() string {
	return p.name
}

func (p *PCIBus) Scan() {
	files, err := ioutil.ReadDir(PATH_SYSFS_PCI_DEVICES)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
		pciAddress := parsePCIAddressFormat(file.Name())
		fmt.Println(pciAddress)
		p.scanOne(pciAddress)
	}
}

type PCIAddress struct {
	Domain   uint16
	Bus      uint8
	Device   uint8
	Function uint8
}

func parsePCIAddressFormat(s string) PCIAddress {
	ss := strings.Split(s, ":")
	domain, _ := strconv.ParseUint(ss[0], 10, 16)
	bus, _ := strconv.ParseUint(ss[1], 10, 8)
	ss = strings.Split(ss[2], ".")
	device, _ := strconv.ParseUint(ss[0], 10, 8)
	function, _ := strconv.ParseUint(ss[1], 10, 8)
	return PCIAddress{
		Domain:   uint16(domain),
		Bus:      uint8(bus),
		Device:   uint8(device),
		Function: uint8(function),
	}
}

func parseSysfsValue(filename string) (uint64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return 0, err
	}
	if bs[len(bs)-1] == '\n' && string(bs[:2]) == "0x" {
		bs = bs[2 : len(bs)-1]
	}
	i, err := strconv.ParseUint(string(bs), 16, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (p *PCIBus) scanOne(addr PCIAddress) {
	var device PCIDevice
	var tmp uint64
	var err error

	device.name = fmt.Sprintf("%.4d:%.2d:%.2d.%d", addr.Domain, addr.Bus, addr.Device, addr.Function)
	dirname := filepath.Join(PATH_SYSFS_PCI_DEVICES, device.name)
	tmp, err = parseSysfsValue(filepath.Join(dirname, "vendor"))
	if err != nil {
		panic(err)
	}
	device.Vendor = uint16(tmp)

	tmp, err = parseSysfsValue(filepath.Join(dirname, "device"))
	if err != nil {
		panic(err)
	}
	device.Device = uint16(tmp)

	tmp, err = parseSysfsValue(filepath.Join(dirname, "subsystem_vendor"))
	if err != nil {
		panic(err)
	}
	device.SubsystemVendor = uint16(tmp)

	tmp, err = parseSysfsValue(filepath.Join(dirname, "subsystem_device"))
	if err != nil {
		panic(err)
	}
	device.SubsystemDevice = uint16(tmp)

	tmp, err = parseSysfsValue(filepath.Join(dirname, "class"))
	if err != nil {
		panic(err)
	}
	device.Class = uint32(tmp)

	tmp, err = parseSysfsValue(filepath.Join(dirname, "max_vfs"))
	if err != nil {
		tmp, _ = parseSysfsValue(filepath.Join(dirname, "sriov_numvfs"))
	}
	device.MaxVFs = uint16(tmp)

	tmp, _ = parseSysfsValue(filepath.Join(dirname, "numa_node"))
	device.NumaNode = tmp

	device.Mem = parsePCIResource(filepath.Join(dirname, "resource"))

	device.Driver = parsePCIDriver(filepath.Join(dirname, "driver"))
	fmt.Println(device)
	// add device
	p.devices[device.Name()] = device
}

func parsePCIDriver(filename string) string {
	path, err := os.Readlink(filename)
	if err != nil {
		fmt.Println("no driver found")
		return ""
	}
	_, driver := filepath.Split(path)
	return driver
}

func parsePCIResource(file string) []MemResource {
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
		end, err := strconv.ParseUint(ss[1][2:len(ss[1])-1], 16, 64)
		if err != nil {
			panic(err)
		}
		flags, err := strconv.ParseUint(ss[2][2:len(ss[2])-1], 16, 64)
		if err != nil {
			panic(err)
		}
		if flags&IORESOURCE_MEM == 0 {
			continue
		}
		mrs = append(mrs, MemResource{
			PhysicalAddress: addr,
			Length:          end - addr + 1,
		})
	}
	return mrs
}

func (p *PCIBus) Probe() {
	for _, device := range p.devices {
		p.probeAllPCIDrivers(device)
	}
}

// TODO: PCIDevice can be named after Device
func (p *PCIBus) probeAllPCIDrivers(device PCIDevice) {
	// only on driver
	// driver := NewEmDriver()
	// TODO: id table match
	// currently, there is one e1000 driver
	// driver.Probe()
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
		fmt.Printf("addr 0x%x\n", addr)
		end, err := strconv.ParseUint(ss[1][2:len(ss[1])-1], 16, 64)
		if err != nil {
			panic(err)
		}
		fmt.Printf("end 0x%x\n", end)
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
		fmt.Printf("CTRL 0x%x\n", uint32(address[0])+uint32(address[1])<<8+uint32(address[2])<<16+uint32(address[3])<<24)
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

func NewPCIBus() Bus {
	return &PCIBus{
		name:    "pci",
		devices: make(map[string]PCIDevice),
		// drivers:make(map[string]Driver),
	}
}
