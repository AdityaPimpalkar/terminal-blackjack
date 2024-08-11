package main

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type bettingKeys struct {
	Reset key.Binding
	Enter key.Binding
	Quit  key.Binding
}

type Betting struct {
	InputBet     textinput.Model
	Balance      int
	Bet          int
	keys         bettingKeys
	help         help.Model
	WindowHeight int
	WindowWidth  int
}


var bettingModeKeys = bettingKeys{
	Reset: key.NewBinding(
		key.WithKeys("ctrl+k"),
		key.WithHelp("ctrl+k", "Reset"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<enter>", "Play!"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "Quit"),
	),
}

func (k bettingKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Enter, k.Reset}, // second column
	}
}

func (k bettingKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Enter, k.Reset}
}


var (
	textInputStyle = lipgloss.NewStyle().
		Width(25)
)

func validateInput(s string) error {
	if _, err := strconv.Atoi(s); err != nil {
		return fmt.Errorf("invalid input")
	}
	return nil
}

func BettingView(m MainModel) string {
	if m.betting.Bet == 0 {
		m.betting.keys.Enter.SetEnabled(false)
	}
	if m.betting.Balance == 0 {
		m.betting.keys.Enter.SetEnabled(false)
		m.betting.keys.Reset.SetEnabled(true)
		m.betting.InputBet.Blur()
	}
	remainingBalance := m.betting.Balance - m.betting.Bet
	var errorText string
	balanceTextColor := lipgloss.Color("36")

	if remainingBalance < 0 {
		remainingBalance = 0
		errorText = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("Bigger than balance!\n")
		balanceTextColor = lipgloss.Color("#3C3C3C")
		m.betting.keys.Enter.SetEnabled(false)
	}
	return fmt.Sprint(
		lipgloss.Place(
			m.windowWidth,
			m.windowHeight,
			lipgloss.Position(lipgloss.Center),
			lipgloss.Position(lipgloss.Center),
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5")).Render("Your bet? \n"),
				textInputStyle.Render(m.betting.InputBet.View()),
				errorText,
				lipgloss.NewStyle().Bold(true).Foreground(balanceTextColor).Render(
					fmt.Sprintf("Balance: $%d\n", remainingBalance),
				),
				m.betting.help.View(m.betting.keys),
			),
		),
	)
}
