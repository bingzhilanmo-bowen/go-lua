package utils

import (
	"bufio"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	"os"
	"strings"
	"sync"
)

var setFunCacheLock sync.Mutex
var setVmCacheLock sync.Mutex

var vmCacheKey = make(map[string]int8)

func DirectExecute(luaScript, callMethod string, args ...lua.LValue) (lua.LValue, error) {
	l := lua.NewState()
	defer l.Close()

	err := l.DoString(luaScript)

	if err != nil {
		log.Errorf("lua do string error %s", err.Error())
		return nil, err
	}

	if callMethod != "" {
		if err := l.CallByParam(lua.P{
			Fn:      l.GetGlobal(callMethod),
			NRet:    lua.MultRet,
			Protect: true,
		}, args...); err != nil {
			log.Errorf("execute lua method [%s] error %s", callMethod, err.Error())
			return nil, err
		}
	}

	result := l.Get(-1)
	l.Pop(1)

	return result, nil
}

// 目前默认都是一直缓存
func CacheCompileExecute(cacheKey, luaScript, callMethod string, args ...lua.LValue) (lua.LValue, error) {

	cacheKey =FormatCacheKey(cacheKey, NLG_LUA_FUNC)

	l := lua.NewState()
	defer l.Close()

	lfunCache, _ := Get(cacheKey)
	if lfunCache == nil {
		setFunCacheLock.Lock()
		defer setFunCacheLock.Unlock()
		// double check
		funCache, _ := Get(cacheKey)
		if funCache == nil {
			log.Infof("set cache key [%s]", cacheKey)
			proto, err := CompileString(luaScript)
			if err != nil {
				return nil, err
			}
			funCache := l.NewFunctionFromProto(proto)

			lfunCache = funCache
			//
			Set(cacheKey, funCache, cache.NoExpiration)
		} else {
			lfunCache = funCache
		}
	}

	cacheFunction := lfunCache.(*lua.LFunction)

	l.Push(cacheFunction)
	err := l.PCall(0, lua.MultRet, nil)

	if err != nil {
		log.Errorf("set cache function error %s", err.Error())
		return nil, err
	}

	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal(callMethod),
		NRet:    lua.MultRet,
		Protect: true,
	}, args...); err != nil {
		log.Errorf("execute lua method [%s] error %s", callMethod, err.Error())
		return nil, err
	}

	result := l.Get(-1)
	l.Pop(1)

	return result, nil
}

func CacheLuaVmExecute(cacheKey, luaScript, callMethod string, args ...lua.LValue) (lua.LValue, error) {

	cacheKey =FormatCacheKey(cacheKey, NLG_LUA_VM)

	lCacheVm, _ := Get(cacheKey)

	if lCacheVm == nil {
		setVmCacheLock.Lock()
		defer setVmCacheLock.Unlock()
		// double check
		vmCache, _ := Get(cacheKey)
		if vmCache == nil {
			log.Infof("set vm cache key [%s]", cacheKey)
			vmCache := lua.NewState()

			err := vmCache.DoString(luaScript)

			if err != nil {
				log.Errorf("lua vm do string error %s", err.Error())
				return nil, err
			}

			// Set
			vmCacheKey[cacheKey] = 1

			lCacheVm = vmCache
			//
			Set(cacheKey, vmCache, cache.NoExpiration)
		} else {
			lCacheVm = vmCache
		}
	}

	cacheVm := lCacheVm.(*lua.LState)

	if err := cacheVm.CallByParam(lua.P{
		Fn:      cacheVm.GetGlobal(callMethod),
		NRet:    lua.MultRet,
		Protect: true,
	}, args...); err != nil {
		log.Errorf("execute lua method [%s] error %s", callMethod, err.Error())
		return nil, err
	}

	result := cacheVm.Get(-1)
	cacheVm.Pop(1)

	return result, nil
}

