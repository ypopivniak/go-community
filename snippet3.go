//go:build linux
// +build linux

package main

import (
	"github.com/jsimonetti/rtnetlink"
	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
	"io"
	"log"
)

const (
	rtmGroupLink = 0x1
)

func main() {
	watchConn, err := netlink.Dial(unix.NETLINK_ROUTE, nil)
	if err != nil {
		log.Fatalf("Unable to connect to netlink socket: %s", err)
	}
	if err = watchConn.JoinGroup(rtmGroupLink); err != nil {
		log.Fatalf("Unable to join multicast group: %s", err)
	}
	for {
		raw, err := watchConn.Receive()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		for _, msg := range raw {
			lm := rtnetlink.LinkMessage{}
			err = lm.UnmarshalBinary(msg.Data)
			if err != nil {
				continue
			}
			log.Printf("Interface [%s] changed state to [%s]",
				lm.Attributes.Name, getState(lm.Attributes.OperationalState))
		}
	}
	_ = watchConn.Close()
}

func getState(state rtnetlink.OperationalState) string {
	switch state {
	case 2:
		return "DOWN"
	case 3:
		return "NO CABLE"
	case 4:
		return "TESTING"
	case 5:
		return "DORMANT"
	case 6:
		return "UP"
	default:
		return "UNKNOWN"
	}
}
