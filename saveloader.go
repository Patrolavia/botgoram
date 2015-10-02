package botgoram

// StateInitializer returns an initialized state data.
type StateInitializer func(uid int) interface{}

// SaveLoader defines how state data are persisted into storage.
// It should return an initialized state data when data not found.
// Returns error only when something goes wrong, eg: connection lost when saving data to db.
type SaveLoader interface {
	Save(uid int, sid StateId, data interface{}) error
	Load(uid int) (sid StateId, data interface{}, err error)
}

type memoryStore struct {
	data map[int]interface{}
	state map[int]StateId
	init StateInitializer
}

// MemoryStore provides default, memory based SaveLoader implementation.
func MemoryStore(init StateInitializer) SaveLoader {
	return &memoryStore{
		make(map[int]interface{}),
		make(map[int]StateId),
		init,
	}
}

func (m *memoryStore) Save(uid int, sid StateId, data interface{}) error {
	m.data[uid] = data
	m.state[uid] = sid
	return nil
}

func (m *memoryStore) Load(uid int) (sid StateId, data interface{}, err error) {
	data, ok := m.data[uid]
	sid = m.state[uid]
	if !ok {
		data = m.init(uid)
	}
	return
}
