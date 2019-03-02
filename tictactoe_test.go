package main

import (
	"testing"
)

type gameMoves struct {
	moves          []move
	expectedWinner Tile
}

type move struct {
	x Tile
	y Tile
}

func TestGames(test *testing.T) {
	initDb()

	var games []gameMoves
	games = append(games, gameMoves{ // Game 0, test vertical winner
		[]move{
			move{0, 0}, //x
			move{1, 1},
			move{1, 0}, //x
			move{2, 2},
			move{2, 0}, //x
		},
		XVar,
	})

	games = append(games, gameMoves{ // Game 1, test diagonal winner
		[]move{
			move{0, 0}, //x
			move{0, 1},
			move{1, 1}, //x
			move{1, 0},
			move{2, 2}, //x
		},
		XVar,
	})

	games = append(games, gameMoves{ // Game 2, test horizontal winner
		[]move{
			move{0, 0}, //x
			move{0, 2},
			move{1, 1}, //x
			move{1, 2},
			move{1, 0}, //x
			move{2, 2},
		},
		OVar,
	})

	games = append(games, gameMoves{ // Game 2, test cats game
		[]move{
			move{0, 0}, //x
			move{1, 0},
			move{2, 0}, //x
			move{2, 1},
			move{0, 1}, //x
			move{0, 2},
			move{2, 2}, //x
			move{1, 1},
			move{1, 2}, //x
		},
		BlankVar,
	})

	for i, game := range games {
		j, g := NewGame(1)
		for _, move := range game.moves {
			var err error
			g, err = g.Update(move.x, move.y)
			if err != nil {
				test.Errorf("test %d: update failed for game id %d", i, j)
			}
		}
		if g.Winner != game.expectedWinner {
			test.Errorf("test %d: winner not selected for game id %d", i, j)
		}
	}

	db.Close()
}

func TestComputer(test *testing.T) {
	initDb()

	turn := 0
	for gameNumber := 0; gameNumber < 10; gameNumber++ {
		i, g := NewGame(1)
		for g.Winner != BlankVar && turn < 9 {
			x, y, err := g.getComputerChoice()
			if err != nil {
				test.Errorf("test %d: getComputerChoice failed for game id %d", gameNumber, i)
			}
			g, err = g.Update(x, y)
			if err != nil {
				test.Errorf("test %d: update failed for game id %d", gameNumber, i)
			}
		}
	}

	db.Close()
}
