package sessions

type Store struct {
	channels map[string]*Channel
	sessions map[string]*Session
}

func NewStore() Store {
	return Store{}
}

func (store *Store) GetSession(id string) *Session {
	session, ok := store.sessions[id]
	if ok {
		return session
	}
	return nil
}

func (store *Store) GetChannel(id string) *Channel {
	channel, ok := store.channels[id]
	if ok {
		return channel
	}
	return nil
}

func (store *Store) SetSession(id string, session *Session) {
	store.sessions[id] = session
}

func (store *Store) DeleteSession(id string) {
	session := store.GetSession(id)

	if session != nil {
		session.Close()
		delete(store.sessions, id)
	}
}

func (store *Store) DeleteChannel(id string) {
	channel := store.GetChannel(id)
	if channel != nil {
		channel.UnsubscribeAll()
		delete(store.channels, id)
	}
}
