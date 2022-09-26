package main

import (
	"github.com/mdlayher/genetlink"
	"log"
)

func main() {
	conn, err := genetlink.Dial(nil)
	if err != nil {
		log.Fatalf("Unable to connect to generic netlink socket: %s", err)
	}
	families, err := conn.ListFamilies()
	if err != nil {
		log.Fatalf("Unable to get fimilies list: %s", err)
	}
	for _, family := range families {
		log.Printf("Family name: %s, ID: %d, version: %d", family.Name, family.ID, family.Version)
		for _, group := range family.Groups {
			log.Printf("\tGroup name: %s, ID: %d", group.Name, group.ID)
		}
	}
	_ = conn.Close()
}
