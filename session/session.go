package session

import (
	"fmt"

	"github.com/dustin/randbo"
	"github.com/gin-gonic/gin"
)

const (
	// SessionKey cookie name
	SessionKey string = "ZQSESSID"
	// SessionTime session有效期(单位:分钟), 默认7天
	SessionTimeout      int = 7 * 24 * 60
	SessionCookieDomain     = ""
	sessionIDLen            = 36
	DefaultKey              = "kevin09002/gin-kit/session"
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

func Default(c *gin.Context) *Session {
	return c.MustGet(DefaultKey).(*Session)
}

func SessionHandler(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := &Session{
			ID:    GetSessionID(c),
			name:  name,
			store: Store,
		}
		c.Set(DefaultKey, s)
		c.Next()
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

func GetSessionID(c *gin.Context) string {
	cookieValue, _ := c.Cookie(SessionKey)
	if cookieValue == "" {
		cookieValue = newSessionID()
		c.SetCookie(SessionKey, cookieValue, SessionTimeout*60, "/", SessionCookieDomain, false, false)
	}
	return cookieValue
}

func newSessionID() string {
	buf := make([]byte, sessionIDLen)
	randbo.New().Read(buf)
	return fmt.Sprintf("%x", buf)
}
