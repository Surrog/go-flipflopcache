package goflipflopcache

import "time"

type FlipFlopCache struct {
	cacheA map[string][]byte
	cacheB map[string][]byte

	mainCacheA       bool
	timeBetweenReset time.Duration
	timeLastFlip     time.Time
}

func NewFlipFlopCache(timeBetweenReset time.Duration) *FlipFlopCache {
	return &FlipFlopCache{
		cacheA:           map[string][]byte{},
		cacheB:           map[string][]byte{},
		mainCacheA:       true,
		timeBetweenReset: timeBetweenReset,
		timeLastFlip:     time.Now(),
	}
}

func (cache *FlipFlopCache) Flip() {
	if cache.mainCacheA {
		cache.cacheB = map[string][]byte{}
	} else {
		cache.cacheA = map[string][]byte{}
	}
	cache.mainCacheA = !cache.mainCacheA
	cache.timeLastFlip = time.Now()
}

func (cache *FlipFlopCache) Reset() {
	cache.mainCacheA = true
	cache.cacheA = map[string][]byte{}
	cache.cacheB = map[string][]byte{}
	cache.timeLastFlip = time.Now()
}

func get(mainCache *map[string][]byte, secondCache *map[string][]byte, key string) ([]byte, bool) {
	val, ok := (*mainCache)[key]
	if ok {
		return val, ok
	}

	oldVal, oldOk := (*secondCache)[key]
	if oldOk {
		(*mainCache)[key] = oldVal
	}
	return oldVal, oldOk
}

func (cache *FlipFlopCache) expireFlip() {
	if time.Since(cache.timeLastFlip) > 2*cache.timeBetweenReset {
		cache.Reset()
		return
	}

	if time.Since(cache.timeLastFlip) > cache.timeBetweenReset {
		cache.Flip()
	}
}

func (cache *FlipFlopCache) Get(key string) ([]byte, bool) {
	cache.expireFlip()

	if cache.mainCacheA {
		return get(&cache.cacheA, &cache.cacheB, key)
	} else {
		return get(&cache.cacheB, &cache.cacheA, key)
	}
}

func (cache *FlipFlopCache) Append(key string, value []byte) {
	cache.expireFlip()

	if cache.mainCacheA {
		cache.cacheA[key] = value
	} else {
		cache.cacheB[key] = value
	}
}
