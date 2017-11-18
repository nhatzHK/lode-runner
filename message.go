package main

import (
	"encoding/json"
	"strings"

	"github.com/abdelq/lode-runner/game"
	msg "github.com/abdelq/lode-runner/message"
)

type message msg.Message

func newMessage(event, data string) message {
	return message(*msg.NewMessage(event, data))
}

// TODO Move sections to game package
func (m *message) parse(sender *client) {
	switch evt := strings.ToLower(strings.TrimSpace(m.Event)); evt {
	case "join":
		parseJoin(m.Data, sender)
	case "move":
		parseMove(m.Data, sender)
	case "dig":
		parseDig(m.Data, sender)
	default:
		sender.out <- newMessage("error", "invalid event")
	}
}

func parseJoin(data json.RawMessage, sender *client) {
	message := new(msg.JoinMessage)
	if err := message.Parse(data); err != nil {
		sender.out <- newMessage("error", err.Error())
		return
	}

	// Find/Create room
	room, ok := rooms[message.Room]
	if !ok {
		room = newRoom(message.Room)
	}

	room.join <- &join{sender, game.NewPlayer(message.Name, message.Role)}
}

// TODO Move to game package
func parseMove(data json.RawMessage, sender *client) {
	message := new(msg.GameMessage)
	if err := message.Parse(data); err != nil {
		sender.out <- newMessage("error", err.Error())
		return
	}

	// Find room name if none sent
	if message.Room == "" {
		if message.Room = findRoom(sender); message.Room == "" {
			sender.out <- newMessage("error", "not in a room")
			return
		}
	}

	if room, ok := rooms[message.Room]; ok {
		if !room.game.Started() {
			sender.out <- newMessage("error", "game not started")
			return
		}

		if player := room.clients[sender]; player != nil {
			player.UpdateAction(game.MOVE, message.Direction)
		} else {
			sender.out <- newMessage("error", "not a player")
		}
	}
}

// TODO Move to game package
func parseDig(data json.RawMessage, sender *client) {
	message := new(msg.GameMessage)
	if err := message.Parse(data); err != nil {
		sender.out <- newMessage("error", err.Error())
		return
	}

	// Find room name if none sent
	if message.Room == "" {
		if message.Room = findRoom(sender); message.Room == "" {
			sender.out <- newMessage("error", "not in a room")
			return
		}
	}

	if room, ok := rooms[message.Room]; ok {
		if !room.game.Started() {
			sender.out <- newMessage("error", "game not started")
			return
		}

		if runner, ok := room.clients[sender].(*game.Runner); ok {
			runner.UpdateAction(game.DIG, message.Direction)
		} else {
			sender.out <- newMessage("error", "not a runner")
		}
	}
}
