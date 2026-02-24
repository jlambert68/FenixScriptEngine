package scriptEngine

import (
	"fmt"
	"time"
)

var currentTimeProvider = time.Now

func init() {
	// Register built-in Go handlers at package load time.
	if err := registerDefaultGoPlaceholderFunctions(); err != nil {
		panic(err)
	}
}

// registerDefaultGoPlaceholderFunctions wires all built-in placeholders to Go handlers.
func registerDefaultGoPlaceholderFunctions() error {
	if err := RegisterGoPlaceholderFunction("Fenix_TodayShiftDay", goFenixTodayShiftDay); err != nil {
		return fmt.Errorf("failed to register placeholder 'Fenix_TodayShiftDay': %w", err)
	}
	if err := RegisterGoPlaceholderFunction("Fenix_ControlledUniqueId", goFenixControlledUniqueID); err != nil {
		return fmt.Errorf("failed to register placeholder 'Fenix_ControlledUniqueId': %w", err)
	}
	if err := RegisterGoPlaceholderFunction("Fenix_RandomPositiveDecimalValue", goFenixRandomPositiveDecimalValue); err != nil {
		return fmt.Errorf("failed to register placeholder 'Fenix_RandomPositiveDecimalValue': %w", err)
	}
	if err := RegisterGoPlaceholderFunction("Fenix_RandomPositiveDecimalValue_Sum", goFenixRandomPositiveDecimalValueSum); err != nil {
		return fmt.Errorf("failed to register placeholder 'Fenix_RandomPositiveDecimalValue_Sum': %w", err)
	}

	return nil
}

// goFenixTodayShiftDay replicates Fenix_TodayShiftDay behavior from Lua.
// It supports zero or one integer argument and returns YYYY-MM-DD.
func goFenixTodayShiftDay(input GoPlaceholderInput) (string, error) {
	if len(input.ArrayIndexes) > 0 {
		return "", fmt.Errorf("Error - array index is not supported. arrayIndexes: %v", input.ArrayIndexes)
	}

	var shiftDays int
	switch len(input.Arguments) {
	case 0:
		shiftDays = 0
	case 1:
		parsedShift, err := parseSingleIntegerArgument(input.Arguments[0])
		if err != nil {
			return "", fmt.Errorf("Error - function argument is not an Integer: '%s'", input.Arguments[0])
		}
		shiftDays = parsedShift
	default:
		return "", fmt.Errorf("Error - more than 1 parameter argument. arguments: %v", input.Arguments)
	}

	// Work with a date-only value in local time to avoid clock-time side effects.
	now := currentTimeProvider().In(time.Local)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return today.AddDate(0, 0, shiftDays).Format("2006-01-02"), nil
}
