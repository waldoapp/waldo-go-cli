package lib

import (
	"net/url"

	"github.com/gorilla/websocket"
)

func openWebSocket() error {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/echo"}

	_, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	return err
}
