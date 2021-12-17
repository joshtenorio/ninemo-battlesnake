package minimax

import (
	"github.com/joshtenorio/ninemo-battlesnake/datatypes"
	"github.com/joshtenorio/ninemo-battlesnake/snake/api"
	"github.com/joshtenorio/ninemo-battlesnake/snake/astar"
)

/*
- board 		: the current position
- depth			: self explanatory
- maximizing 	: the current snake
- ourId 		: always refers to our snake, used for evaluating the position from our perspective
*/
func Minimax(board datatypes.Board, depth int, maximizing bool, ourId string) int {
	id := ""
	// if we are not maximizing, we want to look at our moves
	if !maximizing {
		id = ourId
	} else { // if we are maximizing, look at opponent's moves
		snakes := board.Snakes
		for i := 0; i < len(snakes); i++ {
			if snakes[i].ID != ourId {
				id = snakes[i].ID
				break
			}
		}
	}
	if depth == 0 || IsGameResolved(&board) {
		return Eval(board, ourId)
	}

	if maximizing {
		// from opponent's moves, get the maximum score
		var scores [4]int
		scores[0] = Minimax(MakeMove(id, "up", board), depth-1, false, ourId)
		scores[1] = Minimax(MakeMove(id, "down", board), depth-1, false, ourId)
		scores[2] = Minimax(MakeMove(id, "left", board), depth-1, false, ourId)
		scores[3] = Minimax(MakeMove(id, "right", board), depth-1, false, ourId)
		maxEval := api.GetMax(scores)
		return maxEval
	} else {
		// from our moves, get the minimum score
		var scores [4]int
		scores[0] = Minimax(MakeMove(id, "up", board), depth-1, true, ourId)
		scores[1] = Minimax(MakeMove(id, "down", board), depth-1, true, ourId)
		scores[2] = Minimax(MakeMove(id, "left", board), depth-1, true, ourId)
		scores[3] = Minimax(MakeMove(id, "right", board), depth-1, true, ourId)
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
	var us *datatypes.Battlesnake
	var opp *datatypes.Battlesnake // so we have quick reference
	for i := 0; i < len(board.Snakes); i++ {
		if board.Snakes[i].ID == id {
			us = &board.Snakes[i]
		} else {
			opp = &board.Snakes[i]
		}
	}
	// first, check if either of us are dead
	if us.Health == 0 {
		return -999
	} else if opp.Health == 0 {
		return 999
	}

	score := 0
	// if we are longer, prioritise going for their head
	if us.Length > opp.Length {
		score += 100
		// TODO: use the actual astar thing instead of euclidean distance
		score += (100 - astar.Heuristic(board, us.Head, opp.Head))
	} else { // else, prioritise getting to food
		score -= 100
		// get closest food
		food := datatypes.Coord{X: -99, Y: -99}
		for i := 0; i < len(board.Food); i++ {
			currdist := (us.Head.X-food.X)*(us.Head.X-food.X) + (us.Head.X-food.Y)*(us.Head.X-food.Y)
			newdist := (us.Head.X-board.Food[i].X)*(us.Head.X-board.Food[i].X) + (us.Head.Y-board.Food[i].Y)*(us.Head.Y-board.Food[i].Y)
			if newdist <= currdist {
				food = board.Food[i]
			}
		}
		score += (100 - astar.Heuristic(board, us.Head, food))
	}
	return score
}

func IsGameResolved(board *datatypes.Board) bool {
	for i := 0; i < len(board.Snakes); i++ {
		if board.Snakes[i].Health <= 0 {
			return true
		}
	}
	return false
}
