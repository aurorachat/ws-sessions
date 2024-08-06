package sessions

type Store struct {
	channels map[string]*Channel
	sessions map[string]*Session
}

func NewStore() Store {
	return Store{}
}

func (store *Store) GetSession(id string) *Session {
	return store.sessions[id]
}

func (store *Store) GetChannel(id string) *Channel {
	return store.channels[id]
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
