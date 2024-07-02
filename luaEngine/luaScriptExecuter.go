package luaScriptEngine

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/yuin/gopher-lua"
	"hash/crc32"
	"log"
	"strings"
)

// InitiateLuaScriptEngine
// Initiate the Lua Script Engine
func InitiateLuaScriptEngine(luaScriptFiles [][]byte) (err error) {

	// Load Fenix Lua Script files
	var fenixLuaScripts [][]byte
	fenixLuaScripts = loadFenixLuaScripts()

	// Concatenate Fenix Lua scripts with Domain supported Lua scripts
	for _, fenixLuaScript := range fenixLuaScripts {
		luaScriptFiles = append(luaScriptFiles, fenixLuaScript)
	}

	// Save all script into one byte array
	luaScriptFilesAsByteArray = luaScriptFiles

	// Initiate Gopher-Lua Script Engine state variables
	luaState = lua.NewState()

	// Load standard libraries
	luaState.OpenLibs()

	// List preloaded libraries
	listLibraries(luaState)

	// Remove or stub unsafe functions
	//luaState.SetGlobal("io", lua.LNil) // Remove the 'io' library
	//luaState.SetGlobal("os", lua.LNil) // Remove the 'os' library
	//luaState.SetGlobal("dofile", lua.LNil)   // Remove the 'dofile' function
	//luaState.SetGlobal("loadfile", lua.LNil) // Remove the 'loadfile' function

	// Load the Lua scripts
	for _, luaScriptFileAsByteArray := range luaScriptFilesAsByteArray {
		//err = luaState.Load(string(luaScriptFileAsByteArray))
		//_, err = luaState.Load(bytes.NewReader(luaScriptFileAsByteArray), "script")
		loadAndExecuteScript(luaState, luaScriptFileAsByteArray)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Now list all the functions in the global environment
	global := luaState.Get(lua.GlobalsIndex) // Access the global table
	if tbl, ok := global.(*lua.LTable); ok {
		tbl.ForEach(func(key, value lua.LValue) {
			if _, ok := value.(*lua.LFunction); ok {
				fmt.Printf("Function: %s\n", key.String())
			}
		})
	}

	// Replace the default 'print' with our custom function
	luaState.SetGlobal("print", luaState.NewFunction(customPrint))

	// Run a Lua script
	if err := luaState.DoString(`print("Hello", "world", 123)`); err != nil {
		fmt.Println("Error running Lua script:", err)
	}

	return err
}

// customPrint replaces the default Lua print function to capture output in Go.
func customPrint(L *lua.LState) int {
	top := L.GetTop()
	parts := make([]string, top)
	for i := 1; i <= top; i++ {
		parts[i-1] = L.ToStringMeta(L.Get(i)).String()
	}
	fmt.Println(strings.Join(parts, "\t"))
	return 0 // Number of results
}

// CloseDownLuaScriptEngine
// Close down the Lua Script Engine in a correct way
func CloseDownLuaScriptEngine() {
	luaState.Close()
}

func listLibraries(L *lua.LState) {
	// Directly access the global 'package' table
	packageTable := L.GetGlobal("package")
	if packageTable.Type() == lua.LTNil {
		fmt.Println("No 'package' global found.")
		return
	}

	// Access 'preload' table inside 'package'
	preloadTable := L.GetField(packageTable, "preload")
	if tbl, ok := preloadTable.(*lua.LTable); ok {
		fmt.Println("Preloaded libraries:")
		tbl.ForEach(func(key lua.LValue, value lua.LValue) {
			fmt.Println(key.String())
		})
	} else {
		fmt.Println("No preloaded libraries found or 'preload' is not a table.")
	}
}

// ExecuteLuaScriptBasedOnPlaceholder
// Execute a specific Lua function
func ExecuteLuaScriptBasedOnPlaceholder(inputParameterArray []interface{}, testCaseExecutionUuid string) (responseValue string) {

	var err error
	var luaFunctionToCall string
	var addExtraEntropyValue uint64
	var useEntropyFromTestCaseExecutionUuid bool
	var entropyToUse uint64

	luaFunctionToCall = inputParameterArray[1].(string)
	useEntropyFromTestCaseExecutionUuid = inputParameterArray[4].(bool)
	addExtraEntropyValue = inputParameterArray[5].(uint64)

	// Decide how much entropy to use
	if useEntropyFromTestCaseExecutionUuid == true {

		// Create entropy form TestCaseExecutionUuid by converting string into 32-bit hash
		var entropyBasedOnTestCaseExecutionUuid uint32
		entropyBasedOnTestCaseExecutionUuid = crc32.ChecksumIEEE([]byte(testCaseExecutionUuid))

		// Create final entropy
		entropyToUse = uint64(entropyBasedOnTestCaseExecutionUuid) + addExtraEntropyValue

	} else {
		entropyToUse = addExtraEntropyValue
	}

	// Create Lua Entropy table
	entropyTable := luaState.NewTable() // Instantiate the Lua table

	// Append a boolean value
	luaState.SetTable(entropyTable, lua.LNumber(1), lua.LBool(useEntropyFromTestCaseExecutionUuid))

	// Append an integer value
	luaState.SetTable(entropyTable, lua.LNumber(2), lua.LNumber(entropyToUse))

	var goArrayToBeConvertedIntoLuaTable []interface{}
	goArrayToBeConvertedIntoLuaTable = inputParameterArray[1:4]

	// Create the Lua Input Table
	var luaInputTable *lua.LTable
	luaInputTable = convertToLuaTableRecursively(luaState, goArrayToBeConvertedIntoLuaTable)

	// Append 'entropyTable' to 'luaInputTable'

	// Get number of elements in table
	var numberOfElementsInTable int
	numberOfElementsInTable = luaInputTable.Len()

	// Append an entropy table to 'luaInputTable'
	luaState.SetTable(luaInputTable, lua.LNumber(numberOfElementsInTable+1), entropyTable)

	/*
		// Add the 'entropyBasedOnTestCaseExecutionUuid' to the entropy position, which is position 4
		originalValue := luaState.GetField(luaInputTable, "4")
		fmt.Println("Original value at position 4:", originalValue)


		// Add the Go value to the Lua value at position 4
		lv, ok := originalValue.(lua.LNumber)
		if ok == true {
			newValue := lua.LNumber(entropyBasedOnTestCaseExecutionUuid) + lv
			luaState.SetField(luaInputTable, "4", newValue) // Set the new value at position 4
		} else {
			// Original value is not a number
			err = errors.New(fmt.Sprintf("original value, '%s' in position 4 is not a number", originalValue))

			return err.Error()
		}



		// Access and modify the number in the nested table
		nestedTable, ok := originalValue.(*lua.LTable)
		if ok == true {
			nestedValue := luaState.GetField(nestedTable, "1")
			lv, ok := nestedValue.(lua.LNumber)
			if ok == true {
				newValue := lua.LNumber(entropyBasedOnTestCaseExecutionUuid) + lv
				luaState.SetField(nestedTable, "1", newValue) // Update the number in the nested table
			} else {
				// Original value is not a number
				err = errors.New(fmt.Sprintf("original value, '%s' in position 1 in table is not a number", nestedValue))

				return err.Error()
			}
		} else {
			// Original value is not a lua table
			err = errors.New(fmt.Sprintf("original value, '%s' in position 4 is not a table", originalValue))
			return err.Error()
		}
	*/

	// Call lua function based on Placeholder
	responseValue, err = callPlaceholderFunctionWithInputTable(luaState, luaFunctionToCall, luaInputTable)

	// If there is an error then create error message with input table as response
	if err != nil {

		responseValue = err.Error() //+ "\n" + printLuaTable(luaState, luaInputTable, "-")
		return responseValue
	}

	return responseValue
}

// printLuaTable recursively prints a Lua table and returns the result as a string
func printLuaTable(L *lua.LState, table *lua.LTable, indent string) string {
	var builder strings.Builder
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		if tbl, ok := value.(*lua.LTable); ok {
			// If the value is a table, recurse
			builder.WriteString(fmt.Sprintf("%s%s:\n", indent, key.String()))
			builder.WriteString(printLuaTable(L, tbl, indent+"  "))
		} else {
			// Otherwise, print the key and value
			builder.WriteString(fmt.Sprintf("%s%s: %s\n", indent, key.String(), value.String()))
		}
	})
	return builder.String()
}

