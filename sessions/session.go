package sessions

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"time"
)

const (
	DefaultKey  = "github.com/jlytwhx/yellowDuck"
	errorFormat = "[sessions] ERROR! %s\n"
)

type Key struct {
	Data            interface{}
	ExpireTimeStamp int64
}

type Store interface {
	sessions.Store
	Options(Options)
}

type Options struct {
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

type Session interface {
	// Get returns the sessions value associated to the given key.
	Get(key interface{}) interface{}
	// Set sets the sessions value associated to the given key.
	Set(key interface{}, val interface{}, expire int)
	// Delete removes the sessions value associated to the given key.
	Delete(key interface{})
	// Clear deletes all values in the sessions.
	Clear()
	// Options sets configuration for a sessions.
	Options(Options)
	// Save saves all sessions used during the current request.
	Save() error
}

func SessionMiddleware(name string, store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := &session{name, c.Request, store, nil, false, c.Writer}
		c.Set(DefaultKey, s)
		defer context.Clear(c.Request)
		c.Next()
	}
}

func SessionsManyMiddleware(names []string, store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		newSession := make(map[string]Session, len(names))
		for _, name := range names {
			newSession[name] = &session{name, c.Request, store, nil, false, c.Writer}
		}
		c.Set(DefaultKey, newSession)
		defer context.Clear(c.Request)
		c.Next()
	}
}

type session struct {
	name    string
	request *http.Request
	store   Store
	session *sessions.Session
	written bool
	writer  http.ResponseWriter
}

func (s *session) Get(key interface{}) interface{} {
	if data := s.Session().Values[key]; data != nil {
		result := data.(Key)
		expireTime := result.ExpireTimeStamp
		now := time.Now().Unix()
		if expireTime > now {
			return result.Data
		}
	}
	return nil
}

func (s *session) Written() bool {
	return s.written
}

func Default(c *gin.Context) Session {
	return c.MustGet(DefaultKey).(Session)
}

func (s *session) Session() *sessions.Session {
	if s.session == nil {
		var err error
		s.session, err = s.store.Get(s.request, s.name)
		if err != nil {
			log.Printf(errorFormat, err)
		}
	}
	return s.session
}

func (s *session) Options(options Options) {
	s.Session().Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

func (s *session) Save() error {
	if s.Written() {
		e := s.Session().Save(s.request, s.writer)
		if e == nil {
			s.written = false
		}
		return e
	}
	return nil
}

func (s *session) Set(key interface{}, val interface{}, expire int) {
	if expire > s.Session().Options.MaxAge || expire <= 0 {
		expire = s.Session().Options.MaxAge
	}
	s.Session().Values[key] = Key{val, int64(expire) + time.Now().Unix()}
	s.written = true
}

func (s *session) Delete(key interface{}) {
	delete(s.Session().Values, key)
	s.written = true
}

func (s *session) Clear() {
	for key := range s.Session().Values {
		s.Delete(key)
	}
}

func init() {
	gob.Register(Key{})
}
