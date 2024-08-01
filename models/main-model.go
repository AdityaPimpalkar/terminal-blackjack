package models

import (
	// "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type MainModel struct {
	Width  int
	Height int
}

func (m MainModel) Init() tea.Cmd {
	return nil
}
