//go:build linux
// +build linux

package main

import (
	"github.com/jsimonetti/rtnetlink"
	"golang.org/x/sys/unix"
	"net"
)

func main() {
	iface, _ := net.InterfaceByName("lan0")
	conn, _ := rtnetlink.Dial(nil)
	cidr := net.IPNet{
		IP:   []byte{192, 168, 0, 100},
		Mask: []byte{255, 255, 255, 0},
	}
	pLen, _ := cidr.Mask.Size()
	_ = conn.Address.New(&rtnetlink.AddressMessage{
		Family:       uint8(unix.AF_INET),
		PrefixLength: uint8(pLen),
		Scope:        unix.RT_SCOPE_UNIVERSE,
		Index:        uint32(iface.Index),
		Attributes: &rtnetlink.AddressAttributes{
			Address: cidr.IP,
		},
	})
	_ = conn.Link.Set(&rtnetlink.LinkMessage{
		Family: unix.AF_UNSPEC,
		Index:  uint32(iface.Index),
		Flags:  unix.IFF_UP,
		Change: unix.IFF_UP,
	})
	_ = conn.Close()
}
