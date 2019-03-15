package session

import (
	"fmt"

	"github.com/dustin/randbo"
)

const (
	// SessionTime session有效期(单位:分钟), 默认7天
	SessionTimeout int = 7 * 24 * 60
	sessionIDLen       = 36
)

var Store SessionStorage

// SessionStore
type SessionStorage interface {
	// SetSession set
	SetSession(sessionID string, key string, data []byte) error
	// GetSession get
	GetSession(sessionID string, key string) ([]byte, bool, error)
	// ClearSession clear
	ClearSession(sessionID string, key string) error
}

type Session struct {
	ID    string
	name  string
	store SessionStorage
}

func NewSession(id, name string, store SessionStorage) *Session {
	return &Session{
		ID:    id,
		name:  name,
		store: store,
	}
}

func (s *Session) Get(key string) ([]byte, bool, error) {
	return s.store.GetSession(s.name+":"+s.ID, key)
}

func (s *Session) Set(key string, data []byte) error {
	return s.store.SetSession(s.name+":"+s.ID, key, data)
}

func (s *Session) Clear(key string) error {
	return s.store.ClearSession(s.name+":"+s.ID, key)
}

func NewSessionID() string {
	buf := make([]byte, sessionIDLen)
	randbo.New().Read(buf)
	return fmt.Sprintf("%x", buf)
}
