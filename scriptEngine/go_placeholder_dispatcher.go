package scriptEngine

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
	"sync"
)

type GoPlaceholderInput struct {
	// Raw placeholder as written in template, for diagnostics and logging.
	Placeholder string
	// Canonical runtime function name (dots replaced by underscores).
	FunctionName string
	// Optional array indexes from placeholder syntax.
	ArrayIndexes []int
	// Positional function arguments as normalized strings.
	Arguments []string
	// Controls whether execution UUID contributes to deterministic entropy.
	UseEntropyFromExecutionUUID bool
	// User-supplied entropy offset.
	ExtraEntropy uint64
	// Final entropy used by handlers after UUID hash + extra entropy.
	Entropy uint64
	// Original execution UUID used for deterministic entropy generation.
	TestCaseExecutionUUID string
}

type GoPlaceholderFunction func(input GoPlaceholderInput) (string, error)

var (
	goPlaceholderFunctionsMutex sync.RWMutex
	// Global registry used by ExecuteLuaScriptBasedOnPlaceholder for Go-first dispatch.
	goPlaceholderFunctions = map[string]GoPlaceholderFunction{}
)

// RegisterGoPlaceholderFunction registers or replaces a Go handler for a function name.
func RegisterGoPlaceholderFunction(functionName string, fn GoPlaceholderFunction) error {
	functionName = strings.TrimSpace(functionName)
	if functionName == "" {
		return fmt.Errorf("function name can not be empty")
	}
	if fn == nil {
		return fmt.Errorf("go placeholder function for '%s' is nil", functionName)
	}

	goPlaceholderFunctionsMutex.Lock()
	goPlaceholderFunctions[functionName] = fn
	goPlaceholderFunctionsMutex.Unlock()

	return nil
}

// executeGoPlaceholderFunction attempts to route a placeholder call to a registered Go function.
// Returns handled=false when no Go handler is registered, allowing Lua fallback.
func executeGoPlaceholderFunction(inputParameterArray []interface{}, testCaseExecutionUuid string) (responseValue string, handled bool, err error) {
	functionName, exists := tryExtractFunctionName(inputParameterArray)
	if exists == false {
		return "", false, nil
	}

	goPlaceholderFunctionsMutex.RLock()
	goFunction, exists := goPlaceholderFunctions[functionName]
	goPlaceholderFunctionsMutex.RUnlock()
	if exists == false {
		return "", false, nil
	}

	parsedInput, err := parseGoPlaceholderInput(inputParameterArray, testCaseExecutionUuid)
	if err != nil {
		return "", true, err
	}

	responseValue, err = goFunction(parsedInput)
	return responseValue, true, err
}

// tryExtractFunctionName returns the canonical function name from parsed placeholder input.
func tryExtractFunctionName(inputParameterArray []interface{}) (functionName string, exists bool) {
	if len(inputParameterArray) < 2 {
		return "", false
	}

	functionName, exists = inputParameterArray[1].(string)
	return functionName, exists
}

// parseGoPlaceholderInput converts the shared placeholder input format into a typed structure
// used by all Go handlers.
func parseGoPlaceholderInput(inputParameterArray []interface{}, testCaseExecutionUuid string) (goInput GoPlaceholderInput, err error) {
	if len(inputParameterArray) < 6 {
		return goInput, fmt.Errorf("expected at least 6 input parameters, got %d", len(inputParameterArray))
	}

	placeholder, ok := inputParameterArray[0].(string)
	if ok == false {
		return goInput, fmt.Errorf("input parameter 0 ('placeholder') must be a string")
	}

	functionName, ok := inputParameterArray[1].(string)
	if ok == false {
		return goInput, fmt.Errorf("input parameter 1 ('functionName') must be a string")
	}

	arrayIndexesRaw, ok := inputParameterArray[2].([]interface{})
	if ok == false {
		return goInput, fmt.Errorf("input parameter 2 ('arrayIndexes') must be []interface{}")
	}

	arrayIndexes := make([]int, 0, len(arrayIndexesRaw))
	for _, rawIndex := range arrayIndexesRaw {
		index, ok := rawIndex.(int)
		if ok == false {
			return goInput, fmt.Errorf("all array indexes must be integers")
		}
		arrayIndexes = append(arrayIndexes, index)
	}

	argumentsRaw, ok := inputParameterArray[3].([]interface{})
	if ok == false {
		return goInput, fmt.Errorf("input parameter 3 ('arguments') must be []interface{}")
	}

	arguments := make([]string, 0, len(argumentsRaw))
	for _, rawArg := range argumentsRaw {
		arguments = append(arguments, fmt.Sprint(rawArg))
	}

	useEntropyFromExecutionUuid, ok := inputParameterArray[4].(bool)
	if ok == false {
		return goInput, fmt.Errorf("input parameter 4 ('useEntropyFromExecutionUuid') must be bool")
	}

	extraEntropy, ok := inputParameterArray[5].(uint64)
	if ok == false {
		return goInput, fmt.Errorf("input parameter 5 ('extraEntropy') must be uint64")
	}

	entropy := extraEntropy
	if useEntropyFromExecutionUuid == true {
		entropy = uint64(crc32.ChecksumIEEE([]byte(testCaseExecutionUuid))) + extraEntropy
	}

	goInput = GoPlaceholderInput{
		Placeholder:                 placeholder,
		FunctionName:                functionName,
		ArrayIndexes:                arrayIndexes,
		Arguments:                   normalizeArguments(arguments),
		UseEntropyFromExecutionUUID: useEntropyFromExecutionUuid,
		ExtraEntropy:                extraEntropy,
		Entropy:                     entropy,
		TestCaseExecutionUUID:       testCaseExecutionUuid,
	}

	return goInput, nil
}

// normalizeArguments trims each argument and treats a single empty token as "no arguments".
func normalizeArguments(arguments []string) []string {
	if len(arguments) == 1 && strings.TrimSpace(arguments[0]) == "" {
		return []string{}
	}

	cleanedArguments := make([]string, 0, len(arguments))
	for _, argument := range arguments {
		cleanedArguments = append(cleanedArguments, strings.TrimSpace(argument))
	}

	return cleanedArguments
}

// parseSingleIntegerArgument validates a single integer-style string argument.
func parseSingleIntegerArgument(argument string) (int, error) {
	argument = strings.TrimSpace(argument)
	if argument == "" {
		return 0, fmt.Errorf("argument is empty")
	}

	value, err := strconv.Atoi(argument)
	if err != nil {
		return 0, err
	}

	return value, nil
}
