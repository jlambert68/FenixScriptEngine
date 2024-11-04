package luaEngine

import (
	_ "embed"
)

// Embed files

//go:embed src/date.lua
var date []byte

//go:embed src/Fenix_ControlledUniqueId.lua
var fenix_ControlledUniqueId []byte

//go:embed src/Fenix_RandomPositiveDecimalValue.lua
var fenix_RandomPositiveDecimalValue []byte

//go:embed src/Fenix_TodayDateShift.lua
var fenix_TodayDateShift []byte

type LuaScriptsStruct struct {
	LuaScriptName string
	LuaScript     []byte
}

// Add all files into one slice
func loadFenixLuaScripts() (fenixLuaScripts []LuaScriptsStruct) {

	//fenixLuaScripts = append(fenixLuaScripts, LuaScriptsStruct{"date", date})
	fenixLuaScripts = append(fenixLuaScripts, LuaScriptsStruct{"fenix_ControlledUniqueId", fenix_ControlledUniqueId})
	fenixLuaScripts = append(fenixLuaScripts, LuaScriptsStruct{"fenix_ControlledUniqueId", fenix_ControlledUniqueId})
	fenixLuaScripts = append(fenixLuaScripts, LuaScriptsStruct{"fenix_TodayDateShift", fenix_TodayDateShift})

	return fenixLuaScripts
}
