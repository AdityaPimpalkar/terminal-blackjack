package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type MainModel struct {
	betting      Betting
	game         Game
	windowWidth  int
	windowHeight int
}

func InitModel() MainModel {
	input := textinput.New()
	input.Placeholder = " Place your bet!"
	input.Focus()
	input.Prompt = "$"
	input.CharLimit = 12
	input.Validate = validateInput
	bettingModeKeys.Reset.SetEnabled(false)
	return MainModel{
		betting: Betting{
			InputBet: input,
			Balance:  1000,
			keys:     bettingModeKeys,
			help:     help.New(),
		},
		game: Game{
			help: help.New(),
			keys: gameModeKeys,
		},
	}
}

func (m MainModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if !m.game.isActiveRound && m.betting.Bet > 0 && m.betting.Bet <= m.betting.Balance {
				m.game.isActiveRound = true
				m.game.Deck.Init()
				m.game.Deck.Shuffle()
				m.game.Player.Bet = m.betting.Bet
				m.game.Player.Balance = m.betting.Balance - m.betting.Bet
				err := m.game.DealFirstHand()
				if err != nil {
					panic(err)
				}
				return m, m.game.Status()
			}
		case tea.KeySpace:
			if m.game.isActiveRound && !m.game.playerStood {
				m.game.PlayerHand()
				return m, m.game.Status()
			}
		case tea.KeyTab:
			if m.game.isActiveRound && !m.game.playerStood {
				m.game.playerStood = true
				m.game.DealerHand()
				return m, m.game.Status()
			}
		case tea.KeyCtrlK:
			m.betting.Balance = 1000
			m.betting.InputBet.Focus()
			m.betting.keys.Enter.SetEnabled(true)
			m.betting.keys.Reset.SetEnabled(false)
			return m, nil
		}
	case GameStatusMsg:
		switch msg.status {
		case tie:
			time.Sleep(2 * time.Second)
			m.betting.Balance = m.game.Player.Balance + m.betting.Bet
			m.betting.Bet = 0
			m.betting.InputBet.Reset()
			m.game.Reset()
			return m, nil
		case playerWon:
			time.Sleep(2 * time.Second)
			m.betting.Balance = m.game.Player.Balance + (m.betting.Bet * 2)
			m.betting.Bet = 0
			m.betting.InputBet.Reset()
			m.game.Reset()
			return m, nil
		case playerBusted:
			time.Sleep(2 * time.Second)
			m.betting.Balance = m.game.Player.Balance
			m.betting.Bet = 0
			m.betting.InputBet.Reset()
			m.game.Reset()
			return m, nil
		case dealerDrawing:
			if m.game.playerStood {
				time.Sleep(2 * time.Second)
				m.game.DealerHand()
				return m, m.game.Status()
			}
		case dealerWon:
			time.Sleep(2 * time.Second)
			m.betting.Balance = m.game.Player.Balance
			m.betting.Bet = 0

			m.betting.InputBet.Reset()
			m.game.Reset()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.betting.InputBet, cmd = m.betting.InputBet.Update(msg)
	m.betting.Bet, _ = strconv.Atoi(m.betting.InputBet.Value())
	return m, cmd
}

func (m MainModel) View() string {
	if m.game.isActiveRound {
		return GameView(m)
	}
	return BettingView(m)
}

func main() {
	p := tea.NewProgram(InitModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
