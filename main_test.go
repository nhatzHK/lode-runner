package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func sendMsg(t *testing.T, conn io.Writer, sent message) {
	t.Helper()

	if err := json.NewEncoder(conn).Encode(&sent); err != nil {
		t.Error(err)
	}
}

func receiveMsg(t *testing.T, conn io.Reader, expected message) {
	t.Helper()

	var received message
	if err := json.NewDecoder(conn).Decode(&received); err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("expected: %s, received: %s", expected, received)
	}
}
