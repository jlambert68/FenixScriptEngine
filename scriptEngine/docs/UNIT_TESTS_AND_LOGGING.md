# Unit Tests And Logging

This document maps current unit tests to behavior and logging in the repository.

## Test Packages

- `scriptEngine`
- `placeholderReplacementEngine`

## ScriptEngine Tests

### Dispatcher

File: `scriptEngine/go_placeholder_dispatcher_test.go`

Covers:

- Go handler dispatch for known functions.
- Unknown function fallback handling.
- Parse validation for entropy input types.
- Entropy calculation from `(useEntropy, extraEntropy)` tail.

Logging:

- `logDispatcherInputMatrix(...)`
- `logDispatcherExecutionResult(...)`
- `logDispatcherParseResult(...)`

### TodayShiftDay

File: `scriptEngine/go_placeholder_fenix_today_shift_day_test.go`

Covers:

- Basic shifts `0`, `-1`, `+1`.
- Input validation: array index not allowed, non-integer, wrong arg count.

Logging:

- `logPlaceholderInputMatrix(...)`
- `logPlaceholderExecutionResult(...)`

### ControlledUniqueId

File: `scriptEngine/go_placeholder_fenix_controlled_unique_id_test.go`

Covers:

- Date/time token replacement.
- Jira random token replacement patterns.
- Determinism.
- Entropy argument behavior.
- Unsupported legacy token behavior.
- Input validation.

Logging:

- `logPlaceholderInputMatrix(...)`
- `logPlaceholderExecutionResult(...)`

### RandomPositiveDecimalValue

File: `scriptEngine/go_placeholder_fenix_random_positive_decimal_value_test.go`

Covers:

- Deterministic output.
- Zero-padding and decimal formatting patterns.
- Decimal point replacement.
- Input validation errors.

Logging:

- `logPlaceholderInputMatrix(...)`
- `logPlaceholderExecutionResult(...)`

### RandomPositiveDecimalValue.Sum

File: `scriptEngine/go_placeholder_fenix_random_positive_decimal_value_sum_test.go`

Covers:

- Deterministic summed output.
- Add/subtract behavior for signed indexes.
- Pattern/format assertions.
- Negative padded formatting regression (`-044.613`).
- Input validation errors.

Logging:

- `logPlaceholderInputMatrix(...)`
- `logPlaceholderExecutionResult(...)`

### Shared Logging Helpers

File: `scriptEngine/go_placeholder_input_matrix_logger_test.go`

Defines shared helpers used by placeholder tests:

- `logPlaceholderInputMatrix(...)`
- `logPlaceholderExecutionResult(...)`

## PlaceholderReplacementEngine Tests

File: `placeholderReplacementEngine/placeholderReplacementEngine_test.go`

Covers:

- New TestData format resolution: `TestData.Context.Column`.
- Legacy TestData format resolution: `Context.TestData.Column`.
- Malformed TestData reference handling.

Logging:

- `logParseAndFormatInput(...)`
- `logParseAndFormatOutput(...)`

## Running Tests With Logs

Use verbose mode to print input/output logs:

```bash
go test -v ./placeholderReplacementEngine ./scriptEngine
```
