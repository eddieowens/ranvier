package ws

import (
	"fmt"
	"github.com/eddieowens/axon"
	"github.com/eddieowens/ranvier/server/app/pubsub"
	"github.com/gorilla/websocket"
	"net/http"
)

const WebsocketKey = "Websocket"

type Websocket interface {
	Connect(topic string, wr http.ResponseWriter, req *http.Request, h http.Header) error
}

type websocketImpl struct {
	PubSub   pubsub.PubSub `inject:"PubSub"`
	upgrader websocket.Upgrader
}

func (w *websocketImpl) Connect(topic string, wr http.ResponseWriter, req *http.Request, h http.Header) error {
	ws, err := w.upgrader.Upgrade(wr, req, h)
	if err != nil {
		return err
	}

	defer ws.Close()

	for c := range w.PubSub.Subscribe(topic) {
		err = ws.WriteJSON(c)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return nil
}

func websocketFactory(_ axon.Injector, _ axon.Args) axon.Instance {
	w := &websocketImpl{}

	w.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return axon.StructPtr(w)
}
