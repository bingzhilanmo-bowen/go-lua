package utils

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"strconv"
	"testing"
	"time"
)

func TestInitCache(t *testing.T) {
	InitCache()

	proto, _ := CompileFile("../../sc/test_table.lua")

	l := lua.NewState()

	defer l.Close()



	lfunc := l.NewFunctionFromProto(proto)
	l.Push(lfunc)
	l.PCall(0, lua.MultRet, nil)


	tbl := l.NewTable()

	st := Unix2ShortTime(time.Now().Unix())
	ti , _ :=strconv.ParseInt(st, 10, 64)

	tbl.RawSetString("gender", lua.LString("woman"))
	tbl.RawSetString("requesttime", lua.LNumber(ti))

	l.CallByParam(lua.P{
		Fn:      l.GetGlobal("run"),
		NRet:    1,
		Protect: true,
	}, tbl)
	ret := l.Get(-1)

	res, ok := ret.(lua.LString)
	if ok {
		fmt.Println(res)
	} else {
		fmt.Println("unexpected result")
	}



/*
	Set("bowen", lfunc, 60*time.Second)

	//bts, _ :=  Get("bowen")

	lfunCache, _ := Get("bowen")
	cacheFunction := lfunCache.(*lua.LFunction)

	if cacheFunction != nil {

		i := 0
		l.Push(cacheFunction)
		l.PCall(0, lua.MultRet, nil)

		for i < 10 {

			l.CallByParam(lua.P{
				Fn:      l.GetGlobal("run"),
				NRet:    1,
				Protect: true,
			}, lua.LNumber(i))

			// 获取返回结果
			ret := l.Get(-1)

			res, ok := ret.(lua.LString)
			if ok {
				fmt.Println(res)
			} else {
				fmt.Println("unexpected result")
			}

			i++
		}

	}*/

}


func TestRegr(t *testing.T) {

	strs :=	GetExprList("{{TEST_START}},{{TEST_END}}", "{{[a-zA-Z_]+}}")

	for _,va := range strs{
		fmt.Println(va)
	}
	fmt.Println(len(strs))
}

func TestTimeCv(t *testing.T)  {
	st := Unix2ShortTime(1136149445)
	ti , _ :=strconv.ParseInt(st, 10, 64)
	fmt.Println(ti)
}

