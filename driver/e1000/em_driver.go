package e1000

import (
	"fmt"
)

type EmDriver struct {
	HwAddress []byte // Is this consistent ?
}

func (e *EmDriver) Probe() {
	fmt.Println("em driver probe")
	/* 注册回调函数，开启中断 */
	return
}

func (e *EmDriver) Remove() {
	fmt.Println("em dirver remove")
	return
}

func NewEmDriver() *EmDriver {
	return &EmDriver{}
}
