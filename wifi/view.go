package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type CombinedTables []tea.Model

func (ct CombinedTables) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, t := range ct {
		cmds = append(cmds, t.Init())
	}
	return tea.Batch(cmds...)
}

func (ct CombinedTables) View() string {
	var rendered string
	for _, t := range ct {
		rendered += t.View() + "\n\n"
	}
	return rendered
}

func (ct CombinedTables) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	for i, t := range ct {
		m, cmd := t.Update(msg)
		ct[i] = m
		cmds = append(cmds, cmd)
	}
	return ct, tea.Batch(cmds...)
}
