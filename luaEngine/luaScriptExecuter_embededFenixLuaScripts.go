package luaScriptEngine

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

// Add all files into one slice
func loadFenixLuaScripts() (fenixLuaScripts [][]byte) {

	fenixLuaScripts = append(fenixLuaScripts, date)
	fenixLuaScripts = append(fenixLuaScripts, fenix_ControlledUniqueId)
	fenixLuaScripts = append(fenixLuaScripts, fenix_RandomPositiveDecimalValue)
	fenixLuaScripts = append(fenixLuaScripts, fenix_TodayDateShift)

	return fenixLuaScripts
}
