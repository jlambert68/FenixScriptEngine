package scriptEngine

import (
	"fmt"
	"time"
)

// goFenixTodayShiftDay replicates Fenix_TodayShiftDay behavior from Lua.
// Jira contract: exactly one integer argument and returns YYYY-MM-DD.
func goFenixTodayShiftDay(input GoPlaceholderInput) (string, error) {
	if len(input.ArrayIndexes) > 0 {
		return "", fmt.Errorf("Error - array index is not supported. arrayIndexes: %v", input.ArrayIndexes)
	}

	if len(input.Arguments) != 1 {
		return "", fmt.Errorf("Error - there must be exact 1 parameter argument. arguments: %v", input.Arguments)
	}
	parsedShift, err := parseSingleIntegerArgument(input.Arguments[0])
	if err != nil {
		return "", fmt.Errorf("Error - function argument is not an Integer: '%s'", input.Arguments[0])
	}
	shiftDays := parsedShift

	// Work with a date-only value in local time to avoid clock-time side effects.
	now := currentTimeProvider().In(time.Local)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return today.AddDate(0, 0, shiftDays).Format("2006-01-02"), nil
}
