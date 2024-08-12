package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type gameKeys struct {
	Space key.Binding
	Tab   key.Binding
	Quit  key.Binding
}

type Game struct {
	Deck          Deck
	Player        Player
	Dealer        Dealer
	keys          gameKeys
	help          help.Model
	playerStood   bool
	isActiveRound bool
}

type GameStatus string

type GameStatusMsg struct {
	status GameStatus
}

const (
	playerWon     GameStatus = "playerWon"
	playerBusted  GameStatus = "playerBusted"
	dealerWon     GameStatus = "dealerWon"
	dealerDrawing GameStatus = "dealerDrawing"
	dealerBusted  GameStatus = "dealerBusted"
	tie           GameStatus = "tie"
)

var gameModeKeys = gameKeys{
	Space: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("<space>", "Hit!"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("<tab>", "Stand"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "Quit"),
	),
}

func (k gameKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Space, k.Tab}, // second column
	}
}

func (k gameKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Space, k.Tab}
}

func (g *Game) DealFirstHand() error {
	if len(g.Deck.Cards) == 0 {
		return fmt.Errorf("Deck of cards not found")
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

func (g *Game) Status() tea.Cmd {
	var status GameStatus
	playerPoints := g.Player.GetPoints()
	dealerPoints := g.Dealer.GetPoints()
	if playerPoints == dealerPoints {
		status = tie
	} else if playerPoints == 21 {
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

func GameView(m MainModel) string {
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
						Render("?"),
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
		lipgloss.NewStyle().
			Width(m.windowWidth).
			Align(lipgloss.Center).
			Render(
				m.game.help.View(m.game.keys),
			),
	)

	top := lipgloss.Place(m.windowWidth, m.windowHeight/2, lipgloss.Center, lipgloss.Top, dealerView)
	bottom := lipgloss.Place(m.windowWidth, m.windowHeight/2, lipgloss.Center, lipgloss.Bottom, playerView)

	finalView := lipgloss.JoinVertical(lipgloss.Center, top, bottom)

	return finalView
}
