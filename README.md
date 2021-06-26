# ninemo
[Battlesnake](https://play.battlesnake.com/) written using Go.

[Battlesnake server url (heroku)](https://ninemo-bot.herokuapp.com/)

[Battlesnake page](https://play.battlesnake.com/u/tenmo/ninemo/)

- the live version is on branch `main`
- in-progress code is on branch `dev`

## Strategy
1. Calculate a score for each tile I can move into - favorable tiles will have a higher score, unfavorable tiles will have a lower score
2. Move into the tile whose score is highest

## Resources I used

- [Flood fill](https://en.wikipedia.org/wiki/Flood_fill#Moving_the_recursion_into_a_data_structure)
- [Mojave - local testing](https://github.com/smallsco/mojave)
- [battlesnake board generator - local testing in specific conditions](https://github.com/Nettogrof/battle-snake-board-generator)