package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/edsrzf/mmap-go"
)

func main() {
	fUIO, err := os.OpenFile("/dev/uio1", os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer fUIO.Close()

	fa, err := os.Open("/sys/class/uio/uio1/maps/map0/addr")
	if err != nil {
		log.Fatal(err)
	}
	defer fa.Close()

	fs, err := os.Open("/sys/class/uio/uio1/maps/map0/size")
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()
	addrBuf, err := ioutil.ReadAll(fa)
	if err != nil {
		log.Fatal(err)
	}
	sizeBuf, err := ioutil.ReadAll(fs)
	if err != nil {
		log.Fatal(err)
	}
	addr, err := strconv.ParseUint(string(addrBuf[2:len(addrBuf)-1]), 16, 64)
	if err != nil {
		log.Fatal(err)
	}
	size, err := strconv.ParseUint(string(sizeBuf[2:len(sizeBuf)-1]), 16, 64)
	if err != nil {
		log.Fatal(err)
	}
	m, err := mmap.MapRegion(fUIO, int(size), mmap.RDWR, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("virtual address: %p\n", &m[0])
	fmt.Printf("physical address: %x\n", addr)
	fmt.Printf("mapping size: %d\n", size)
	var intBuf = make([]byte, 4)
	for i := 0; i < 10; i++ {
		fmt.Println("read count", i)
		_, err = fUIO.Read(intBuf)
		if err != nil {
			log.Fatal(err)
			break
		}
	}
}
