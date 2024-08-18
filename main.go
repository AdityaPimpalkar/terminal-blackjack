package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProgressStatus string

type ProgressBarMsg struct {
	status ProgressStatus
}

type MainModel struct {
	betting           Betting
	game              Game
	progress          progress.Model
	celebrate         Celebrate
	isProgressView    bool
	isCelebrationView bool
	phase             int
	windowWidth       int
	windowHeight      int
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
		progress:  progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C")),
		phase:     1,
		celebrate: InitCelebration(),
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
		m.progress.Width = msg.Width - 2*2 - 4
		if m.progress.Width > 70 {
			m.progress.Width = 70
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.isProgressView = true
			return m, tea.Batch(
				tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
					return ProgressBarMsg{
						status: "loading",
					}
				}),
				tea.Tick(time.Millisecond*1500, func(t time.Time) tea.Msg {
					return ProgressBarMsg{
						status: "done",
					}
				}),
				tea.Tick(time.Millisecond*2500, func(t time.Time) tea.Msg {
					return ProgressBarMsg{
						status: "started",
					}
				}),
			)
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
	case ProgressBarMsg:
		switch msg.status {
		case "loading":
			m.phase++
			return m, nil
		case "done":
			m.phase++
			return m, nil
		case "started":
			m.phase++
		}
		if m.phase > 3 {
			m.isProgressView = false
			m.phase = 1
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
			m.isCelebrationView = true
			m.celebrate.confetti.Particles = append(
				m.celebrate.confetti.Particles,
				Spawn(m.windowWidth, m.windowHeight)...,
			)
			m.celebrate.confetti.Frame.Width = m.windowWidth
			m.celebrate.confetti.Frame.Height = m.windowHeight

			return m, tea.Batch(
				AnimateCelebration(),
				Winner(),
			)
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
	case CelebrationMsg:
		m.celebrate.confetti.Update()
		return m, AnimateCelebration()

	case WinnerMsg:
		m.isCelebrationView = false
		return m, nil
	}

	var cmd tea.Cmd
	m.betting.InputBet, cmd = m.betting.InputBet.Update(msg)
	m.betting.Bet, _ = strconv.Atoi(m.betting.InputBet.Value())
	return m, cmd
}

func ProgressView(m MainModel) string {
	var msg string
	switch m.phase {
	case 1:
		msg = "Shuffling cards..."
	case 2:
		msg = "Starting round..."
	case 3:
		msg = "Good luck!"
	}

	return lipgloss.Place(
		m.windowWidth,
		m.windowHeight,
		lipgloss.Position(lipgloss.Center),
		lipgloss.Position(lipgloss.Center),
		lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().Bold(true).Render(msg+"\n"),
			m.progress.ViewAs(float64(m.phase)/3.0),
		),
	)
}

func (m MainModel) View() string {
	if m.game.isActiveRound {
		return GameView(m)
	} else if m.isProgressView {
		return ProgressView(m)
	} else if m.isCelebrationView {
		return CelebrationView(m)
	}
	return BettingView(m)
}

func main() {
	p := tea.NewProgram(InitModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
