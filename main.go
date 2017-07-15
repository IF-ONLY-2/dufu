package main

import (
    "encoding/binary"
    "fmt"
    "log"
    "os"
)

func main() {
    var command uint16
    fUIO, err := os.OpenFile("/dev/uio0", os.O_RDWR, 0755)
    if err != nil {
        log.Fatal(err)
    }
    defer fUIO.Close()
    fConfig, err := os.OpenFile("/sys/class/uio/uio0/device/config", os.O_RDWR, 0755)
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
