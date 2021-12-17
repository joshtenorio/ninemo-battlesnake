package snake

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/joshtenorio/ninemo-battlesnake/datatypes"
	"github.com/joshtenorio/ninemo-battlesnake/snake/api"
	"github.com/joshtenorio/ninemo-battlesnake/snake/floodfill"
	"github.com/joshtenorio/ninemo-battlesnake/snake/minimax"
)

func OnMove(w http.ResponseWriter, r *http.Request) {
	request := datatypes.GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// define our head
	head := request.You.Head
	move := "null"

	if len(request.Board.Snakes) > 2 {
		// if there is a potential head to head, go for it if we can win, else avoid
		move = DetectHeadToHead(&head, &request.Board, request.You.Length)

		// else, if we are in hazard and health is <=50, find the closest not-hazard square and move towards it if possible
		if move == "null" {
			move = HandleHazard(&head, int(request.You.Health), &request.Board)
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
			move = MoveInDirection(&head, &closestFood, &request.Board)
		}

		// if all other cases don't apply, pick a move that results in moving towards the most amount of space
		if move == "null" {
			move = HandleDefaultMove(&request.You.Head, &request.Board, 5, 6)
		}

	} else {
		// sadge
		ourId := request.You.ID
		scores := [4]int{0, 0, 0, 0}
		scores[0] = minimax.Minimax(minimax.MakeMove(ourId, "up", request.Board), 4, true, ourId)
		scores[1] = minimax.Minimax(minimax.MakeMove(ourId, "down", request.Board), 4, true, ourId)
		scores[2] = minimax.Minimax(minimax.MakeMove(ourId, "left", request.Board), 4, true, ourId)
		scores[3] = minimax.Minimax(minimax.MakeMove(ourId, "right", request.Board), 4, true, ourId)
		maxScore := api.GetMax(scores)
		index := 0
		for i := 0; i < len(scores); i++ {
			if maxScore == scores[i] {
				index = i
				break
			}
		}
		switch index {
		case 0:
			move = "up"
		case 1:
			move = "down"
		case 2:
			move = "left"
		case 3:
			move = "right"
		}

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

func IsMoveTrap(board *datatypes.Board, head *datatypes.Coord, move string, searchDist int, minSpaces int) bool {
	futurePos := api.MoveToCoord(move, head)
	numFreeSpaces := floodfill.CountFreeSpaces(board, futurePos, searchDist)
	if numFreeSpaces < minSpaces {
		return true
	} else {
		return false
	}
}

/*
returns true if move is physically possible, false if otherwise
*/
func IsMovePossible(head *datatypes.Coord, board *datatypes.Board, move string) bool {
	// calculate end position for move
	var position datatypes.Coord
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

	return !api.IsBlocking(board, position, false)
}

/*
returns valid move if we win a head-to-head, else returns a move that avoids it
if there is no head-to-head, return null
*/
func DetectHeadToHead(us *datatypes.Coord, board *datatypes.Board, ourLength int32) string {
	// get list of heads that are close to us (not including ourself)
	var heads []datatypes.Coord
	var lengths []int32
	snakes := board.Snakes
	for i := 0; i < len(snakes); i++ {
		if snakes[i].Head.X != us.X || snakes[i].Head.Y != us.Y { // don't include ourself in list of snakes
			heads = append(heads, snakes[i].Head)
			lengths = append(lengths, snakes[i].Length)
		}
	}
	// iterate through all the heads, if d^2 is == 2 or 4 then there is a possibility of h2h
	// find first head that matches the above condition
	// limitation: only considers one possible h2h at a time - if there are >1 possible h2h i only consider one for now
	var enemyHead = datatypes.Coord{X: -1, Y: -1}
	var enemyLength int32
	for i := 0; i < len(heads); i++ {
		distSquared := (heads[i].X-us.X)*(heads[i].X-us.X) + (heads[i].Y-us.Y)*(heads[i].Y-us.Y)
		if distSquared == 2 || distSquared == 4 {
			if distSquared == 4 && api.IsBlocking(board, datatypes.Coord{X: (us.X + heads[i].X) / 2, Y: (us.Y + heads[i].Y) / 2}, false) {
				// special case for d^2=4: make sure there isn't a body between us
				continue
			} else {
				enemyHead = heads[i]
				enemyLength = lengths[i]
				break
			}

		}
	} // end for loop

	// if enemyhead is -1 return null because there is no h2h collision possible
	if enemyHead.X == -1 {
		return "null"
	}

	// before continuing, determine all possible moves since we need it for both cases
	movesUs := []datatypes.Coord{{X: us.X, Y: us.Y + 1}, {X: us.X, Y: us.Y - 1}, {X: us.X - 1, Y: us.Y}, {X: us.X + 1, Y: us.Y}}
	movesEnemy := []datatypes.Coord{
		{X: enemyHead.X, Y: enemyHead.Y + 1},
		{X: enemyHead.X, Y: enemyHead.Y - 1},
		{X: enemyHead.X - 1, Y: enemyHead.Y},
		{X: enemyHead.X + 1, Y: enemyHead.Y}}

	// determine if we can beat them
	if ourLength < enemyLength {
		fmt.Printf("in h2h: we lose so avoid\n")
		// pick something that avoids them because we'll lose
		for i := 0; i < len(movesUs); i++ {
			futureUs := movesUs[i]
			escapes := true
			for j := 0; j < len(movesEnemy); j++ {
				futureEnemy := movesEnemy[j]
				if (futureUs.X == futureEnemy.X && futureUs.Y == futureEnemy.Y) || !IsMovePossible(us, board, api.IndexToMove(i)) {
					escapes = false
				}
			} // end for j
			if escapes {
				return api.IndexToMove(i)
			}
		} // end for i
	} else if ourLength > enemyLength { // if we are > we win
		fmt.Printf("in h2h: we win so attempt\n")
		// pick the move that results in h2h collision
		// if there is food adjacent, move to eat the food instead
		foodAdjacent, move := api.IsFoodAdjacent(board, *us)
		if foodAdjacent {
			return move
		}

		// iterate through all possible moves, if it results in enemy collision then return that move
		for i := 0; i < len(movesUs); i++ {
			futureUs := movesUs[i]
			for j := 0; j < len(movesEnemy); j++ {
				futureEnemy := movesEnemy[j]
				if (futureUs.X == futureEnemy.X && futureUs.Y == futureEnemy.Y) && IsMovePossible(us, board, api.IndexToMove(i)) {
					return api.IndexToMove(i)
				}
			} // end for j
		} // end for i
	} else { // if lengths are == and there is a food adjacent to us, go for the food even if h2h collision possible
		foodAdjacent, move := api.IsFoodAdjacent(board, *us)
		if foodAdjacent {
			return move
		} else { // if food is not adjacent, avoid collision
			for i := 0; i < len(movesUs); i++ {
				futureUs := movesUs[i]
				escapes := true
				for j := 0; j < len(movesEnemy); j++ {
					futureEnemy := movesEnemy[j]
					if (futureUs.X == futureEnemy.X && futureUs.Y == futureEnemy.Y) || !IsMovePossible(us, board, api.IndexToMove(i)) {
						escapes = false
					}
				} // end for j
				if escapes {
					return api.IndexToMove(i)
				}
			} // end for i
		}
	} // end ourLength __ enemyLength else
	return "null"
}

func MoveInDirection(head *datatypes.Coord, target *datatypes.Coord, board *datatypes.Board) string {
	var dx, dy int = target.X - head.X, target.Y - head.Y
	if dx > 0 && IsMovePossible(head, board, "right") && !IsMoveTrap(board, head, "right", 5, 6) {
		return "right"
	} else if dx < 0 && IsMovePossible(head, board, "left") && !IsMoveTrap(board, head, "left", 5, 6) {
		return "left"
	} else if dy > 0 && IsMovePossible(head, board, "up") && !IsMoveTrap(board, head, "up", 5, 6) {
		return "up"
	} else if dy < 0 && IsMovePossible(head, board, "down") && !IsMoveTrap(board, head, "down", 5, 6) {
		return "down"
	} else {
		return "null"
	}
}

func HandleHazard(head *datatypes.Coord, health int, board *datatypes.Board) string {
	if api.IsHazard(board, *head) && health <= 60 {

		// first, check if there is food adjacent
		foodAdjacent, move := api.IsFoodAdjacent(board, *head)
		if foodAdjacent {
			return move
		}

		// find closest non-hazard square
		// check a 5x5 region around our head for the closest non-hazard square (25 loops)
		safeCoord := datatypes.Coord{X: -1, Y: -1}
		distSquared := 90000 // TODO: change this to int's max value
		for i := head.X - 2; i < head.X+2; i++ {
			for j := head.Y - 2; j < head.Y+2; j++ {
				// sanity check i and j
				if i < 0 || i >= board.Width || j < 0 || j >= board.Height {
					continue
				}
				var dx, dy int = i - head.X, j - head.Y
				if !api.IsHazard(board, datatypes.Coord{X: i, Y: j}) && (dx*dx+dy*dy < distSquared) {
					safeCoord = datatypes.Coord{X: i, Y: j}
					distSquared = dx*dx + dy*dy
				}
			} // end for j
		} // end for i

		return MoveInDirection(head, &safeCoord, board)
	}
	return "null"
}

func HandleDefaultMove(head *datatypes.Coord, board *datatypes.Board, searchDist int, minSpaces int) string {
	// if all other cases don't apply, pick a move that results in moving towards the most amount of space
	var spaces []int
	var legalMoves []string
	if IsMovePossible(head, board, "up") && !IsMoveTrap(board, head, "up", searchDist, minSpaces) {
		legalMoves = append(legalMoves, "up")
		spaces = append(spaces, floodfill.CountFreeSpaces(board, api.MoveToCoord("up", head), searchDist))
	} else if IsMovePossible(head, board, "down") && !IsMoveTrap(board, head, "down", searchDist, minSpaces) {
		legalMoves = append(legalMoves, "down")
		spaces = append(spaces, floodfill.CountFreeSpaces(board, api.MoveToCoord("down", head), searchDist))
	} else if IsMovePossible(head, board, "left") && !IsMoveTrap(board, head, "left", searchDist, minSpaces) {
		legalMoves = append(legalMoves, "left")
		spaces = append(spaces, floodfill.CountFreeSpaces(board, api.MoveToCoord("left", head), searchDist))
	} else if IsMovePossible(head, board, "right") && !IsMoveTrap(board, head, "right", searchDist, minSpaces) {
		legalMoves = append(legalMoves, "right")
		spaces = append(spaces, floodfill.CountFreeSpaces(board, api.MoveToCoord("right", head), searchDist))
	}

	bestMove := "null"
	maxSpaces := -1
	for i := 0; i < len(legalMoves); i++ {
		if spaces[i] > maxSpaces {
			bestMove = legalMoves[i]
			maxSpaces = spaces[i]
		}
	}
	return bestMove
}
