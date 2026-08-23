package storage

import (
	"strings"
	"sync"

	"iot-zero/services/tcp-server/internal/protocol/model"
)

type SegmentCache struct {
	cacheByKey map[string]*model.Segment
	mu         sync.Mutex
}

var segmentCacheSingleton *SegmentCache
var segmentCacheInitOnce sync.Once

func getSegmentCache() *SegmentCache {
	segmentCacheInitOnce.Do(func() {
		segmentCacheSingleton = &SegmentCache{
			cacheByKey: make(map[string]*model.Segment),
		}
	})
	return segmentCacheSingleton
}

func CacheSegment(seg *model.Segment) bool {
	b := new(strings.Builder)
	b.WriteString(seg.Phone)
	b.WriteString("/")
	b.WriteRune(rune(seg.MsgID))
	key := b.String()

	cache := getSegmentCache()
	cache.mu.Lock()
	defer cache.mu.Unlock()

	s, ok := cache.cacheByKey[key]
	if !ok {
		cache.cacheByKey[key] = seg
		return s.IsComplete()
	}

	s.Merge(seg)

	// cache.cacheByKey[key] = s

	return s.IsComplete()
}
