package dufu

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/edsrzf/mmap-go"
)

func testmain() {
	var uio string
	flag.StringVar(&uio, "uio", "uio0", "uio device name")

	fUIO, err := os.OpenFile(fmt.Sprintf("/dev/%s", uio), os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer fUIO.Close()

	fa, err := os.Open(fmt.Sprintf("/sys/class/uio/%s/maps/map0/addr", uio))
	if err != nil {
		log.Fatal(err)
	}
	defer fa.Close()

	fs, err := os.Open(fmt.Sprintf("/sys/class/uio/%s/maps/map0/size", uio))
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
