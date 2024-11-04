package luaEngine

import lua "github.com/yuin/gopher-lua"

// Holds all lua script file that is used
var luaScriptFilesAsByteArray []LuaScriptsStruct

// The shared Lua state used for execution
var luaState *lua.LState
