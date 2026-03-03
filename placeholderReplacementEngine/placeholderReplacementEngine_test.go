package placeholderReplacementEngine

import (
	"strings"
	"testing"
)

func logParseAndFormatInput(t *testing.T, callLabel string, template string, testDataMap map[string]string, executionUUID string) {
	t.Helper()
	t.Logf(
		"Input [%s]\n  Template: %q\n  TestDataMap: %v\n  ExecutionUUID: %q",
		callLabel,
		template,
		testDataMap,
		executionUUID,
	)
}

func logParseAndFormatOutput(t *testing.T, callLabel string, pureText string) {
	t.Helper()
	t.Logf(
		"Output [%s]\n  PureText: %q",
		callLabel,
		pureText,
	)
}

func TestParseAndFormatPlaceholders_ShouldResolveNewTestDataOrder(t *testing.T) {
	testDataMap := map[string]string{
		"FirstName": "Alice",
	}
	template := "Name: {{TestData.Customer.FirstName}}"
	executionUUID := "execution-uuid"

	logParseAndFormatInput(t, "resolve-new-order", template, testDataMap, executionUUID)
	_, _, pureText := ParseAndFormatPlaceholders(
		template,
		&testDataMap,
		executionUUID,
	)
	logParseAndFormatOutput(t, "resolve-new-order", pureText)

	if pureText != "Name: Alice" {
		t.Fatalf("expected resolved TestData value, got: %s", pureText)
	}
}

func TestParseAndFormatPlaceholders_ShouldDetectMalformedNewTestDataReference(t *testing.T) {
	testDataMap := map[string]string{}
	template := "Name: {{TestData.}}"
	executionUUID := "execution-uuid"

	logParseAndFormatInput(t, "detect-malformed-new-order", template, testDataMap, executionUUID)
	_, _, pureText := ParseAndFormatPlaceholders(
		template,
		&testDataMap,
		executionUUID,
	)
	logParseAndFormatOutput(t, "detect-malformed-new-order", pureText)

	if strings.Contains(pureText, "is not a correct TestData-reference") == false {
		t.Fatalf("expected malformed TestData-reference message, got: %s", pureText)
	}
}

func TestParseAndFormatPlaceholders_ShouldStillResolveLegacyTestDataOrder(t *testing.T) {
	testDataMap := map[string]string{
		"FirstName": "Alice",
	}
	template := "Name: {{Customer.TestData.FirstName}}"
	executionUUID := "execution-uuid"

	logParseAndFormatInput(t, "resolve-legacy-order", template, testDataMap, executionUUID)
	_, _, pureText := ParseAndFormatPlaceholders(
		template,
		&testDataMap,
		executionUUID,
	)
	logParseAndFormatOutput(t, "resolve-legacy-order", pureText)

	if pureText != "Name: Alice" {
		t.Fatalf("expected resolved legacy TestData value, got: %s", pureText)
	}
}
