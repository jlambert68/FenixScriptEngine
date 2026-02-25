package placeholderReplacementEngine

import (
	"strings"
	"testing"
)

func TestParseAndFormatPlaceholders_ShouldResolveNewTestDataOrder(t *testing.T) {
	testDataMap := map[string]string{
		"FirstName": "Alice",
	}

	_, _, pureText := ParseAndFormatPlaceholders(
		"Name: {{TestData.Customer.FirstName}}",
		&testDataMap,
		"execution-uuid",
	)

	if pureText != "Name: Alice" {
		t.Fatalf("expected resolved TestData value, got: %s", pureText)
	}
}

func TestParseAndFormatPlaceholders_ShouldDetectMalformedNewTestDataReference(t *testing.T) {
	testDataMap := map[string]string{}

	_, _, pureText := ParseAndFormatPlaceholders(
		"Name: {{TestData.}}",
		&testDataMap,
		"execution-uuid",
	)

	if strings.Contains(pureText, "is not a correct TestData-reference") == false {
		t.Fatalf("expected malformed TestData-reference message, got: %s", pureText)
	}
}

func TestParseAndFormatPlaceholders_ShouldStillResolveLegacyTestDataOrder(t *testing.T) {
	testDataMap := map[string]string{
		"FirstName": "Alice",
	}

	_, _, pureText := ParseAndFormatPlaceholders(
		"Name: {{Customer.TestData.FirstName}}",
		&testDataMap,
		"execution-uuid",
	)

	if pureText != "Name: Alice" {
		t.Fatalf("expected resolved legacy TestData value, got: %s", pureText)
	}
}
