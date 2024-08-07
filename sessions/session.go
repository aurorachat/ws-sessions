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
	ctx              context.Context
	cancelFunc       context.CancelFunc
}

func NewSession(id string, conn *websocket.Conn) *Session {
	ctx, cancelCtx := context.WithCancel(context.TODO())
	return &Session{id: id, conn: conn, ctx: ctx, cancelFunc: cancelCtx}
}

func (session *Session) Id() string {
	return session.id
}

func (session *Session) Context() context.Context {
	return session.ctx
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

func (session *Session) Receive() (interface{}, error) {
	var msg interface{}
	err := session.conn.ReadJSON(&msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (session *Session) Send(msg interface{}) error {
	return session.conn.WriteJSON(msg)
}

func (session *Session) StartListening() {
	for {
		select {
		case <-session.ctx.Done():
			return
		case data := <-session.listeningChannel:
			err := session.Send(data)
			if err != nil {
				session.Close()
			}
		}
	}
}
