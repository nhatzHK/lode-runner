package message

import (
	"encoding/json"
	"errors"
	"strings"
)

type GameMessage struct {
	Direction uint8
	Room      string
}

func (m *GameMessage) Parse(data json.RawMessage) error {
	if err := json.Unmarshal(data, m); err != nil {
		return err
	}

	if m.Direction > 4 { // NONE is a valid direction
		return errors.New("invalid direction")
	}
	m.Room = strings.TrimSpace(m.Room)

	return nil
}
