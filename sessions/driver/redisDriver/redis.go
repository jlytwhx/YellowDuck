package redisDriver

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	gsessions "github.com/gorilla/sessions"
	"github.com/jlytwhx/YellowDuck/sessions"
)

type Store interface {
	sessions.Store
}

func NewStore(size int, network, address, password string) (Store, error) {
	s, err := NewRediStore(size, network, address, password)
	if err != nil {
		return nil, err
	}
	return &store{s}, nil
}

// NewStoreWithDB - like NewStore but accepts `DB` parameter to select
// redis DB instead of using the default one ("0")
//
// Ref: https://godoc.org/github.com/boj/redistore#NewRediStoreWithDB
func NewStoreWithDB(size int, network, address, password, DB string) (Store, error) {
	s, err := NewRediStoreWithDB(size, network, address, password, DB)
	if err != nil {
		return nil, err
	}
	return &store{s}, nil
}

// NewStoreWithPool instantiates a RediStore with a *redis.Pool passed in.
//
// Ref: https://godoc.org/github.com/boj/redistore#NewRediStoreWithPool
func NewStoreWithPool(pool *redis.Pool) (Store, error) {
	s, err := NewRediStoreWithPool(pool)
	if err != nil {
		return nil, err
	}
	return &store{s}, nil
}

type store struct {
	*RediStore
}

// GetRedisStore get the actual woking store.
//
// Ref: https://godoc.org/github.com/boj/redistore#RediStore
func GetRedisStore(s Store) (err error, rediStore *RediStore) {
	realStore, ok := s.(*store)
	if !ok {
		err = errors.New("unable to get the redis store: Store isn't *store")
		return
	}

	rediStore = realStore.RediStore
	return
}

// SetKeyPrefix sets the key prefix in the redis database.
func SetKeyPrefix(s Store, prefix string) error {
	err, rediStore := GetRedisStore(s)
	if err != nil {
		return err
	}

	rediStore.SetKeyPrefix(prefix)
	return nil
}

func (c *store) Options(options sessions.Options) {
	c.RediStore.Options = &gsessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
