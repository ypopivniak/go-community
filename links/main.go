//go:build linux
// +build linux

package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func main() {
	svc, err := NewService()
	if err != nil {
		log.Fatalf("could not create links service: %s", err)
	}
	ui := NewUI(svc)
	if err = tea.NewProgram(ui, tea.WithAltScreen()).Start(); err != nil {
		log.Fatal(err)
	}
}
