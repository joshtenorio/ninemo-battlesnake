package minimax

import (
	"fmt"

	"github.com/joshtenorio/ninemo-battlesnake/datatypes"
	"github.com/joshtenorio/ninemo-battlesnake/snake/api"
)

func Minimax(board datatypes.Board, depth int, maximizing bool, ourId string) int {
	id := ""
	if maximizing {
		id = ourId
	} else {
		snakes := board.Snakes
		for i := 0; i < len(snakes); i++ {
			if snakes[i].ID != ourId {
				id = snakes[i].ID
				break
			}
		}
	}
	if depth == 0 || IsGameResolved(&board) {
		return Eval(board, id)
	}

	if maximizing {
		var scores [4]int
		scores[0] = Minimax(MakeMove(id, "up", board), depth-1, false, id)
		scores[1] = Minimax(MakeMove(id, "down", board), depth-1, false, id)
		scores[2] = Minimax(MakeMove(id, "left", board), depth-1, false, id)
		scores[3] = Minimax(MakeMove(id, "right", board), depth-1, false, id)
		maxEval := api.GetMax(scores)
		return maxEval
	} else {
		var scores [4]int
		scores[0] = Minimax(MakeMove(id, "up", board), depth-1, true, id)
		scores[1] = Minimax(MakeMove(id, "down", board), depth-1, true, id)
		scores[2] = Minimax(MakeMove(id, "left", board), depth-1, true, id)
		scores[3] = Minimax(MakeMove(id, "right", board), depth-1, true, id)
		minEval := api.GetMin(scores)
		return minEval
	}
}

func MakeMove(id string, move string, board datatypes.Board) datatypes.Board {
	updatedBoard := board
	var us *datatypes.Battlesnake
	var opp *datatypes.Battlesnake // so we have quick reference
	for i := 0; i < len(updatedBoard.Snakes); i++ {
		if updatedBoard.Snakes[i].ID == id {
			us = &updatedBoard.Snakes[i]
		} else {
			opp = &updatedBoard.Snakes[i]
		}
	}

	// get final position of snake and update body
	posFinal := api.MoveToCoord(move, &us.Head)
	us.Head = posFinal
	for i := len(us.Body) - 1; i >= 1; i-- {
		us.Body[i] = us.Body[i-1]
	}
	us.Health -= 1

	// check if in hazard
	if api.IsHazard(&updatedBoard, posFinal) {
		us.Health -= 15
	}
	// check for collision with wall, our body, or opponent body
	if api.IsBlocking(&updatedBoard, posFinal, true) {
		// we are dead lol
		us.Health = 0
		return updatedBoard
	}

	// resolve head to head collision
	if posFinal == opp.Head {
		if us.Length > opp.Length {
			opp.Health = 0
		} else if us.Length < opp.Length {
			us.Health = 0
		} else {
			opp.Health = 0
			us.Health = 0
		}
		return updatedBoard
	}

	// get food
	if api.IsFood(&updatedBoard, posFinal) {
		foodIndex := 0
		for i := 0; i < len(updatedBoard.Food); i++ {
			if posFinal.X == updatedBoard.Food[i].X && posFinal.Y == updatedBoard.Food[i].Y {
				foodIndex = i
				break
			}
		}
		us.Health = 100
		us.Length += 1
		us.Body = append(us.Body, datatypes.Coord{X: us.Body[len(us.Body)-1].X, Y: us.Body[len(us.Body)-1].Y})
		updatedBoard.Food = api.RemoveCoord(updatedBoard.Food, foodIndex)
	}

	return updatedBoard
}

func Eval(board datatypes.Board, id string) int {
	fmt.Printf("hello!")
	return 0
}

func IsGameResolved(board *datatypes.Board) bool {
	for i := 0; i < len(board.Snakes); i++ {
		if board.Snakes[i].Health <= 0 {
			return true
		}
	}
	return false
}