// 'convertToLuaTableRecursively' recursively converts a Go slice of `[]interface{}` to a Lua table.
func convertToLuaTableRecursively(tempLuaState *lua.LState, goSlice []interface{}) *lua.LTable {
	var luaTable *lua.LTable
	luaTable = tempLuaState.NewTable()
	for _, item := range goSlice {
		switch v := item.(type) {
		case []interface{}:
			// If the item is a slice, recursively convert it to a Lua table
			nestedTable := convertToLuaTableRecursively(tempLuaState, v)
			tempLuaState.RawSet(luaTable, lua.LNumber(luaTable.Len()+1), nestedTable)
		default:
			// Otherwise, add the item directly to the Lua table
			tempLuaState.RawSet(luaTable, lua.LNumber(luaTable.Len()+1), lua.LString(fmt.Sprint(v)))
		}
	}
	return luaTable
}

// Loads and executes Lua code from byte slice to define functions or execute initializations
func loadAndExecuteScript(L *lua.LState, script []byte) error {
	if fn, err := L.Load(bytes.NewReader(script), ""); err != nil {
		return err
	} else {
		L.Push(fn)
		return L.PCall(0, lua.MultRet, nil)
	}
}

// Calls a Lua function with parameters and prints the result
func callPlaceholderFunctionWithInputTable(L *lua.LState, funcName string, placeholderInputTable *lua.LTable) (luaFunctionResponse string, err error) {
	//L.GetGlobal(funcName)
	//L.Push(placeholderInputTable)
	//err = L.PCall(1, 1, nil)

	// Execute the function
	err = L.CallByParam(lua.P{
		Fn:      L.GetGlobal(funcName),
		NRet:    1,
		Protect: true,
	},
		placeholderInputTable)

	if err != nil {
		return err.Error(), err
	}

	// Extract the response
	var luaResponseTable *lua.LTable
	tbl, ok := L.Get(-1).(*lua.LTable)
	L.Pop(1) // Remove the result from the stack

	if ok == true {
		luaResponseTable = tbl

		// Extract response
		var success, value, errorMessage lua.LValue
		var successAsBool bool
		var valueString string
		var errorMessageAsString string

		success = luaResponseTable.RawGetString("success")
		value = luaResponseTable.RawGetString("value")
		errorMessage = luaResponseTable.RawGetString("errorMessage")

		// Check that 'success' of type boolean, and if so then convert into a boolean
		if success.Type() != lua.LTBool {
			err = errors.New(fmt.Sprintf("In response from placeholder function: '%s' the responseTable.success is not of type Boolean. Instead the type seems to be a '%s'", funcName, value.Type().String()))

			return "", err
		} else {

			successAsBool = lua.LVAsBool(success)
		}

		// Check that 'value' of type string, and if so then convert into a string
		if value.Type() != lua.LTString {
			err = errors.New(fmt.Sprintf("In response from placeholder function: '%s' the responseTable.value is not of type String. Instead the type seems to be a '%s'", funcName, value.Type().String()))

			return "", err
		} else {

			valueString = lua.LVAsString(value)
		}

		// Check that 'errorMessage' of type string, and if so then convert into a string
		if errorMessage.Type() != lua.LTString {
			err = errors.New(fmt.Sprintf("In response from placeholder function: '%s' the responseTable.errorMessage is not of type String. Instead the type seems to be a '%s'", funcName, value.Type().String()))

			return "", err
		} else {

			errorMessageAsString = lua.LVAsString(errorMessage)
		}

		// Check if we got any error message back
		if len(errorMessageAsString) > 0 {
			return "", errors.New(errorMessageAsString)
		}

		// Check if we didn't get a OK response and the errorMessage is empty
		if len(errorMessageAsString) == 0 && successAsBool == false {
			return "", errors.New(fmt.Sprintf("'errorMessage' from function '%s' is empty but responseTable.success is 'false'. This shouldn't happen", funcName))
		}

		// Return the response value from Lua
		return valueString, nil

	} else {
		err = errors.New(fmt.Sprintf("Expected a table, but didn't get one as a response from the Lua execution for Placeholder function: '%s'", funcName))
		return "", err
	}

}
