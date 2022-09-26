package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mdlayher/wifi"
	"time"
)

const updateInterval = 100 * time.Millisecond

type updateTableMessage struct {
	Name     string
	Stations []*wifi.StationInfo
}

type periodicalUpdateTableMessage updateTableMessage

type Model struct {
	svc   *Service
	name  string
	table table.Model
}

func NewUI(svc *Service, name string) Model {
	columns := []table.Column{
		{Title: "MAC", Width: 18},
		{Title: "Signal (dB)", Width: 12},
		{Title: "RX bitrate", Width: 11},
		{Title: "TX bitrate", Width: 11},
		{Title: "RX packets", Width: 10},
		{Title: "TX packets", Width: 10},
		{Title: "Active (s)", Width: 11},
		{Title: "Inactive (s)", Width: 13},
		{Title: "TX fail", Width: 8},
		{Title: "TX retry", Width: 9},
		{Title: "Beacon loss", Width: 11},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithHeight(10),
		table.WithFocused(false),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := Model{
		svc:   svc,
		name:  name,
		table: t,
	}
	return m
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.updateStations(),
		m.periodicalUpdateStations(),
	)
}

func (m Model) View() string {
	return m.getHeader() + "\n" + baseStyle.Render(m.table.View()) + "\n"
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case updateTableMessage:
		if msg.Name == m.name {
			m.table.SetRows(m.getRows(msg.Stations))
		}
	case periodicalUpdateTableMessage:
		if msg.Name == m.name {
			m.table.SetRows(m.getRows(msg.Stations))
			cmds = append(cmds, m.periodicalUpdateStations())
		}
	}
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) updateStations() tea.Cmd {
	return func() tea.Msg {
		stations, _ := m.svc.GetStations()
		return updateTableMessage{
			Name:     m.name,
			Stations: stations,
		}
	}
}

func (m Model) periodicalUpdateStations() tea.Cmd {
	return tea.Tick(updateInterval, func(t time.Time) tea.Msg {
		stations, _ := m.svc.GetStations()
		return periodicalUpdateTableMessage{
			Name:     m.name,
			Stations: stations,
		}
	})
}

func (m Model) getRow(station *wifi.StationInfo) table.Row {
	return []string{
		fmt.Sprintf("%s", station.HardwareAddr),
		fmt.Sprintf("%d", station.Signal),
		fmt.Sprintf("%d", station.ReceiveBitrate/1000),
		fmt.Sprintf("%d", station.TransmitBitrate/1000),
		fmt.Sprintf("%d", station.ReceivedPackets),
		fmt.Sprintf("%d", station.TransmittedPackets),
		fmt.Sprintf("%f", station.Connected.Seconds()),
		fmt.Sprintf("%f", station.Inactive.Seconds()),
		fmt.Sprintf("%d", station.TransmitFailed),
		fmt.Sprintf("%d", station.TransmitRetries),
		fmt.Sprintf("%d", station.BeaconLoss),
	}
}

func (m Model) getRows(stations []*wifi.StationInfo) []table.Row {
	rows := make([]table.Row, 0, len(stations))
	for _, link := range stations {
		rows = append(rows, m.getRow(link))
	}
	return rows
}

func (m Model) getHeader() string {
	return fmt.Sprintf(" Name: %s\tMAC:%s\tFreq: %d\tType: %s",
		m.svc.iface.Name, m.svc.iface.HardwareAddr, m.svc.iface.Frequency, m.svc.iface.Type)
}
