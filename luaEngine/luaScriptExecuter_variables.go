package luaScriptEngine

import lua "github.com/yuin/gopher-lua"

// Holds all lua script file that is used
var luaScriptFilesAsByteArray [][]byte

// The shared Lua state used for execution
var luaState *lua.LState
