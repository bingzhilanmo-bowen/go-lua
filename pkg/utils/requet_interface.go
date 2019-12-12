package utils

import lua "github.com/yuin/gopher-lua"

type RequestInterface interface {
	 Struct2luaTable(l *lua.LState) *lua.LTable
}
