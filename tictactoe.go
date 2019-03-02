package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Tile is used as the cells that the board is made from
type Tile int

// Board is used as the tic tac toe board
type Board [][]Tile

// Tiles can have any of these values
const (
	BlankVar = iota
	XVar     = iota
	OVar     = iota
)

var (
	dbport = os.Getenv("POSTGRES_PORT_5432_TCP_PORT")
	dbuser = os.Getenv("POSTGRES_ENV_POSTGRES_USER")
	dbpass = os.Getenv("POSTGRES_ENV_POSTGRES_PASSWORD")
	dbname = os.Getenv("POSTGRES_ENV_POSTGRES_DB")
)

var db *sql.DB

func initDb() *sql.DB {
	var err error
	db, err = sql.Open(
		"postgres",
		fmt.Sprintf("host=postgres port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbport, dbuser, dbpass, dbname),
	)
	if err != nil {
		log.Fatalf("The data source arguments are not valid.")
	}

	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err != nil {
			log.Printf("The database is not up yet")
			if i == 9 {
				panic(err)
			}
			time.Sleep(6 * time.Second)
		} else {
			break
		}
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS tictactoe (id integer, board jsonb, winner integer, turn integer, players integer);")
	if err != nil {
		fmt.Println("Cannot create table")
	}

	return db
}

// Game tick tack toe board
type Game struct {
	Winner    Tile
	Turn      Tile
	SessionID int
	Players   int
	Board     Board
}

// Update the game
func (g *Game) Update(r, c Tile) (*Game, error) {
	var err error
	if g.Board[r][c] != BlankVar { // If the choice was already chosen then we no updates are needed
		return g, errors.New("invalid move that position has already been taken")
	}

	if g.Winner != BlankVar {
		return g, errors.New("invalid move a winner has already been selected")
	}

	g.Board[r][c] = g.Turn

	g.checkForWinner()

	type gameJSON struct {
		Board [][]Tile `json:"board"`
	}
	marshalJSON := gameJSON{g.Board}
	boardJsonb, err := json.Marshal(marshalJSON)
	if err != nil {
		log.Fatal(err)
	}

	g.passTurn()

	_, err = db.Exec(
		"UPDATE tictactoe SET board=$1, winner=$2, turn=$3, players=$4 WHERE id=$5;",
		string(boardJsonb), g.Winner, g.Turn, g.Players, g.SessionID,
	)
	if err != nil {
		log.Println(err)
	}

	return g, nil
}

// passTurn changes the turn from "X" to "O" or vice versa
func (g *Game) passTurn() {
	if g.Turn == XVar {
		g.Turn = OVar
	} else if g.Turn == OVar {
		g.Turn = XVar
	} else {
		log.Println(string(g.Turn) + " is invalid")
	}
}

// DeepCopy returns a deep copy of game
func (g *Game) DeepCopy() *Game {
	newBoard := make([][]Tile, 3)
	for i := 0; i < 3; i++ {
		newBoard[i] = make([]Tile, 3)
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			newBoard[i][j] = g.Board[i][j]
		}
	}

	return &Game{
		g.Winner,
		g.Turn,
		g.SessionID,
		g.Players,
		newBoard,
	}
}

func (g *Game) checkForWinner() {
	winner := BlankVar

	// Horizontal and vertical checks
	for i := 0; i < 3; i++ {
		xRowCheck := true
		xColumnCheck := true
		oRowCheck := true
		oColumnCheck := true
		for j := 0; j < 3; j++ {
			if g.Board[i][j] != XVar {
				xRowCheck = false
			}
			if g.Board[i][j] != OVar {
				oRowCheck = false
			}
			if g.Board[j][i] != XVar {
				xColumnCheck = false
			}
			if g.Board[j][i] != OVar {
				oColumnCheck = false
			}
		}
		if xColumnCheck || xRowCheck {
			winner = XVar
		}
		if oColumnCheck || oRowCheck {
			winner = OVar
		}
	}

	// Diagonal checks
	if g.Board[0][2] == XVar && g.Board[1][1] == XVar && g.Board[2][0] == XVar {
		winner = XVar
	} else if g.Board[0][2] == OVar && g.Board[1][1] == OVar && g.Board[2][0] == OVar {
		winner = OVar
	}
	if g.Board[0][0] == XVar && g.Board[1][1] == XVar && g.Board[2][2] == XVar {
		winner = XVar
	} else if g.Board[0][0] == OVar && g.Board[1][1] == OVar && g.Board[2][2] == OVar {
		winner = OVar
	}

	g.Winner = Tile(winner)
}

