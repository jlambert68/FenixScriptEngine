package main

import (
	"fmt"
	"github.com/jlambert68/FenixScriptEngine/scriptEngine"
	"log"
)

type placeholderExample struct {
	source      string
	description string
	input       []interface{}
}

func main() {

	var err error

	// No external Lua-Libraries
	var fenixLuaScripts []scriptEngine.LuaScriptsStruct
	fenixLuaScripts = []scriptEngine.LuaScriptsStruct{}

	// Load and initiate Lua-engine
	err = scriptEngine.InitiateLuaScriptEngine(fenixLuaScripts)
	if err != nil {
		log.Fatalln("Error", err)
	}

	runExamplesFromLuaFiles()

}

func runExamplesFromLuaFiles() {
	testCaseExecutionUuid := "f8c06f7e-0a8a-4d75-9f25-5e5fb8d2a6d3"

	examples := []placeholderExample{}
	examples = append(examples, examplesFromTodayDateShiftLua()...)
	examples = append(examples, examplesFromControlledUniqueIDLua()...)
	examples = append(examples, examplesFromRandomPositiveDecimalValueLua()...)

	fmt.Println("Running ScriptEngine examples from Lua files")
	fmt.Println("===========================================")
	for exampleIndex, example := range examples {
		response := scriptEngine.ExecuteLuaScriptBasedOnPlaceholder(example.input, testCaseExecutionUuid)

		fmt.Printf(
			"\n[%d] %s :: %s\ninput: %v\noutput: %s\n",
			exampleIndex+1,
			example.source,
			example.description,
			example.input,
			response,
		)
	}
}

func examplesFromTodayDateShiftLua() []placeholderExample {
	return []placeholderExample{
		{
			source:      "Fenix_TodayDateShift.lua",
			description: "Example invocation at end of file",
			input: buildPlaceholderInput(
				"{{Fenix.TodayShiftDay()}}",
				"Fenix_TodayShiftDay",
				[]int{},
				[]string{},
				true,
				0,
			),
		},
		{
			source:      "Fenix_TodayDateShift.lua",
			description: "Commented example invocation in file",
			input: buildPlaceholderInput(
				"{{Fenix.TodayShiftDay(-1)}}",
				"Fenix_TodayShiftDay",
				[]int{},
				[]string{"-1"},
				true,
				0,
			),
		},
	}
}

func examplesFromControlledUniqueIDLua() []placeholderExample {
	return []placeholderExample{
		{
			source:      "Fenix_ControlledUniqueId.lua",
			description: "Simple date replacement example",
			input: buildPlaceholderInput(
				"{{Fenix.ControlledUniqueId(Date: %YYYY-MM-DD%)}}",
				"Fenix_ControlledUniqueId",
				[]int{},
				[]string{"Date: %YYYY-MM-DD%"},
				true,
				0,
			),
		},
		{
			source:      "Fenix_ControlledUniqueId.lua",
			description: "Extended token replacement example",
			input: buildPlaceholderInput(
				"{{Fenix.ControlledUniqueId(Date: %YYYY-MM-DD%, Date: %YYYYMMDD%, Date: %YYMMDD%, Time: %hh:mm:ss%, Time: %hhmmss%, Time: %hhmm%, Random Number: %nnnnn%, Random String: %a(5; 11)%, Random String Uppercase: %A(5; 10)%, Time: %hh:mm:ss%, Time: %hh.mm.ss% )}}",
				"Fenix_ControlledUniqueId",
				[]int{0},
				[]string{"Date: %YYYY-MM-DD%, Date: %YYYYMMDD%, Date: %YYMMDD%, Time: %hh:mm:ss%, Time: %hhmmss%, Time: %hhmm%, Random Number: %nnnnn%, Random String: %a(5; 11)%, Random String Uppercase: %A(5; 10)%, Time: %hh:mm:ss%, Time: %hh.mm.ss% "},
				true,
				0,
			),
		},
	}
}

