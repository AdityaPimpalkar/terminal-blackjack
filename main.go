package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strconv"
	"time"
)

type gameModelKeys struct {
	Space key.Binding
	Tab   key.Binding
	Quit  key.Binding
}

var gameKeys = gameModelKeys{
	Space: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("<space>", "Hit!"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("<tab>", "Stand"),
	),
}

func (k gameModelKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Space, k.Tab}, // second column
	}
}

func (k gameModelKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Space, k.Tab}
}

type Game struct {
	Deck          Deck
	Player        Player
	Dealer        Dealer
	playerStood   bool
	isActiveRound bool
	keys          gameModelKeys
	help          help.Model
}

type GameStatus string

type GameStatusMsg struct {
	status GameStatus
}

func (g *Game) DealFirstHand() error {
	if len(g.Deck.Cards) == 0 {
		return fmt.Errorf("Deck of cards not found.")
	}
	g.Player.Cards = append(g.Player.Cards, g.Deck.Cards[0], g.Deck.Cards[2])
	g.Dealer.FaceUpCards = append(g.Dealer.FaceUpCards, g.Deck.Cards[1])
	g.Dealer.FaceDownCards = append(g.Dealer.FaceDownCards, g.Deck.Cards[3])
	g.Deck.Cards = g.Deck.Cards[4:len(g.Deck.Cards)]

	return nil
}

func (g *Game) PlayerHand() {
	g.Player.Cards = append(g.Player.Cards, g.Deck.Cards[0])
	g.Deck.Cards = g.Deck.Cards[1:len(g.Deck.Cards)]
}

const (
	playerWon     GameStatus = "playerWon"
	playerBusted  GameStatus = "playerBusted"
	dealerWon     GameStatus = "dealerWon"
	dealerDrawing GameStatus = "dealerDrawing"
	dealerBusted  GameStatus = "dealerBusted"
	tie           GameStatus = "tie"
)

func (g *Game) Status() tea.Cmd {
	var status GameStatus
	playerPoints := g.Player.GetPoints()
	dealerPoints := g.Dealer.GetPoints()
	if playerPoints == 21 {
		status = playerWon
	} else if playerPoints > 21 {
		status = playerBusted
	} else if dealerPoints > 21 {
		status = playerWon
	} else if dealerPoints < 17 {
		status = dealerDrawing
	} else if playerPoints > dealerPoints {
		status = playerWon
	} else {
		status = dealerWon
	}
	return func() tea.Msg {
		return GameStatusMsg{
			status: status,
		}
	}
}

func (g *Game) DealerHand() {
	dealer := g.Dealer
	deck := g.Deck
	if len(dealer.FaceDownCards) > 0 {
		dealer.FaceUpCards = append(dealer.FaceUpCards, dealer.FaceDownCards...)
		dealer.FaceDownCards = []Card{}
	} else {
		dealer.FaceUpCards = append(dealer.FaceUpCards, deck.Cards[0])
		deck.Cards = deck.Cards[1:len(deck.Cards)]
	}
	g.Dealer = dealer
	g.Deck = deck
}

func (g *Game) Reset() {
	g.Player.Cards = []Card{}
	g.Player.Points = 0
	g.Dealer.FaceUpCards = []Card{}
	g.Dealer.FaceDownCards = []Card{}
	g.isActiveRound = false
	g.playerStood = false
	g.Deck.Init()
}

const margin = 8
const height = 13

var (
	dealerCardStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("9"))
	playerCardStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("12"))
	cardHiddenStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3C3C3C"))
	cardTopFaceStyle = lipgloss.NewStyle().
				Bold(true).
				PaddingLeft(1)
	cardBottomFaceStyle = lipgloss.NewStyle().
				Bold(true).
				PaddingRight(1).
				Align(lipgloss.Right)
	cardSuitStyle = lipgloss.NewStyle().
			Bold(true).
			Height(11).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center)
)