func (g *Game) getComputerChoice() (Tile, Tile, error) {
	var availableSpaces [][]int
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.Board[i][j] == BlankVar {
				availableSpaces = append(availableSpaces, []int{i, j})
				testGameO := g.DeepCopy()
				testGameO.Board[i][j] = OVar
				testGameO.checkForWinner()
				if testGameO.Winner != Tile(BlankVar) {
					return Tile(i), Tile(j), nil
				}
			}
		}
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.Board[i][j] == BlankVar {
				testGameX := g.DeepCopy()
				testGameX.Board[i][j] = XVar
				testGameX.checkForWinner()
				if testGameX.Winner != Tile(BlankVar) {
					return Tile(i), Tile(j), nil
				}
			}
		}
	}
	rand.Seed(time.Now().Unix())
	if len(availableSpaces) > 0 {
		n := rand.Int() % len(availableSpaces)
		return Tile(availableSpaces[n][0]), Tile(availableSpaces[n][1]), nil
	}
	return 0, 0, errors.New("No available spaces")
}

// NewGame returns an empty board object
func NewGame(players int) (int, *Game) {
	statement := "SELECT id FROM tictactoe ORDER BY id DESC LIMIT 1;"
	row := db.QueryRow(statement)
	var id int
	err := row.Scan(&id)
	if err != nil {
		log.Println(err)
	}
	id = id + 1
	log.Println("New Game D is: ", id)

	game := &Game{
		BlankVar,
		XVar,
		id,
		players,
		Board{
			[]Tile{BlankVar, BlankVar, BlankVar},
			[]Tile{BlankVar, BlankVar, BlankVar},
			[]Tile{BlankVar, BlankVar, BlankVar},
		},
	}

	type gameJSON struct {
		Board [][]Tile `json:"board"`
	}
	marshalJSON := gameJSON{game.Board}
	boardJsonb, err := json.Marshal(marshalJSON)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(
		"INSERT INTO tictactoe(winner, id, turn, players, board) VALUES ($1, $2, $3, $4, $5);",
		game.Winner, game.SessionID, game.Turn, game.Players, string(boardJsonb),
	)
	if err != nil {
		log.Println(err)
	}

	return id, game
}

// GetGame retuns the game based off of the id
func GetGame(id int) *Game {
	var boardJsonb []uint8
	var board Board
	var players int
	var turn, winner Tile
	board = make([][]Tile, 3)
	for i := 0; i < 3; i++ {
		board[i] = make([]Tile, 3)
	}

	statement := "SELECT board, turn, winner, players FROM tictactoe WHERE id=$1 LIMIT 1;"
	row := db.QueryRow(statement, id)
	err := row.Scan(&boardJsonb, &turn, &winner, &players)
	if err != nil {
		log.Println(err)
	}
	var i map[string]interface{}
	json.Unmarshal(boardJsonb, &i)
	for name, v := range i {
		if name == "board" {
			for r, e := range v.([]interface{}) {
				for c, f := range e.([]interface{}) {
					var tile Tile
					f, ok := f.(float64)
					if ok {
						tile = Tile(f)
					}
					board[r][c] = tile
				}
			}
		}
	}
	return &Game{winner, turn, id, players, board}
}
