package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func wait() {
	time.Sleep(3000 * time.Millisecond)
}

func main() {

	deck := Deck{}
	deck.Init()

	keepPlaying := true

	for {
		isNewRound := true
		player := Player{
			Balance: 1000,
		}
		dealer := Dealer{}
		for {
			if !keepPlaying {
				break
			}
			fmt.Printf("Account balance: $%d\n", player.Balance)

			fmt.Println("Place your bet: ")
			var bet int
			fmt.Scan(&bet)

			player.Bet = bet
			fmt.Printf("You bet $%d!\n", bet)

			fmt.Println("GOOD LUCK!!")

			wait()
			clear()

			if isNewRound {
				fmt.Println("Shuffling cards...")
				deck.Shuffle()
				isNewRound = false
			}

			time.Sleep(1500 * time.Millisecond)

			player.Cards = append(player.Cards, deck.Cards[0], deck.Cards[2])
			dealer.FaceUpCards = append(dealer.FaceUpCards, deck.Cards[1])
			dealer.FaceDownCards = append(dealer.FaceDownCards, deck.Cards[3])
			deck.Cards = deck.Cards[4:len(deck.Cards)]

			for {
				clear()

				dealerCards := dealer.GetCardNames()
				playerCards := player.GetCardNames()

				dealerPoints := dealer.GetPoints()
				playerPoints := player.GetPoints()

				fmt.Println("Dealer(" + strconv.Itoa(dealerPoints) + ")")
				for i := 0; i < len(dealerCards); i++ {
					fmt.Println(dealerCards[i])
				}
				if len(dealer.FaceDownCards) > 0 {
					fmt.Println("(1 hidden card)")
				}
				fmt.Println()
				fmt.Println()

				fmt.Println("Player(" + strconv.Itoa(playerPoints) + ")")
				for i := 0; i < len(playerCards); i++ {
					fmt.Println(playerCards[i])
				}

				fmt.Println()
				fmt.Println()

				playerWon := playerPoints == 21
				playerBusted := playerPoints > 21
				dealerBusted := dealerPoints > 21
				tie := len(dealer.FaceDownCards) == 0 && playerPoints == dealerPoints
				dealerIsLow := len(dealer.FaceDownCards) == 0 && dealerPoints <= 16
				dealerIsPar := len(dealer.FaceDownCards) == 0 && dealerPoints >= 17

				if playerWon {
					fmt.Printf("You WON $%d!!\n", bet)
					player.Balance = player.Balance + (player.Bet * 2)
					wait()
					clear()
					break
				} else if playerBusted {
					fmt.Println("You lost...")
					player.Balance = player.Balance - player.Bet
					wait()
					clear()
					break
				} else if dealerBusted {
					fmt.Printf("You WON $%d!!\n", bet)
					player.Balance = player.Balance + (player.Bet * 2)
					wait()
					clear()
					break
				} else if dealerIsLow {
					fmt.Println("Dealers hand...")
					wait()
					dealer.FaceUpCards = append(dealer.FaceUpCards, deck.Cards[0])
					deck.Cards = deck.Cards[1:len(deck.Cards)]
					continue
				} else if tie {
					fmt.Println("It was a tie.")
					wait()
					clear()
					break
				} else if dealerIsPar {
					if playerPoints > dealerPoints {
						fmt.Printf("You WON $%d!!\n", bet)
						player.Balance = player.Balance + (player.Bet * 2)
					} else {
						fmt.Println("You lost...")
						player.Balance = player.Balance - player.Bet
					}
					wait()
					clear()
					break
				}

				fmt.Println("Hit? [Y/n]")
				var action string
				fmt.Scan(&action)

				switch action {
				case "Y":
					{
						fmt.Println("Players hand...")
						wait()
						player.Cards = append(player.Cards, deck.Cards[0])
						deck.Cards = deck.Cards[1:len(deck.Cards)]
						continue
					}
				case "n":
					{
						dealer.FaceUpCards = append(dealer.FaceUpCards, dealer.FaceDownCards...)
						dealer.FaceDownCards = []Card{}
						break
					}
				}
			}

			player.Cards = []Card{}
			dealer.FaceUpCards = []Card{}
			dealer.FaceDownCards = []Card{}
			deck.Init()

			fmt.Printf("Your account balance is - $%d\n", player.Balance)
			fmt.Println()
			fmt.Println("Continue? [Y/n]")
			var action string
			fmt.Scan(&action)

			switch action {
			case "Y":
				{
					isNewRound = true
					clear()
					continue
				}
			case "n":
				{
					keepPlaying = false
					break
				}
			}
		}
		if !keepPlaying {
			break
		}
	}

}
