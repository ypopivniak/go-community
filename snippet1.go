//go:build linux
// +build linux

package main

import (
	"github.com/vishvananda/netlink"
	"net"
)

func main() {
	iface, _ := netlink.LinkByName("eth0")
	_ = netlink.AddrAdd(iface, &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   []byte{192, 168, 0, 100},
			Mask: []byte{255, 255, 255, 0},
		},
	})
	_ = netlink.LinkSetUp(iface)
}
