package sessions

import (
	"context"
	"github.com/gorilla/websocket"
	"slices"
)

type ListeningConnectionData struct {
	connId string
	conn   *websocket.Conn
}

type Session struct {
	id                 string
	connections        map[string]*ListeningConnectionData
	listeningTo        map[string]*Channel
	listeningChannel   chan interface{}
	connectionsChannel chan sessionReceivedData
	ctx                context.Context
	cancelFunc         context.CancelFunc
}

type sessionReceivedData struct {
	sessionId    string
	receivedData interface{}
}

func NewSession(id string) *Session {
	ctx, cancelCtx := context.WithCancel(context.TODO())
	inst := Session{id: id, ctx: ctx, cancelFunc: cancelCtx, connections: make(map[string]*ListeningConnectionData), listeningTo: make(map[string]*Channel), listeningChannel: make(chan interface{}), connectionsChannel: make(chan sessionReceivedData)}
	go inst.startListeningToChannels()
	return &inst
}

func (session *Session) RegisterConnection(connectionId string, conn *websocket.Conn) {
	data := ListeningConnectionData{connectionId, conn}
	session.connections[connectionId] = &data
	go session.startListeningToWebsocket(&data)
}

func (session *Session) CloseSpecificConnections(specificConnections ...string) {
	for _, connData := range session.connections {
		if slices.Contains(specificConnections, connData.connId) {
			_ = connData.conn.Close()
			delete(session.connections, connData.connId)
		}
	}
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
	for _, data := range session.connections {
		_ = data.conn.Close()
	}
	session.cancelFunc()
}

func (session *Session) Receive() (string, interface{}) {
	data := <-session.connectionsChannel
	return data.sessionId, data.receivedData
}

func (session *Session) Send(msg interface{}, specificConnections ...string) {
	for _, connData := range session.connections {
		if len(specificConnections) != 0 {
			if !slices.Contains(specificConnections, connData.connId) {
				continue
			}
		}

		_ = connData.conn.WriteJSON(msg)
	}
}

func (session *Session) startListeningToChannels() {
	for {
		select {
		case <-session.ctx.Done():
			return
		case data := <-session.listeningChannel:
			session.Send(data)
		}
	}
}

func (session *Session) startListeningToWebsocket(data *ListeningConnectionData) {
	for {
		var msg interface{}
		err := data.conn.ReadJSON(&msg)
		if err != nil {
			delete(session.connections, data.connId)
		}
		session.connectionsChannel <- sessionReceivedData{
			data.connId, msg,
		}
	}
}
