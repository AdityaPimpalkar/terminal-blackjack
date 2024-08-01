package main

type Player struct {
	Balance int
	Bet     int
	Cards   []Card
	Points  int
}

func (player *Player) GetPoints() int {
	var aces int
	cards := player.Cards
	var totalPoints int
	for i := 0; i < len(cards); i++ {
		cardFace := cards[i].Face
		if cardFace == "Ace" {
			aces += 1
		}
		totalPoints += cards[i].Point
	}
	return totalPoints
}


func (player *Player) GetCardNames() []string {
	var cardNames []string
	cards := player.Cards
	for i := 0; i < len(cards); i++ {
		cardFace := cards[i].Face
		cardSuit := cards[i].Suit
		cardNames = append(cardNames, cardFace+" of "+cardSuit)
	}
	return cardNames
}