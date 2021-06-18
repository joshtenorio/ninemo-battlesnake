package floodfill

import (
	"github.com/joshtenorio/ninemo-bot/datatypes"
	"github.com/joshtenorio/ninemo-bot/snake/api"
)

/*
Uses flood fill to count the number of free spaces around pos within maxDist
maxDist is essentially dx+dy, or the number of moves it takes to reach a coord from pos
uses slices as a way to implement queues
references:
https://en.wikipedia.org/wiki/Flood_fill#Moving_the_recursion_into_a_data_structure

TODO: not very efficient so look for other ways to implement floodfill
*/
func CountFreeSpaces(board *datatypes.Board, pos datatypes.Coord, maxDist int) int {
	count := 0
	queue := make([]datatypes.Coord, 0)
	queue = append(queue, pos)
	visited := make([]datatypes.Coord, 0)
	for len(queue) != 0 {
		n := queue[0]     // set n to first coord in queue
		queue = queue[1:] // pop queue
		dist := (n.X - pos.X) + (n.Y - pos.Y)
		if dist < 0 { // if dist is negative, make it positive
			dist *= -1
		}

		if !api.IsBlocking(board, n) && dist < maxDist && !LinearSearch(visited, n) {
			count++
			queue = append(queue, api.MoveToCoord("up", &n))
			queue = append(queue, api.MoveToCoord("down", &n))
			queue = append(queue, api.MoveToCoord("left", &n))
			queue = append(queue, api.MoveToCoord("right", &n))

			visited = append(visited, n)
		}
	} // end for len(queue) != 0
	return count
}

// helper function for CountFreeSpaces
func LinearSearch(list []datatypes.Coord, key datatypes.Coord) bool {
	for i := 0; i < len(list); i++ {
		if list[i].X == key.X && list[i].Y == key.Y {
			return true
		}
	}
	return false
}
