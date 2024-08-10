package main

type Dealer struct {
	FaceUpCards   []Card
	FaceDownCards []Card
	Points        int
}

func (dealer *Dealer) GetPoints() int {
	cards := dealer.FaceUpCards
	var totalPoints int
	for i := 0; i < len(cards); i++ {
		cardFace := cards[i].Face
		if cardFace == "Ace" {
			if totalPoints+10 <= 21 {
				totalPoints += 10
			} else {
				totalPoints += 1
			}
		} else {
			totalPoints += cards[i].Point
		}
	}
	return totalPoints
}
