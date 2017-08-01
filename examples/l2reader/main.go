package main

import (
	"github.com/ggaaooppeenngg/dufu"
)

func main() {
	tap, err := dufu.NewTAP("tap_dufu")
	if err != nil {
		panic(err)
	}
	l2 := &dufu.L2Layer{TapDevice: tap}
	l2.Loop()
}
