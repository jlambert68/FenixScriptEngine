package main

import (
	"github.com/jlambert68/FenixScriptEngine/luaEngine"
	"log"
)

func main() {

	var err error

	// No external Lua-Libraries
	var fenixLuaScripts []luaEngine.LuaScriptsStruct
	fenixLuaScripts = []luaEngine.LuaScriptsStruct{}

	// Load and initiate Lua-engine
	err = luaEngine.InitiateLuaScriptEngine(fenixLuaScripts)
	if err != nil {
		log.Fatalln("Error", err)
	}

}
