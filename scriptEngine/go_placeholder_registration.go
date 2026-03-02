package scriptEngine

import "fmt"

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
