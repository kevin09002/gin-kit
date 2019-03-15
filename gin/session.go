package ginkits

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinma2010/gkits/session"
)

const (
	// SessionKey cookie name
	SessionKey string = "ZQSESSID"
	// SessionTime session有效期(单位:分钟), 默认7天
	SessionTimeout      int = 7 * 24 * 60
	SessionCookieDomain     = ""
	sessionIDLen            = 36
	DefaultKey              = "kevinma2010/gkits/gin/session"
)

func Default(c *gin.Context) *session.Session {
	return c.MustGet(DefaultKey).(*session.Session)
}

func SessionHandler(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := session.NewSession(GetSessionID(c), name, session.Store)
		c.Set(DefaultKey, s)
		c.Next()
	}
}

func GetSessionID(c *gin.Context) string {
	cookieValue, _ := c.Cookie(SessionKey)
	if cookieValue == "" {
		cookieValue = session.NewSessionID()
		c.SetCookie(SessionKey, cookieValue, SessionTimeout*60, "/", SessionCookieDomain, false, false)
	}
	return cookieValue
}
