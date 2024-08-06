package sessions

import (
	"context"
	"github.com/gorilla/websocket"
)

type Session struct {
	id               string
	conn             *websocket.Conn
	listeningTo      map[string]*Channel
	listeningChannel chan interface{}
	runningContext   context.Context
	cancelFunc       context.CancelFunc
}

func NewSession(id string, conn *websocket.Conn) *Session {
	ctx, cancelCtx := context.WithCancel(context.TODO())
	return &Session{id: id, conn: conn, runningContext: ctx, cancelFunc: cancelCtx}
}

func (session *Session) Id() string {
	return session.id
}

func (session *Session) Subscribe(channel *Channel) {
	session.listeningTo[session.Id()] = channel
	channel.Subscribe(session, session.listeningChannel)
}

func (session *Session) Unsubscribe(channel *Channel) {
	delete(session.listeningTo, channel.Id())
	channel.Unsubscribe(session.Id())
}

func (session *Session) Close() {
	_ = session.conn.Close()
	session.cancelFunc()
}

func (session *Session) StartListening() {
	for {
		select {
		case <-session.runningContext.Done():
			return
		case data := <-session.listeningChannel:
			err := session.conn.WriteJSON(data)
			if err != nil {
				session.Close()
			}
		}
	}
}
