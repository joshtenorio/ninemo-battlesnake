package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	Height int           `json:"height"`
	Width  int           `json:"width"`
	Food   []Coord       `json:"food"`
	Snakes []Battlesnake `json:"snakes"`
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
// TODO: Use the information in the GameRequest object to determine your next move.
func HandleMove(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// get board size and current position of our head
	// board size info is in here and not HandleStart in case we are playing two different games at once with different board sizes
	var xMin, yMin int = 0, 0
	var xMax int = request.Board.Width
	var yMax int = request.Board.Height

	var xHead int = request.You.Head.X
	var yHead int = request.You.Head.Y

	// declare list of legal moves
	legalMoves := []string{"up", "down", "left", "right"}

	// eliminate illegal moves
	var (
		endUp    = Coord{xHead, yHead + 1}
		endDown  = Coord{xHead, yHead - 1}
		endLeft  = Coord{xHead - 1, yHead}
		endRight = Coord{xHead + 1, yHead}
	)

	// eliminate moves that result in colliding with self
	body := request.You.Body
	for i := 0; i < len(body); i++ {
		coord := body[i]
		if coord.X == endUp.X && coord.Y == endUp.Y {
			legalMoves[0] = "null"
			fmt.Printf("body is above head\n")
		} else if coord.X == endDown.X && coord.Y == endDown.Y {
			legalMoves[1] = "null"
			fmt.Printf("body is below head\n")
		} else if coord.X == endLeft.X && coord.Y == endLeft.Y {
			legalMoves[2] = "null"
			fmt.Printf("body is left of head\n")
		} else if coord.X == endRight.X && coord.Y == endRight.Y {
			legalMoves[3] = "null"
			fmt.Printf("body is right of head\n")
		}
	}

	// eliminate moves  that result in colliding with wall
	for i := 0; i < len(legalMoves); i++ {
		if legalMoves[i] != "null" {
			switch legalMoves[i] {
			case "up":
				if endUp.Y >= yMax {
					legalMoves[i] = "null"
				}
			case "down":
				if endDown.Y < yMin {
					legalMoves[i] = "null"
				}
			case "left":
				if endLeft.X < xMin {
					legalMoves[i] = "null"
				}
			case "right":
				if endRight.X >= xMax {
					legalMoves[i] = "null"
				}
			} // end switch
		}
	} // end for

	// select a legal move that isn't null
	//move := legalMoves[rand.Intn(len(legalMoves))]
	// pick the first move that isn't null
	move := "null"
	for i := 0; i < len(legalMoves); i++ {
		if legalMoves[i] != "null" {
			move = legalMoves[i]
		}
	}
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
