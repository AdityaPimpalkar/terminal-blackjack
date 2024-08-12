package main

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
