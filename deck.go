package main

import "math/rand"

type Card struct {
	Suit  string
	Face  string
	Point int
}

func (card *Card) GetCard() (string, string) {
	var face string
	var suit string
	switch card.Face {
	case "Ace":
		face = "A"
	case "King":
		face = "K"
	case "Jack":
		face = "J"
	case "Queen":
		face = "Q"
	default:
		face = card.Face
	}

	switch card.Suit {
	case "Spades":
		suit = "♠"
	case "Hearts":
		suit = "♥"
	case "Diamonds":
		suit = "♦"
	case "Clubs":
		suit = "♣"
	}
	return face, suit
}

type Deck struct {
	Cards []Card
}

func (deck *Deck) Init() {
	suit := [4]string{"Spades", "Hearts", "Diamonds", "Clubs"}
	cards := [13]string{"Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King"}
	cardValues := map[string]int{
		"Ace":   0,
		"2":     2,
		"3":     3,
		"4":     4,
		"5":     5,
		"6":     6,
		"7":     7,
		"8":     8,
		"9":     9,
		"10":    10,
		"Jack":  10,
		"Queen": 10,
		"King":  10,
	}

	for i := 0; i <= len(suit)-1; i++ {
		for j := 0; j <= len(cards)-1; j++ {
			deck.Cards = append(deck.Cards, Card{
				suit[i],
				cards[j],
				cardValues[cards[j]],
			})
		}
	}
}

func (deck *Deck) Shuffle() {
	perm := rand.Perm(len(deck.Cards))
	for i := 0; i < len(perm)-1; i++ {
		deck.Cards[i] = deck.Cards[perm[i]]
	}
}