func CacheLuaVmExecute2(cacheKey, luaScript, callMethod string, arg RequestInterface) (lua.LValue, error) {

	cacheKey =FormatCacheKey(cacheKey, NLG_LUA_VM)

	lCacheVm, _ := Get(cacheKey)

	if lCacheVm == nil {
		setVmCacheLock.Lock()
		defer setVmCacheLock.Unlock()
		// double check
		vmCache, _ := Get(cacheKey)
		if vmCache == nil {
			log.Infof("set vm cache key [%s]", cacheKey)
			vmCache := lua.NewState()

			err := vmCache.DoString(luaScript)

			if err != nil {
				log.Errorf("lua vm do string error %s", err.Error())
				return nil, err
			}

			// Set
			vmCacheKey[cacheKey] = 1

			lCacheVm = vmCache
			//
			Set(cacheKey, vmCache, cache.NoExpiration)
		} else {
			lCacheVm = vmCache
		}
	}

	cacheVm := lCacheVm.(*lua.LState)

	if err := cacheVm.CallByParam(lua.P{
		Fn:      cacheVm.GetGlobal(callMethod),
		NRet:    lua.MultRet,
		Protect: true,
	}, arg.Struct2luaTable(cacheVm)); err != nil {
		log.Errorf("execute lua method [%s] error %s", callMethod, err.Error())
		return nil, err
	}

	result := cacheVm.Get(-1)
	cacheVm.Pop(1)

	return result, nil
}


func CacheExecuteWithMap(cacheKey, luaScript, callMethod string, params map[string]string) (lua.LValue, error) {

	cacheKey =FormatCacheKey(cacheKey, NLG_LUA_VM)

	lCacheVm, _ := Get(cacheKey)

	if lCacheVm == nil {
		setVmCacheLock.Lock()
		defer setVmCacheLock.Unlock()
		// double check
		vmCache, _ := Get(cacheKey)
		if vmCache == nil {
			log.Infof("set vm cache key [%s]", cacheKey)
			vmCache := lua.NewState()

			err := vmCache.DoString(luaScript)

			if err != nil {
				log.Errorf("lua vm do string error %s", err.Error())
				return nil, err
			}

			// Set
			vmCacheKey[cacheKey] = 1

			lCacheVm = vmCache
			//
			Set(cacheKey, vmCache, cache.NoExpiration)
		} else {
			lCacheVm = vmCache
		}
	}

	cacheVm := lCacheVm.(*lua.LState)

	if err := cacheVm.CallByParam(lua.P{
		Fn:      cacheVm.GetGlobal(callMethod),
		NRet:    lua.MultRet,
		Protect: true,
	}, map2luaTable(params, cacheVm)); err != nil {
		log.Errorf("execute lua method [%s] error %s", callMethod, err.Error())
		return nil, err
	}

	result := cacheVm.Get(-1)
	cacheVm.Pop(1)

	return result, nil
}

func map2luaTable(params map[string]string, l *lua.LState) *lua.LTable  {
	tbl := l.NewTable()
	for k,v := range params {
		tbl.RawSetString(k, lua.LString(v))
	}
    return tbl
}

func WarmLuaVm(cacheKey, luaScript string) error {
	cacheKey = FormatCacheKey(cacheKey, NLG_LUA_VM)
	closeCache(cacheKey)

	log.Infof("set vm cache key [%s]", cacheKey)
	vmCache := lua.NewState()

	err := vmCache.DoString(luaScript)
	if err != nil {
		log.Errorf("lua vm do string error %s", err.Error())
		return err
	}
	vmCacheKey[cacheKey] = 1
	Set(cacheKey, vmCache, cache.NoExpiration)

	return nil
}

func closeCache(cacheKey string)  {
	lCacheVm, _ := Get(cacheKey)
	if lCacheVm != nil {
		cacheVm := lCacheVm.(*lua.LState)
		cacheVm.Close()
		Del(cacheKey)
	}
}


func CloseCacheVm() {
	for k, _ := range vmCacheKey {
		lCacheVm, _ := Get(k)
		if lCacheVm != nil {
			cacheVm := lCacheVm.(*lua.LState)
			cacheVm.Close()
		}
	}
}

// 编译 lua 代码字段
func CompileString(source string) (*lua.FunctionProto, error) {
	reader := strings.NewReader(source)
	chunk, err := parse.Parse(reader, source)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, source)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

// 编译 lua 代码文件
func CompileFile(filePath string) (*lua.FunctionProto, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, filePath)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, filePath)
	if err != nil {
		return nil, err
	}
	return proto, nil
}
