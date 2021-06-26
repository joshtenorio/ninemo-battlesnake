package tilescore

import (
	"github.com/joshtenorio/ninemo-bot/datatypes"
	//"github.com/joshtenorio/ninemo-bot/snake/floodfill"
	//"github.com/joshtenorio/ninemo-bot/snake/api"
)

func CalculateTileScore(board *datatypes.Board, tile datatypes.Coord, ffMaxSearch int, ffThreshold int) int {
	// if tile is in hazard, give it minus
	// if tile has food, give it positive
	// use flood fill to calculate score, if number of free tiles is less than threshold it is negative score, else positive score
	// if h2h collision possible and we win, positive score
	// score for a winning h2h should be smaller in magnitude than a negative flood fill score
	//        because if the other snake doesn't go into h2h and we get locked into a space, that is pretty bad

	return 0 // placeholder
}
