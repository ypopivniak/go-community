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

const rtmGroupLink = 0x1

type Link *rtnetlink.LinkMessage

type Service struct {
	watchConn *netlink.Conn
	mgmtConn  *rtnetlink.Conn
}

func NewService() (*Service, error) {
	var err error
	var svc = new(Service)

	svc.mgmtConn, err = rtnetlink.Dial(nil)
	if err != nil {
		log.Fatalf("Could not connect to rtnetlink socket: %s", err)
	}

	svc.watchConn, err = netlink.Dial(unix.NETLINK_ROUTE, nil)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func (svc *Service) GetLinks() ([]Link, error) {
	lms, err := svc.mgmtConn.Link.List()
	if err != nil {
		return nil, err
	}
	links := make([]Link, 0, len(lms))
	for _, lm := range lms {
		link := lm
		links = append(links, &link)
	}
	return links, nil
}

func (svc *Service) WatchLinks() (chan Link, chan error) {
	errCh := make(chan error)
	if err := svc.watchConn.JoinGroup(rtmGroupLink); err != nil {
		errCh <- err
		return nil, errCh
	}
	linkCh := make(chan Link)
	go func() {
		for {
			raw, err := svc.watchConn.Receive()
			if err != nil {
				errCh <- io.EOF
				break
			} else if err != nil {
				errCh <- err
				return
			}
			for _, msg := range raw {
				lm := new(rtnetlink.LinkMessage)
				if err = lm.UnmarshalBinary(msg.Data); err == nil {
					linkCh <- lm
				}
			}
		}
		_ = svc.watchConn.LeaveGroup(rtmGroupLink)
	}()
	return linkCh, errCh
}
