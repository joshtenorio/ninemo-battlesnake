package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
)

type Game struct {
	ID      string `json:"id"`
	Timeout int32  `json:"timeout"`
}

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Battlesnake struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Health int32   `json:"health"`
	Body   []Coord `json:"body"`
	Head   Coord   `json:"head"`
	Length int32   `json:"length"`
	Shout  string  `json:"shout"`
}

type Board struct {
	Height  int           `json:"height"`
	Width   int           `json:"width"`
	Food    []Coord       `json:"food"`
	Snakes  []Battlesnake `json:"snakes"`
	Hazards []Coord       `json:"hazards"`
}

type BattlesnakeInfoResponse struct {
	APIVersion string `json:"apiversion"`
	Author     string `json:"author"`
	Color      string `json:"color"`
	Head       string `json:"head"`
	Tail       string `json:"tail"`
}

type GameRequest struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`
}

type MoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout,omitempty"`
}

/*
returns true if move is possible, false if otherwise
doesn't take into consideration head to head collisions
*/
func isMovePossible(head *Coord, board *Board, move string) bool {
	// calculate end position for move
	var position Coord
	switch move {
	case "up":
		position.X = head.X
		position.Y = head.Y + 1
	case "down":
		position.X = head.X
		position.Y = head.Y - 1
	case "left":
		position.X = head.X - 1
		position.Y = head.Y
	case "right":
		position.X = head.X + 1
		position.Y = head.Y
	}

	// check if move collides with wall
	xMax, yMax := board.Width, board.Height
	if position.X >= xMax || position.X < 0 || position.Y >= yMax || position.Y < 0 {
		return false
	}

	// check if move collides with other snakes
	for i := 0; i < len(board.Snakes); i++ {
		coord := board.Snakes[i].Head
		if position.X == coord.X && position.Y == coord.Y {
			return false
		}
		for j := 0; j < len(board.Snakes[i].Body); j++ {
			coord = board.Snakes[i].Body[j]
			if position.X == coord.X && position.Y == coord.Y {
				return false
			}
		} // end inner for
	} // end outer for
	return true
}

/*
returns valid move if we win a head-to-head, else returns a move that avoids it
if there is no head-to-head, return null
*/
func detectHeadToHead(us *Coord, board *Board, validMoves []string) string {
	// get list of heads that are close to us (not including ourself)

	return "null"
}

// HandleIndex is called when your Battlesnake is created and refreshed
// by play.battlesnake.com. BattlesnakeInfoResponse contains information about
// your Battlesnake, including what it should look like on the game board.
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	response := BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "tenmo",
		Color:      "#4287f5",
		Head:       "default",
		Tail:       "default",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

// HandleStart is called at the start of each game your Battlesnake is playing.
// The GameRequest object contains information about the game that's about to start.
// TODO: Use this function to decide how your Battlesnake is going to look on the board.
func HandleStart(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// Nothing to respond with here
	fmt.Print("START\n")
}

// HandleMove is called for each turn of each game.
// Valid responses are "up", "down", "left", or "right".
func HandleMove(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// define our head
	head := request.You.Head

	// define list of legal moves
	var legalMoves []string
	if isMovePossible(&head, &request.Board, "up") {
		legalMoves = append(legalMoves, "up")
	} else if isMovePossible(&head, &request.Board, "down") {
		legalMoves = append(legalMoves, "down")
	} else if isMovePossible(&head, &request.Board, "left") {
		legalMoves = append(legalMoves, "left")
	} else if isMovePossible(&head, &request.Board, "right") {
		legalMoves = append(legalMoves, "right")
	}

	// if there is a potential head to head, go for it if we can win, else avoid
	move := "null"
	move = detectHeadToHead(&head, &request.Board, legalMoves)

	// else, if we are in hazard and health is <=50, find the closest not-hazard square and move towards it if possible
	if move == "null" {
		// put hazard code in here
	}
	// else, find closest food and path to it if possible
	if move == "null" {
		dist := 90000 // TODO: change this to actual max value of int, lookup golang spec
		food := request.Board.Food
		var closestFood Coord
		for i := 0; i < len(food); i++ {
			if (food[i].X-head.X)*(food[i].X-head.X)+(food[i].Y-head.Y)*(food[i].Y-head.Y) < dist {
				dist = (food[i].X-head.X)*(food[i].X-head.X) + (food[i].Y-head.Y)*(food[i].Y-head.Y)
				closestFood = food[i]
			}
		}

		// attempt to go in the direction of the closestFood
		var dx, dy int = closestFood.X - head.X, closestFood.Y - head.Y
		if dx > 0 && isMovePossible(&head, &request.Board, "right") {
			move = "right"
		} else if dx < 0 && isMovePossible(&head, &request.Board, "left") {
			move = "left"
		} else if dy > 0 && isMovePossible(&head, &request.Board, "up") {
			move = "up"
		} else if dy < 0 && isMovePossible(&head, &request.Board, "down") {
			move = "down"
		}
	}

	// if we can't go in direction of closest food, just pick a random move
	if move == "null" {
		move = legalMoves[rand.Intn(len(legalMoves))]
	}

	// set up response
	response := MoveResponse{
		Move: move,
	}

	fmt.Printf("CHOSEN MOVE: %s\n\n", response.Move)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

// HandleEnd is called when a game your Battlesnake was playing has ended.
// It's purely for informational purposes, no response required.
func HandleEnd(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// Nothing to respond with here
	fmt.Print("END\n")
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/start", HandleStart)
	http.HandleFunc("/move", HandleMove)
	http.HandleFunc("/end", HandleEnd)

	fmt.Printf("Starting Battlesnake Server at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
