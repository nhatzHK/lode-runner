package game

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
)

type tile = byte
type level struct {
	num     int
	tiles   [][]tile
	players map[position]tile
	gold    []position
	escape  []position
}

// Tiles
const (
	EMPTY        = ' '
	RUNNER       = '&'
	GUARD        = '0'
	BRICK        = '#'
	SOLIDBRICK   = '@'
	FALSEBRICK   = 'X'
	LADDER       = 'H'
	ESCAPELADDER = 'S'
	ROPE         = '-'
	GOLD         = '$'
)

// Position
type position struct{ x, y int }

func manhattanDist(a, b position) float64 {
	return math.Abs(float64(a.x-b.x)) + math.Abs(float64(a.y-b.y))
}

func newLevel(num int) (*level, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf("levels/%03d.lvl", num))
	if err != nil {
		return nil, err
	}

	// TODO Specify right len/cap
	lvl := &level{
		num, bytes.Split(content, []byte("\n")),
		make(map[position]tile), make([]position, 0), make([]position, 0),
	}

	// Collect data for players/gold
	for i, tiles := range lvl.tiles {
		for j, tile := range tiles {
			switch tile {
			case RUNNER, GUARD:
				lvl.players[position{j, i}] = tile
				lvl.tiles[i][j] = EMPTY
			case GOLD:
				lvl.gold = append(lvl.gold, position{j, i})
				lvl.tiles[i][j] = EMPTY
			case ESCAPELADDER:
				lvl.escape = append(lvl.escape, position{j, i})
				lvl.tiles[i][j] = EMPTY
			}
		}
	}

	return lvl, nil
}

func (l *level) String() string {
	return string(bytes.Join(l.getTiles(), []byte("\n")))
}

func (l *level) emptyBelow(pos position) bool {
	return l.getTiles()[pos.y+1][pos.x] == EMPTY
}

func (l *level) goldCollected() bool {
	return len(l.gold) == 0
}

// TODO Rewrite + Rename
func (l *level) getTiles() [][]tile {
	tiles := make([][]tile, len(l.tiles))
	for i := range tiles {
		tiles[i] = make([]tile, len(l.tiles[i]))
		copy(tiles[i], l.tiles[i])
	}

	// Gold
	for _, pos := range l.gold {
		tiles[pos.y][pos.x] = GOLD
	}

	// Escape ladders
	if !l.goldCollected() {
		for _, pos := range l.escape {
			tiles[pos.y][pos.x] = ESCAPELADDER
		}
	}

	// Players
	for pos, tile := range l.players {
		tiles[pos.y][pos.x] = tile
	}

	return tiles
}

func (l *level) validMove(orig, dest position, dir direction) bool {
	if dest.x < 0 || dest.x >= 28 /*|| dest.y < 0*/ || dest.y >= 16 {
		return false
	}

	// FIXME FIXME FIXME FIXME FIXME FIXME FIXME
	origTile := l.tiles[orig.y][orig.x]

	if !l.goldCollected() && origTile == ESCAPELADDER {
		origTile = EMPTY
	}

	// if dir == DOWN && origTile == ROPE {
	// 	return false
	// }

	if dest.y < 0 {
		return origTile == ESCAPELADDER || origTile == LADDER
	}

	destTile := l.tiles[dest.y][dest.x]

	if !l.goldCollected() && destTile == ESCAPELADDER {
		destTile = EMPTY
	}

	switch destTile {
	case EMPTY, ROPE:
		if dir == UP {
			return origTile == LADDER || origTile == ESCAPELADDER
		}
		return true
	/*
		case BRICK:
			return dir != UP && false // TODO && bricksbrokenat(dest)
	*/
	case BRICK, SOLIDBRICK:
		return false
	case LADDER, ESCAPELADDER:
		return true
		// case LADDER:
		// 	return dir != DOWN
	}

	return false
}

func (l *level) validDig(pos position) bool {
	if pos.x < 0 || pos.x >= 28 {
		return false
	}

	return l.tiles[pos.y][pos.x] == BRICK
}
