package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mdlayher/wifi"
	"log"
)

func main() {
	client, err := wifi.New()
	if err != nil {
		log.Fatal(err)
	}

	ifaces, err := client.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	var ct CombinedTables
	for _, iface := range ifaces {
		svc, err := NewService(client, iface)
		if err != nil {
			log.Fatal(err)
		}
		ui := NewUI(svc, iface.Name)
		ct = append(ct, ui)
	}

	if err = tea.NewProgram(ct, tea.WithAltScreen()).Start(); err != nil {
		log.Fatal(err)
	}
}
