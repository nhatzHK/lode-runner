package game

import (
	"encoding/json"
	"flag"
	//"fmt"
	"log"
	"time"

	msg "github.com/abdelq/lode-runner/message"
)

var tick = flag.String("tick", "250ms", "duration of game tick")

type Game struct {
	level     *level
	runner    *Runner
	guards    map[*Guard]struct{}
	ticker    *time.Ticker
	broadcast chan *msg.Message
	tick_time time.Duration
}

func NewGame(broadcast chan *msg.Message) *Game {
	dur, err := time.ParseDuration(*tick)
	if err != nil {
		dur = 250 * time.Millisecond // XXX
	}

	game := &Game{
		guards:    make(map[*Guard]struct{}),
		ticker:    time.NewTicker(dur),
		broadcast: broadcast,
		tick_time: dur,
	}

	return game
}

func (g *Game) Started() bool {
	return g.level != nil // XXX
}

func (g *Game) filled() bool {
	return g.runner != nil && len(g.guards) == 0 // XXX
}

// FIXME
func (g *Game) start(lvl int) {
	g.ticker.Stop()

	level, err := newLevel(lvl)
	if err != nil {
		log.Println(err)
		// TODO Broadcast error && Stop game
		return
	}

	/* Runner */
	g.runner.init(level.players)

	/* Guards */
	for guard := range g.guards {
		guard.init(level.players)
	}

	// Delete rest
PLAYERS:
	for pos, tile := range level.players {
		if tile == GUARD {
			for guard := range g.guards {
				if guard.pos == pos {
					continue PLAYERS
				}
			}
			delete(level.players, pos)
		}
	}

	g.broadcast <- msg.NewMessage("start", level.String())
	g.level = level // XXX

	time.Sleep(50 * time.Millisecond)

	g.ticker = time.NewTicker(250 * time.Millisecond)
	go g.gameTick()
}

// FIXME
func (g *Game) stop(winner tile) {
	g.ticker.Stop() // XXX Verify garbage collection

	if winner == RUNNER {
		g.broadcast <- msg.NewMessage("quit", "runner wins")
	} else {
		g.broadcast <- msg.NewMessage("quit", "runner looses")
	}

	// switch winner {
	// case RUNNER:
	// 	g.broadcast <- msg.NewMessage("quit", "runner wins") // TODO
	// case GUARD:
	// 	g.broadcast <- msg.NewMessage("quit", "runner looses") // TODO
	// default:
	// 	g.broadcast <- msg.NewMessage("quit", "draw") // TODO
	// }
}

func (g *Game) hasPlayer(name string) bool {
	// Runner
	if g.runner != nil && g.runner.name == name {
		return true
	}

	// Guards
	for guard := range g.guards {
		if guard.name == name {
			return true
		}
	}

	return false
}

func (g *Game) gameTick() {
	for range g.ticker.C {
		//if g.Started() {
		// Runner
		switch runner := g.runner; runner.action.Type {
		case MOVE:
			runner.move(runner.action.Direction, g)
			runner.action = action{}
		case DIG:
			runner.dig(runner.action.Direction, g)
			runner.action = action{}
		}

		// Guards
		for guard := range g.guards {
			guard.move(guard.action.Direction, g)
			guard.action = action{}
		}

		// XXX
		next := struct {
			Runner struct {
				Position struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"position"`
			} `json:"runner"`
		}{}
		next.Runner.Position.X = g.runner.pos.x
		next.Runner.Position.Y = g.runner.pos.y
		//next.Runner.Position = position{g.runner.pos.x, g.runner.pos.y}

		stuff, _ := json.Marshal(next)

		//fmt.Println(next)

		g.runner.out <- &msg.Message{"next", stuff}

		//g.runner.out <- &msg.Message{"next", []byte(`{"runner": {"position": {"x": ` + string(g.runner.pos.x) + `, "y": ` + string(g.runner.pos.y) + `}}}`)} // FIXME
		g.broadcast <- msg.NewMessage("next", g.level.String())
	}
	//}
}
