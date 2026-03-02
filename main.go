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

	runPlaceholderExamples()

}

func runPlaceholderExamples() {
	testCaseExecutionUuid := "f8c06f7e-0a8a-4d75-9f25-5e5fb8d2a6d3"

	examples := []placeholderExample{}
	examples = append(examples, examplesFromTodayShiftDay()...)
	examples = append(examples, examplesFromControlledUniqueID()...)
	examples = append(examples, examplesFromRandomPositiveDecimalValue()...)
	examples = append(examples, examplesFromRandomPositiveDecimalValueSum()...)

	fmt.Println("Running ScriptEngine placeholder examples")
	fmt.Println("========================================")
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

		printRandomPositiveDecimalValueSumComponents(example.input, testCaseExecutionUuid)
	}
}

func examplesFromTodayShiftDay() []placeholderExample {
	return []placeholderExample{
		{
			source:      "go_placeholder_fenix_today_shift_day.go",
			description: "Call: TemplateEngine.TodayShiftDay(0)",
			input: buildPlaceholderInput(
				"{{Fenix.TodayShiftDay(0)}}",
				"Fenix_TodayShiftDay",
				[]int{},
				[]string{"0"},
				true,
				0,
			),
		},
		{
			source:      "go_placeholder_fenix_today_shift_day.go",
			description: "Call: TemplateEngine.TodayShiftDay(-1)",
			input: buildPlaceholderInput(
				"{{Fenix.TodayShiftDay(-1)}}",
				"Fenix_TodayShiftDay",
				[]int{},
				[]string{"-1"},
				true,
				0,
			),
		},
		{
			source:      "go_placeholder_fenix_today_shift_day.go",
			description: "Call: TemplateEngine.TodayShiftDay(1)",
			input: buildPlaceholderInput(
				"{{Fenix.TodayShiftDay(1)}}",
				"Fenix_TodayShiftDay",
				[]int{},
				[]string{"1"},
				true,
				0,
			),
		},
	}
}

func examplesFromControlledUniqueID() []placeholderExample {
	return []placeholderExample{
		{
			source:      "go_placeholder_fenix_controlled_unique_id.go",
			description: "Call: TemplateEngine.ControlledUniqueId(\"%YYYY-MM-DD%\", true, 0)",
			input: buildPlaceholderInput(
				"{{Fenix.ControlledUniqueId(%YYYY-MM-DD%)}(true,0)}",
				"Fenix_ControlledUniqueId",
				[]int{},
				[]string{"%YYYY-MM-DD%"},
				true,
				0,
			),
		},
		{
			source:      "go_placeholder_fenix_controlled_unique_id.go",
			description: "Call: TemplateEngine.ControlledUniqueId(\"Date=%YYYY-MM-DD%, Time=%hh:mm:ss%, Compact=%hhmmss%\", true, 0)",
			input: buildPlaceholderInput(
				"{{Fenix.ControlledUniqueId(Date=%YYYY-MM-DD%, Time=%hh:mm:ss%, Compact=%hhmmss%)}(true,0)}",
				"Fenix_ControlledUniqueId",
				[]int{},
				[]string{"Date=%YYYY-MM-DD%, Time=%hh:mm:ss%, Compact=%hhmmss%"},
				true,
				0,
			),
		},
		{
			source:      "go_placeholder_fenix_controlled_unique_id.go",
			description: "Call: TemplateEngine.ControlledUniqueId(\"%nnnnn%-%a(5; 11)%-%A(5; 10)%\", true, 5) with arrayIndex=2",
			input: buildPlaceholderInput(
				"{{Fenix.ControlledUniqueId[2](%nnnnn%-%a(5; 11)%-%A(5; 10)%)}(true,5)}",
				"Fenix_ControlledUniqueId",
				[]int{2},
				[]string{"%nnnnn%-%a(5; 11)%-%A(5; 10)%"},
				true,
				5,
			),
		},
		{
			source:      "go_placeholder_fenix_controlled_unique_id.go",
			description: "Call: TemplateEngine.ControlledUniqueId(\"%n(4)%-%aA(4)%\", false, 1)",
			input: buildPlaceholderInput(
				"{{Fenix.ControlledUniqueId(%n(4)%-%aA(4)%)}(false,1)}",
				"Fenix_ControlledUniqueId",
				[]int{1},
				[]string{"%n(4)%-%aA(4)%"},
				false,
				1,
			),
		},
	}
}

