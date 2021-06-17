package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/joshtenorio/ninemo-bot/datatypes"
	"github.com/joshtenorio/ninemo-bot/snake"
	"github.com/joshtenorio/ninemo-bot/snake/api"
)

// HandleIndex is called when your Battlesnake is created and refreshed
// by play.battlesnake.com. BattlesnakeInfoResponse contains information about
// your Battlesnake, including what it should look like on the game board.
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	response := datatypes.BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "tenmo",
		Color:      "#4287f5",
		Head:       "evil",
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
	request := datatypes.GameRequest{}
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
	request := datatypes.GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// define our head
	head := request.You.Head

	// define list of legal moves
	var legalMoves []string
	if snake.IsMovePossible(&head, &request.Board, "up") {
		legalMoves = append(legalMoves, "up")
	} else if snake.IsMovePossible(&head, &request.Board, "down") {
		legalMoves = append(legalMoves, "down")
	} else if snake.IsMovePossible(&head, &request.Board, "left") {
		legalMoves = append(legalMoves, "left")
	} else if snake.IsMovePossible(&head, &request.Board, "right") {
		legalMoves = append(legalMoves, "right")
	}

	// if there is a potential head to head, go for it if we can win, else avoid
	move := "null"
	move = snake.DetectHeadToHead(&head, &request.Board, request.You.Length, legalMoves)
	// else, if we are in hazard and health is <=50, find the closest not-hazard square and move towards it if possible
	if move == "null" {
		// check if head is in a hazard and health is <= half
		// TODO: move this to snake.go
		if api.IsHazard(&request.Board, request.You.Head) && request.You.Health <= 60 {
			// find closest non-hazard square
			// check a 5x5 region around our head for the closest non-hazard square (25 loops)
			safeCoord := datatypes.Coord{X: -1, Y: -1}
			distSquared := 90000 // TODO: change this to int's max value
			for i := head.X - 2; i < head.X+2; i++ {
				for j := head.Y - 2; j < head.Y+2; j++ {
					// sanity check i and j
					if i < 0 || i >= request.Board.Width || j < 0 || j >= request.Board.Height {
						continue
					}
					var dx, dy int = i - head.X, j - head.Y
					if !api.IsHazard(&request.Board, datatypes.Coord{X: i, Y: j}) && (dx*dx+dy*dy < distSquared) {
						safeCoord = datatypes.Coord{X: i, Y: j}
						distSquared = dx*dx + dy*dy
					}
				} // end for j
			} // end for i

			var dx, dy int = safeCoord.X - head.X, safeCoord.Y - head.Y
			if dx > 0 && snake.IsMovePossible(&head, &request.Board, "right") {
				move = "right"
			} else if dx < 0 && snake.IsMovePossible(&head, &request.Board, "left") {
				move = "left"
			} else if dy > 0 && snake.IsMovePossible(&head, &request.Board, "up") {
				move = "up"
			} else if dy < 0 && snake.IsMovePossible(&head, &request.Board, "down") {
				move = "down"
			}
		} // end if "in hazard and health low"
	} // end if move == null

	// else, find closest food and path to it if possible
	if move == "null" {
		// TODO: put this in a function in snake.go
		dist := 90000 // TODO: change this to actual max value of int, lookup golang spec
		food := request.Board.Food
		var closestFood datatypes.Coord
		for i := 0; i < len(food); i++ {
			if (food[i].X-head.X)*(food[i].X-head.X)+(food[i].Y-head.Y)*(food[i].Y-head.Y) < dist {
				dist = (food[i].X-head.X)*(food[i].X-head.X) + (food[i].Y-head.Y)*(food[i].Y-head.Y)
				closestFood = food[i]
			}
		}

		// attempt to go in the direction of the closestFood
		var dx, dy int = closestFood.X - head.X, closestFood.Y - head.Y
		if dx > 0 && snake.IsMovePossible(&head, &request.Board, "right") {
			move = "right"
		} else if dx < 0 && snake.IsMovePossible(&head, &request.Board, "left") {
			move = "left"
		} else if dy > 0 && snake.IsMovePossible(&head, &request.Board, "up") {
			move = "up"
		} else if dy < 0 && snake.IsMovePossible(&head, &request.Board, "down") {
			move = "down"
		}
	}

	// if we can't go in direction of closest food, just pick a random move
	if move == "null" {
		move = legalMoves[rand.Intn(len(legalMoves))]
	}

	// set up response
	response := datatypes.MoveResponse{
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
	request := datatypes.GameRequest{}
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
