package snake

import (
	"fmt"

	"github.com/joshtenorio/ninemo-bot/datatypes"
	//"github.com/joshtenorio/ninemo-bot/floodfill"
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
checks if pos is blocking
*/
func IsBlocking(board *datatypes.Board, pos datatypes.Coord) bool {
	// check if snake
	for i := 0; i < len(board.Snakes); i++ {
		head := board.Snakes[i].Head
		if head.X == pos.X && head.Y == pos.Y {
			return true
		}
		for j := 0; j < len(board.Snakes[i].Body); j++ {
			if board.Snakes[i].Body[j].X == pos.X && board.Snakes[i].Body[j].Y == pos.Y {
				return true
			}
		}
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
returns true if move is possible, false if otherwise
doesn't take into consideration head to head collisions
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
func DetectHeadToHead(us *datatypes.Coord, board *datatypes.Board, ourLength int32, validMoves []string) string {
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
			if distSquared == 4 && IsBlocking(board, datatypes.Coord{X: (us.X + heads[i].X) / 2, Y: (us.Y + heads[i].Y) / 2}) { // special case for d^2=4: make sure there isn't a body between us
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
	if ourLength <= enemyLength {
		fmt.Printf("in h2h: we lose so avoid\n")
		// pick something that avoids them because we'll lose
		for i := 0; i < len(movesUs); i++ {
			futureUs := movesUs[i]
			escapes := true
			for j := 0; j < len(movesEnemy); j++ {
				futureEnemy := movesEnemy[j]
				if (futureUs.X == futureEnemy.X && futureUs.Y == futureEnemy.Y) || !IsMovePossible(us, board, IndexToMove(i)) {
					escapes = false
				}
			} // end for j
			if escapes {
				return IndexToMove(i)
			}
		} // end for i
	} else if ourLength > enemyLength {
		fmt.Printf("in h2h: we win so attempt\n")
		// pick the move that results in h2h collision
		// TODO: if there are two possible squares for a collision and there is food in one of them, go for the one w/ food
		for i := 0; i < len(movesUs); i++ {
			futureUs := movesUs[i]
			for j := 0; j < len(movesEnemy); j++ {
				futureEnemy := movesEnemy[j]
				if (futureUs.X == futureEnemy.X && futureUs.Y == futureEnemy.Y) && IsMovePossible(us, board, IndexToMove(i)) {
					return IndexToMove(i)
				}
			} // end for j
		} // end for i
	}
	return "null"
}
