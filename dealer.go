package main

type Dealer struct {
	FaceUpCards   []Card
	FaceDownCards []Card
	Points        int
}

func (dealer *Dealer) GetPoints() int {
	var aces int
	cards := dealer.FaceUpCards
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
