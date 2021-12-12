package minimax

import (
	"fmt"

	"github.com/joshtenorio/ninemo-battlesnake/datatypes"
	"github.com/joshtenorio/ninemo-battlesnake/snake/api"
)

func Minimax(board datatypes.Board, depth int, isUs bool, ourId string) int {
	id := ""
	if isUs {
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
	if depth == 0 {
		return Eval(board, id)
	}

	if isUs {
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
	return board
}

func Eval(board datatypes.Board, id string) int {
	fmt.Printf("hello!")
	return 0
}
