package botgoram

// StateInitializer returns an initialized state data.
type StateInitializer func(uid string) interface{}

// SaveLoader defines how state data are persisted into storage.
// It should return an initialized state data when data not found.
// Returns error only when something goes wrong, eg: connection lost when saving data to db.
type SaveLoader interface {
	Save(uid string, sid string, data interface{}) error
	Load(uid string) (sid string, data interface{}, err error)
}

type memoryStore struct {
	data  map[string]interface{}
	state map[string]string
	init  StateInitializer
}

// MemoryStore provides default, memory based SaveLoader implementation.
func MemoryStore(init StateInitializer) SaveLoader {
	return &memoryStore{
		make(map[string]interface{}),
		make(map[string]string),
		init,
	}
}

func (m *memoryStore) Save(uid string, sid string, data interface{}) error {
	m.data[uid] = data
	m.state[uid] = sid
	return nil
}

func (m *memoryStore) Load(uid string) (sid string, data interface{}, err error) {
	data, ok := m.data[uid]
	sid = m.state[uid]
	if !ok {
		data = m.init(uid)
	}
	return
}
