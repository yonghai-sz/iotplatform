package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"iotplatform/services/tcp-server/internal/protocol/model"
)

var ErrSessionClosed = errors.New("session closed")

const monitorSessionCntInterval = 10 * time.Second

type SessionCache struct {
	cacheByID map[string]*model.Session
	mu        sync.Mutex
}

var sessionCacheSingleton *SessionCache
var sessionCacheInitOnce sync.Once

func getSessionCache() *SessionCache {

	sessionCacheInitOnce.Do(func() {
		sessionCacheSingleton = &SessionCache{
			cacheByID: make(map[string]*model.Session),
		}

		monitorSessionCnt()
	})

	return sessionCacheSingleton
}

func monitorSessionCnt() {
	go func() {
		for {
			fmt.Printf("total_conn_cnt: %d", countSession())
			time.Sleep(monitorSessionCntInterval)
		}
	}()
}

func GetSession(id string) (*model.Session, error) {

	c := getSessionCache()

	c.mu.Lock()
	defer c.mu.Unlock()

	if s, ok := c.cacheByID[id]; ok {
		return s, nil
	}

	return nil, ErrSessionClosed
}

func StoreSession(s *model.Session) {
	c := getSessionCache()

	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheByID[s.ID] = s
}

func ClearSession(id string) {
	c := getSessionCache()
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cacheByID, id)
}

func countSession() int {
	c := getSessionCache()
	return len(c.cacheByID)
}
