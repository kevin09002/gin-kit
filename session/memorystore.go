package session

type memoryStore struct {
	data map[string][]byte
}

func NewMemoryStore() *memoryStore {
	return &memoryStore{
		data: make(map[string][]byte),
	}
}

func (s *memoryStore) GetSession(sessionID, key string) ([]byte, bool, error) {
	data, exist := s.data[sessionID+":"+key]
	return data, exist, nil
}

func (s *memoryStore) SetSession(sessionID, key string, data []byte) error {
	s.data[sessionID+":"+key] = data
	return nil
}

func (s *memoryStore) ClearSession(sessionID, key string) error {
	delete(s.data, sessionID+":"+key)
	return nil
}