func examplesFromRandomPositiveDecimalValueLua() []placeholderExample {
	return []placeholderExample{
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Basic random positive decimal value",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue(2, 3)}}", "Fenix_RandomPositiveDecimalValue", []int{}, []string{"2", "3"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Array index 1",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[1](2, 3)}}", "Fenix_RandomPositiveDecimalValue", []int{1}, []string{"2", "3"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Array index 3",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[3](2, 3)}}", "Fenix_RandomPositiveDecimalValue", []int{3}, []string{"2", "3"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Array index 2",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[2](2, 3)}}", "Fenix_RandomPositiveDecimalValue", []int{2}, []string{"2", "3"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Sum function with one index",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue.Sum[1](2, 3)}}", "Fenix_RandomPositiveDecimalValue_Sum", []int{1}, []string{"2", "3"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Sum function with string-like negative/positive indexes",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue.Sum[-1,2](2, 3)}}", "Fenix_RandomPositiveDecimalValue_Sum", []int{-1, 2}, []string{"2", "3"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Sum function with subtraction",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue.Sum[1,-2](2, 3)}}", "Fenix_RandomPositiveDecimalValue_Sum", []int{1, -2}, []string{"2", "3"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Sum function with only negative indexes",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue.Sum[-1,-2](2, 3)}}", "Fenix_RandomPositiveDecimalValue_Sum", []int{-1, -2}, []string{"2", "3"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Sum function with three indexes",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue.Sum[1,2,3](2, 3)}}", "Fenix_RandomPositiveDecimalValue_Sum", []int{1, 2, 3}, []string{"2", "3"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Sum function with zero-padding",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue.Sum[1,2,3](2, 3, 2, 3)}}", "Fenix_RandomPositiveDecimalValue_Sum", []int{1, 2, 3}, []string{"2", "3", "2", "3"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Sum function with larger zero-padding",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue.Sum[1,2,3](2, 3, 4, 4)}}", "Fenix_RandomPositiveDecimalValue_Sum", []int{1, 2, 3}, []string{"2", "3", "4", "4"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Basic value with index 1",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[1](2, 3)}}", "Fenix_RandomPositiveDecimalValue", []int{1}, []string{"2", "3"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Value with small integer/fraction sizes",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue(1, 2)}}", "Fenix_RandomPositiveDecimalValue", []int{}, []string{"1", "2"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Value with array index 2 and small sizes",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[2](1, 2)}}", "Fenix_RandomPositiveDecimalValue", []int{2}, []string{"1", "2"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "One decimal place",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue(1, 1)}}", "Fenix_RandomPositiveDecimalValue", []int{}, []string{"1", "1"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "One decimal place with extra entropy",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue(1, 1)}(true,1)}", "Fenix_RandomPositiveDecimalValue", []int{}, []string{"1", "1"}, true, 1),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "One decimal place with index + extra entropy",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[1](1, 1)}(true,1)}", "Fenix_RandomPositiveDecimalValue", []int{1}, []string{"1", "1"}, true, 1),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Zero integer digits",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue(0, 1)}}", "Fenix_RandomPositiveDecimalValue", []int{}, []string{"0", "1"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "No decimal part",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[1](1, 0)}}", "Fenix_RandomPositiveDecimalValue", []int{1}, []string{"1", "0"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Zero and zero",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[1](0, 0)}}", "Fenix_RandomPositiveDecimalValue", []int{1}, []string{"0", "0"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Zero with padding",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[1](0, 0, 2, 3)}}", "Fenix_RandomPositiveDecimalValue", []int{1}, []string{"0", "0", "2", "3"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Padding when decimals exist",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[1](0, 2, 3, 4)}}", "Fenix_RandomPositiveDecimalValue", []int{1}, []string{"0", "2", "3", "4"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Larger value with padding",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[1](2, 2, 3, 4)}}", "Fenix_RandomPositiveDecimalValue", []int{1}, []string{"2", "2", "3", "4"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Large integer and fraction sizes",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue(6, 6)}}", "Fenix_RandomPositiveDecimalValue", []int{}, []string{"6", "6"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Large value with 10 decimals",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue(6, 10)}}", "Fenix_RandomPositiveDecimalValue", []int{}, []string{"6", "10"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Error case: one argument only",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[1](0)}}", "Fenix_RandomPositiveDecimalValue", []int{1}, []string{"0"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Error case: empty arguments",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[1]()}}", "Fenix_RandomPositiveDecimalValue", []int{1}, []string{}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Error case: three arguments",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[1](1,2,3)}}", "Fenix_RandomPositiveDecimalValue", []int{1}, []string{"1", "2", "3"}, true, 0),
		},
		{
			source:      "Fenix_RandomPositiveDecimalValue.lua",
			description: "Error case: too many array indexes",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[1,2](2,3)}}", "Fenix_RandomPositiveDecimalValue", []int{1, 2}, []string{"2", "3"}, true, 0),
		},
	}
}

func buildPlaceholderInput(
	placeholder string,
	functionName string,
	arrayIndexes []int,
	functionArguments []string,
	useEntropyFromTestCaseExecutionUuid bool,
	addExtraEntropyValue uint64,
) []interface{} {
	var arrayIndexesAsInterface []interface{}
	for _, arrayIndex := range arrayIndexes {
		arrayIndexesAsInterface = append(arrayIndexesAsInterface, arrayIndex)
	}

	var functionArgumentsAsInterface []interface{}
	for _, functionArgument := range functionArguments {
		functionArgumentsAsInterface = append(functionArgumentsAsInterface, functionArgument)
	}

	return []interface{}{
		placeholder,
		functionName,
		arrayIndexesAsInterface,
		functionArgumentsAsInterface,
		useEntropyFromTestCaseExecutionUuid,
		addExtraEntropyValue,
	}
}