func examplesFromRandomPositiveDecimalValue() []placeholderExample {
	return []placeholderExample{
		{
			source:      "go_placeholder_fenix_random_positive_decimal_value.go",
			description: "Call: TemplateEngine.RandomPositiveDecimalValue(2, 3, 2, 3, \".\")",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue(2, 3, 2, 3, \".\")}}", "Fenix_RandomPositiveDecimalValue", []int{}, []string{"2", "3", "2", "3", "."}, true, 0),
		},
		{
			source:      "go_placeholder_fenix_random_positive_decimal_value.go",
			description: "Call: TemplateEngine.RandomPositiveDecimalValue[2](2, 3, 2, 3, \".\")",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[2](2, 3, 2, 3, \".\")}}", "Fenix_RandomPositiveDecimalValue", []int{2}, []string{"2", "3", "2", "3", "."}, true, 0),
		},
		{
			source:      "go_placeholder_fenix_random_positive_decimal_value.go",
			description: "Call: TemplateEngine.RandomPositiveDecimalValue(1, 2, 3, 4, \".\")",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue(1, 2, 3, 4, \".\")}}", "Fenix_RandomPositiveDecimalValue", []int{}, []string{"1", "2", "3", "4", "."}, true, 0),
		},
		{
			source:      "go_placeholder_fenix_random_positive_decimal_value.go",
			description: "Call: TemplateEngine.RandomPositiveDecimalValue(1, 1, 1, 1, \".\") with entropy(true,1)",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue(1, 1, 1, 1, \".\")}(true,1)}", "Fenix_RandomPositiveDecimalValue", []int{}, []string{"1", "1", "1", "1", "."}, true, 1),
		},
		{
			source:      "go_placeholder_fenix_random_positive_decimal_value.go",
			description: "Call: TemplateEngine.RandomPositiveDecimalValue[1](0)",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue[1](0)}}", "Fenix_RandomPositiveDecimalValue", []int{1}, []string{"0"}, true, 0),
		},
	}
}

func examplesFromRandomPositiveDecimalValueSum() []placeholderExample {
	return []placeholderExample{
		{
			source:      "go_placeholder_fenix_random_positive_decimal_value_sum.go",
			description: "Call: TemplateEngine.RandomPositiveDecimalValue.Sum[1](2, 3, 2, 3, \".\")",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue.Sum[1](2, 3, 2, 3, \".\")}}", "Fenix_RandomPositiveDecimalValue_Sum", []int{1}, []string{"2", "3", "2", "3", "."}, true, 0),
		},
		{
			source:      "go_placeholder_fenix_random_positive_decimal_value_sum.go",
			description: "Call: TemplateEngine.RandomPositiveDecimalValue.Sum[-1,2](2, 3, 3, 3, \".\")",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue.Sum[-1,2](2, 3, 3, 3, \".\")}}", "Fenix_RandomPositiveDecimalValue_Sum", []int{-1, 2}, []string{"2", "3", "3", "3", "."}, true, 0),
		},
		{
			source:      "go_placeholder_fenix_random_positive_decimal_value_sum.go",
			description: "Call: TemplateEngine.RandomPositiveDecimalValue.Sum[1,2,3](2, 3, 4, 4, \".\")",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue.Sum[1,2,3](2, 3, 4, 4, \".\")}}", "Fenix_RandomPositiveDecimalValue_Sum", []int{1, 2, 3}, []string{"2", "3", "4", "4", "."}, true, 0),
		},
		{
			source:      "go_placeholder_fenix_random_positive_decimal_value_sum.go",
			description: "Call: TemplateEngine.RandomPositiveDecimalValue.Sum[1,2](2, 3, 4, 4, \",\")",
			input:       buildPlaceholderInput("{{Fenix.RandomPositiveDecimalValue.Sum[1,2](2, 3, 4, 4, \",\")}}", "Fenix_RandomPositiveDecimalValue_Sum", []int{1, 2}, []string{"2", "3", "4", "4", ","}, true, 0),
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

func printRandomPositiveDecimalValueSumComponents(input []interface{}, testCaseExecutionUUID string) {
	functionName, ok := input[1].(string)
	if ok == false || functionName != "Fenix_RandomPositiveDecimalValue_Sum" {
		return
	}

	arrayIndexesRaw, ok := input[2].([]interface{})
	if ok == false {
		return
	}
	if len(arrayIndexesRaw) == 0 {
		arrayIndexesRaw = []interface{}{1}
	}

	argumentsRaw, ok := input[3].([]interface{})
	if ok == false {
		return
	}

	useEntropyFromExecutionUUID, ok := input[4].(bool)
	if ok == false {
		return
	}

	extraEntropy, ok := input[5].(uint64)
	if ok == false {
		return
	}

	functionArguments := make([]string, 0, len(argumentsRaw))
	for _, argument := range argumentsRaw {
		functionArguments = append(functionArguments, fmt.Sprint(argument))
	}

	fmt.Println("sum components:")
	for _, arrayIndexRaw := range arrayIndexesRaw {
		arrayIndex, ok := arrayIndexRaw.(int)
		if ok == false {
			continue
		}

		absArrayIndex := absoluteInt(arrayIndex)
		componentInput := buildPlaceholderInput(
			fmt.Sprintf("{{Fenix.RandomPositiveDecimalValue[%d](...)}}", absArrayIndex),
			"Fenix_RandomPositiveDecimalValue",
			[]int{absArrayIndex},
			functionArguments,
			useEntropyFromExecutionUUID,
			extraEntropy,
		)

		componentValue := scriptEngine.ExecuteLuaScriptBasedOnPlaceholder(componentInput, testCaseExecutionUUID)
		appliedSign := "+"
		if arrayIndex < 0 {
			appliedSign = "-"
		}

		fmt.Printf("  index %d -> value %s (applied %s%s)\n", arrayIndex, componentValue, appliedSign, componentValue)
	}
}

func absoluteInt(value int) int {
	if value < 0 {
		return -value
	}

	return value
}
