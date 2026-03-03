package scriptEngine

import (
	"regexp"
	"strings"
	"testing"
)

func TestLuaPlaceholder_HappyLuaTime_ShouldReturnExpectedFormat(t *testing.T) {
	err := InitiateLuaScriptEngine([]LuaScriptsStruct{})
	if err != nil {
		t.Fatalf("failed to initiate Lua engine: %v", err)
	}
	defer CloseDownLuaScriptEngine()

	input := []interface{}{
		"{{HappyLuaTime()}}",
		"HappyLuaTime",
		[]interface{}{},
		[]interface{}{},
		true,
		uint64(0),
	}
	testCaseExecutionUUID := "f8c06f7e-0a8a-4d75-9f25-5e5fb8d2a6d3"

	t.Logf("Input [happy-lua-time]\n  Input: %v\n  TestCaseExecutionUUID: %q", input, testCaseExecutionUUID)
	response := ExecuteLuaScriptBasedOnPlaceholder(input, testCaseExecutionUUID)
	t.Logf("Output [happy-lua-time]\n  Response: %q", response)

	const expectedPrefix = "My name is Lua and the time is "
	if strings.HasPrefix(response, expectedPrefix) == false {
		t.Fatalf("expected response prefix %q, got %q", expectedPrefix, response)
	}

	pattern := regexp.MustCompile(`^My name is Lua and the time is [0-9]{2}:[0-9]{2}:[0-9]{2}$`)
	if pattern.MatchString(response) == false {
		t.Fatalf("response did not match expected time format: %q", response)
	}
}
