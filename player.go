package main

type Player struct {
	Balance int
	Bet     int
	Cards   []Card
	Points  int
}

func (player *Player) GetPoints() int {
	cards := player.Cards
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
