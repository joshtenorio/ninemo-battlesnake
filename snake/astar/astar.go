package astar

import (
	"github.com/joshtenorio/ninemo-battlesnake/datatypes"
)

func AStar(start datatypes.Coord, goal datatypes.Coord, board datatypes.Board) []datatypes.Coord {
	openSet := []datatypes.Coord{start}
	cameFrom := make(map[datatypes.Coord]datatypes.Coord)

	gScore := InitializeScoreMap(board.Height, board.Width)
	gScore[start] = 0

	fScore := InitializeScoreMap(board.Height, board.Width)
	fScore[start] = Heuristic(start, goal)

	for len(openSet) != 0 {
		current := SelectCoord(openSet, fScore)
		if current == goal {
			return ReconstructPath(cameFrom, current)
		}
		// TODO: remove current from openset
		neighbors := GenerateNeighbors(current)
		for i := 0; i < len(neighbors); i++ {
			tentativeGScore := 1 + gScore[current]
			if tentativeGScore < gScore[neighbors[i]] {
				cameFrom[neighbors[i]] = current
				gScore[neighbors[i]] = tentativeGScore
				fScore[neighbors[i]] = tentativeGScore + Heuristic(neighbors[i], goal)
			}
			// TODO: if neighbors[i] is not in openSet, add it to openSet
		}

	}
	// if this is reached, openSet is empty but goal was never reached
	failure := []datatypes.Coord{{X: 0, Y: 0}}
	return failure
}

func ReconstructPath(cameFrom map[datatypes.Coord]datatypes.Coord, current datatypes.Coord) []datatypes.Coord {
	totalPath := []datatypes.Coord{current}
	return totalPath
}

func Heuristic(start datatypes.Coord, goal datatypes.Coord) int {
	// first, if start is out of bounds or unpathable, give it an ultra high sore
	return 0
}

func InitializeScoreMap(height int, width int) map[datatypes.Coord]int {
	scoreMap := make(map[datatypes.Coord]int)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			scoreMap[datatypes.Coord{X: i, Y: j}] = 999
		}
	}
	return scoreMap
}

func GenerateNeighbors(coord datatypes.Coord) []datatypes.Coord {
	var neighbors []datatypes.Coord
	neighbors = append(neighbors, datatypes.Coord{X: coord.X, Y: coord.Y + 1})
	neighbors = append(neighbors, datatypes.Coord{X: coord.X, Y: coord.Y - 1})
	neighbors = append(neighbors, datatypes.Coord{X: coord.X + 1, Y: coord.Y})
	neighbors = append(neighbors, datatypes.Coord{X: coord.X - 1, Y: coord.Y})
	return neighbors
}

/* gets the cheapest coord in openset by fScore */
func SelectCoord(openSet []datatypes.Coord, fScore map[datatypes.Coord]int) datatypes.Coord {
	cheapest := openSet[0]
	for i := 0; i < len(openSet); i++ {
		if fScore[openSet[i]] < fScore[cheapest] {
			cheapest = openSet[i]
		}
	}
	return cheapest
}

func GetPathLength(start datatypes.Coord, goal datatypes.Coord, board datatypes.Board) int {
	path := AStar(start, goal, board)
	return len(path)
}
