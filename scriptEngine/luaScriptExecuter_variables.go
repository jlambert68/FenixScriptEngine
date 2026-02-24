package scriptEngine

import lua "github.com/yuin/gopher-lua"

// luaScriptFilesAsByteArray stores all Lua scripts currently loaded by the engine.
var luaScriptFilesAsByteArray []LuaScriptsStruct

// luaState is the shared gopher-lua VM used for placeholder execution.
var luaState *lua.LState