func gameView(m MainModel) string {
	dealer := m.game.Dealer
	player := m.game.Player

	dealerPoints := dealer.GetPoints()
	playerPoints := player.GetPoints()

	var dealerCards []string

	var playerCards []string
	for _, card := range dealer.FaceUpCards {
		face, suit := card.GetCard()
		card := dealerCardStyle.
			Height(height).
			Width(m.windowWidth / margin).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Top,
					cardTopFaceStyle.Render(face),
					cardSuitStyle.
						Width(m.windowWidth/margin).
						Render(suit),
					cardBottomFaceStyle.
						Width(m.windowWidth/margin).
						Align(lipgloss.Right).
						Render(face),
				),
			)
		dealerCards = append(dealerCards, card)
	}

	for _, card := range player.Cards {
		face, suit := card.GetCard()
		card := playerCardStyle.
			Height(height).
			Width(m.windowWidth / margin).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Top,
					cardTopFaceStyle.Render(face),
					cardSuitStyle.
						Width(m.windowWidth/margin).
						Render(suit),
					cardBottomFaceStyle.
						Width(m.windowWidth/margin).
						Align(lipgloss.Right).
						Render(face),
				),
			)
		playerCards = append(playerCards, card)
	}

	if len(dealer.FaceDownCards) > 0 {
		hiddenCard := cardHiddenStyle.
			Height(height).
			Width(m.windowWidth / margin).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Top,
					cardTopFaceStyle.Render(" "),
					cardSuitStyle.
						Width(m.windowWidth/margin).
						Render("???"),
					cardBottomFaceStyle.
						Width(m.windowWidth/margin).
						Align(lipgloss.Right).
						Render(" "),
				),
			)
		dealerCards = append(dealerCards, hiddenCard)

	}
	dealerView := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().
			Width(m.windowWidth).
			Align(lipgloss.Center).
			Bold(true).
			Foreground(lipgloss.Color("9")).
			Render(fmt.Sprintf("Dealer score: %d", dealerPoints)),
		lipgloss.PlaceHorizontal(
			m.windowWidth,
			lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				dealerCards...,
			),
		),
	)

	playerView := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.PlaceHorizontal(
			m.windowWidth,
			lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				playerCards...,
			),
		),
		lipgloss.NewStyle().
			Width(m.windowWidth).
			Align(lipgloss.Center).
			Bold(true).
			Foreground(lipgloss.Color("12")).
			Render(fmt.Sprintf("Your score: %d", playerPoints)),
	)

	top := lipgloss.Place(m.windowWidth, m.windowHeight/2, lipgloss.Center, lipgloss.Top, dealerView)
	bottom := lipgloss.Place(m.windowWidth, m.windowHeight/2, lipgloss.Center, lipgloss.Bottom, playerView)

	finalView := lipgloss.JoinVertical(lipgloss.Center, top, bottom)

	return finalView
}

type bettingModelKeys struct {
	Reset key.Binding
	Enter key.Binding
	Quit  key.Binding
}

var bettingKeys = bettingModelKeys{
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

func (k bettingModelKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Enter, k.Reset}, // second column
	}
}

func (k bettingModelKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Enter, k.Reset}
}

var (
	textInputStyle = lipgloss.NewStyle().
		Width(25)
)

func validateInput(s string) error {
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

type Betting struct {
	InputBet     textinput.Model
	Balance      int
	Bet          int
	keys         bettingModelKeys
	help         help.Model
	WindowHeight int
	WindowWidth  int
}

func bettingView(m MainModel) string {
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
	return fmt.Sprintf(
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
	bettingKeys.Reset.SetEnabled(false)
	return MainModel{
		betting: Betting{
			InputBet: input,
			Balance:  1000,
			keys:     bettingKeys,
			help:     help.New(),
		},
	}
}

func (m MainModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
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
			m.keys.Enter.SetEnabled(true)
			m.keys.Reset.SetEnabled(false)
			return m, nil
		}
	case GameStatusMsg:
		switch msg.status {
		case tie:
			time.Sleep(2 * time.Second)
			m.game.Player.Balance = m.game.Player.Balance + m.betting.Bet
			m.game.Reset()
			return m, nil
		case playerWon:
			time.Sleep(2 * time.Second)
			m.game.Player.Balance = m.game.Player.Balance + (m.betting.Bet * 2)
			m.game.Reset()
			return m, nil
		case playerBusted:
			time.Sleep(2 * time.Second)
			m.game.Player.Balance = m.game.Player.Balance - m.betting.Bet
			m.betting.Bet = 0
			m.betting.InputBet.SetValue(" ")
			m.game.Reset()
			return m, nil
		case dealerDrawing:
			if m.game.playerStood {
				time.Sleep(1 * time.Second)
				m.game.DealerHand()
				return m, m.game.Status()
			}
		case dealerWon:
			time.Sleep(2 * time.Second)
			m.game.Reset()
			return m, nil
		}
	}

	if !m.game.isActiveRound {
		m.betting.InputBet.Validate = validateInput
		m.betting.InputBet, cmd = m.betting.InputBet.Update(msg)
		bet, _ := strconv.ParseInt(m.betting.InputBet.Value(), 10, 64)
		m.betting.Bet = int(bet)
	}

	return m, cmd
}

func (m MainModel) View() string {
	if m.game.isActiveRound {
		return gameView(m)
	}
	return bettingView(m)
}

func main() {
	p := tea.NewProgram(InitModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
