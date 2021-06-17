package api

import (
	"github.com/joshtenorio/ninemo-bot/datatypes"
)

/*
returns move corresponding to index
0: up
1: down
2: left
3: right
default: null
*/
func IndexToMove(num int) string {
	switch num {
	case 0:
		return "up"
	case 1:
		return "down"
	case 2:
		return "left"
	case 3:
		return "right"
	default:
		return "null"
	}
}

/*
given a move and initial position, calculate final position
*/
func MoveToCoord(move string, position *datatypes.Coord) datatypes.Coord {
	output := datatypes.Coord{X: -1, Y: -1}
	switch move {
	case "up":
		output.X = position.X
		output.Y = position.Y + 1
	case "down":
		output.X = position.X
		output.Y = position.Y - 1
	case "left":
		output.X = position.X - 1
		output.Y = position.Y
	case "right":
		output.X = position.X + 1
		output.Y = position.Y
	}
	return output
}

/*
checks if pos is blocking (snake body or wall)
*/
func IsBlocking(board *datatypes.Board, pos datatypes.Coord) bool {
	// check if snake occupies pos
	for i := 0; i < len(board.Snakes); i++ {
		for j := 0; j < len(board.Snakes[i].Body)-1; j++ {
			if board.Snakes[i].Body[j].X == pos.X && board.Snakes[i].Body[j].Y == pos.Y {
				return true
			}
		}

	} // end outer for

	// check if wall
	if pos.X < 0 || pos.Y < 0 || pos.X >= board.Width || pos.Y >= board.Height {
		return true
	}

	return false
}

/*
checks if pos is a hazard
*/
func IsHazard(board *datatypes.Board, pos datatypes.Coord) bool {
	for i := 0; i < len(board.Hazards); i++ {
		if board.Hazards[i].X == pos.X && board.Hazards[i].Y == pos.Y {
			return true
		}
	}
	return false
}

/*
checks if there is food adjacent
*/
func IsFoodAdjacent(board *datatypes.Board, pos datatypes.Coord) (adjacent bool, move string) {
	food := board.Food
	for i := 0; i < len(food); i++ {
		x, y := food[i].X, food[i].Y
		if (x-pos.X*x-pos.X)+(y-pos.Y*y-pos.Y) == 1 { // if d^2 == 1 there is food adjacent
			adjacent = true
			// find the move that results in eating the food
			for j := 0; j < 4; j++ {
				futurePos := MoveToCoord(IndexToMove(j), &pos)
				if futurePos.X == x && futurePos.Y == y {
					move = IndexToMove(j)
					return
				}
			}
		}
	}
	adjacent = false
	move = "null"
	return
}
