package snake

import (
	"fmt"

	"github.com/joshtenorio/ninemo-bot/datatypes"
	"github.com/joshtenorio/ninemo-bot/snake/api"
	"github.com/joshtenorio/ninemo-bot/snake/floodfill"
)

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

	return !api.IsBlocking(board, position)
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
			if distSquared == 4 && api.IsBlocking(board, datatypes.Coord{X: (us.X + heads[i].X) / 2, Y: (us.Y + heads[i].Y) / 2}) {
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
