package utils

import (
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"time"
)

type CacheType int

const (
	NLG_KEYWORD CacheType = iota
	NLG_TEM
	NLG_IT
	NLG_LUA_VM
	NLG_LUA_FUNC
)

var Cache *cache.Cache

func InitCache() {
	c := cache.New(5*time.Minute, 10*time.Minute)

	if c == nil {
		log.Panicf("init cache error")
	}

	Cache = c
}

func Set(key string, entry interface{}, expireTime time.Duration) {
	Cache.Set(key, entry, expireTime)
}

func Get(key string) (interface{}, bool) {
	return Cache.Get(key)
}

func Del(key string)  {
	Cache.Delete(key)
}

func FormatCacheKey(original string, cacheType CacheType) string {

	switch cacheType {
	case NLG_KEYWORD:
		return original + "_KY"
	case NLG_TEM:
		return original + "_TEM"
	case NLG_IT:
		return original + "_IT"
	case NLG_LUA_VM:
		return original + "_LUA_VM"
	case NLG_LUA_FUNC:
		return original + "_LUA_FUNC"
	default:
		return original
	}
}
