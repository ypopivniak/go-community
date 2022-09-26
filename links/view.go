//go:build linux
// +build linux

package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jsimonetti/rtnetlink"
	"sort"
	"time"
)

const updateInterval = 100 * time.Millisecond

type updateTableMessage []Link

type periodicalUpdateTableMessage updateTableMessage

type Model struct {
	svc     *Service
	table   table.Model
	links   map[string]Link
	watchCh chan Link
}

func NewUI(svc *Service) Model {
	columns := []table.Column{
		{Title: "Index", Width: 6},
		{Title: "Name", Width: 10},
		{Title: "State", Width: 10},
		{Title: "MAC", Width: 18},
		{Title: "RX packets", Width: 10},
		{Title: "TX packets", Width: 10},
		{Title: "RX bytes", Width: 14},
		{Title: "TX bytes", Width: 14},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithHeight(10),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	watchCh, _ := svc.WatchLinks()
	m := Model{
		svc:     svc,
		table:   t,
		links:   make(map[string]Link),
		watchCh: watchCh,
	}
	return m
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.updateLinks(m.watchCh),
		m.periodicalUpdateLinks(m.svc),
	)
}

func (m Model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
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
		m.setLinks(msg)
	case periodicalUpdateTableMessage:
		m.setLinks(msg)
		cmds = append(cmds, m.periodicalUpdateLinks(m.svc))
	}
	m.table.SetRows(m.getRows())
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) setLinks(links []Link) {
	for _, link := range links {
		m.links[link.Attributes.Name] = link
	}
}

func (m Model) updateLinks(watchCh chan Link) tea.Cmd {
	return func() tea.Msg {
		return updateTableMessage{<-watchCh}
	}
}

func (m Model) periodicalUpdateLinks(svc *Service) tea.Cmd {
	return tea.Tick(updateInterval, func(t time.Time) tea.Msg {
		links, _ := svc.GetLinks()
		return periodicalUpdateTableMessage(links)
	})
}

func (m Model) getRow(link Link) table.Row {
	return []string{
		fmt.Sprintf("%d", link.Index),
		link.Attributes.Name,
		m.getState(link.Attributes.OperationalState),
		link.Attributes.Address.String(),
		fmt.Sprintf("%d", link.Attributes.Stats.RXPackets),
		fmt.Sprintf("%d", link.Attributes.Stats.TXPackets),
		fmt.Sprintf("%d", link.Attributes.Stats.RXBytes),
		fmt.Sprintf("%d", link.Attributes.Stats.TXBytes),
	}
}

func (m Model) getRows() []table.Row {
	links := make([]Link, 0, len(m.links))
	rows := make([]table.Row, 0, len(m.links))
	for _, link := range m.links {
		links = append(links, link)
	}
	// Had no time to sort in on loop, sorry :)
	sort.Slice(links, func(i, j int) bool {
		return links[i].Index > links[j].Index
	})
	for _, link := range links {
		rows = append(rows, m.getRow(link))
	}
	return rows
}

func (Model) getState(state rtnetlink.OperationalState) string {
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
