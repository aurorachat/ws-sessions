package sessions

type subscribedSessionData struct {
	listeningChannel chan interface{}
	sessionObject    *Session
}

type Channel struct {
	id         string
	subscribed map[string]subscribedSessionData
}

func NewChannel(id string) *Channel {
	return &Channel{id: id}
}

func (channel *Channel) Id() string {
	return channel.id
}

func (channel *Channel) Subscribe(session *Session, ch chan interface{}) {
	channel.subscribed[session.Id()] = subscribedSessionData{
		listeningChannel: ch,
		sessionObject:    session,
	}
}

func (channel *Channel) Unsubscribe(name string) {
	delete(channel.subscribed, name)
}

func (channel *Channel) UnsubscribeAll() {
	for name, sessionListening := range channel.subscribed {
		sessionListening.sessionObject.Unsubscribe(channel)
		delete(channel.subscribed, name)
	}
}

func (channel *Channel) Broadcast(msg interface{}) {
	for name := range channel.subscribed {
		channel.subscribed[name].listeningChannel <- msg
	}
}
