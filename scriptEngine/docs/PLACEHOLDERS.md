# ScriptEngine Placeholder Functions

This document describes placeholders currently supported in the Go-first ScriptEngine path.

## Implementation Mapping

| Jira Story | Runtime Placeholder | Go Function | Go File |
|---|---|---|---|
| `TemplateEngine.TodayShiftDay(arg)` | `Fenix.TodayShiftDay` | `goFenixTodayShiftDay` | `go_placeholder_fenix_today_shift_day.go` |
| `TemplateEngine.ControlledUniqueId(args[])` | `Fenix.ControlledUniqueId` | `goFenixControlledUniqueID` | `go_placeholder_fenix_controlled_unique_id.go` |
| `TemplateEngine.RandomPositiveDecimalValue(args[])` | `Fenix.RandomPositiveDecimalValue` | `goFenixRandomPositiveDecimalValue` | `go_placeholder_fenix_random_positive_decimal_value.go` |
| `TemplateEngine.RandomPositiveDecimalValue.Sum(args[])` | `Fenix.RandomPositiveDecimalValue.Sum` | `goFenixRandomPositiveDecimalValueSum` | `go_placeholder_fenix_random_positive_decimal_value_sum.go` |
| _(Lua placeholder)_ | `HappyLuaTime` | `HappyLuaTime` | `luaFunctions/HappyLuaTime.lua` |

Shared files:

- `go_placeholder_dispatcher.go`
- `go_placeholder_registration.go`
- `go_placeholder_time_provider.go`
- `go_placeholder_fenix_random_positive_decimal_helpers.go`

## Execution Flow

1. Placeholder text is parsed in `placeholderReplacementEngine.match(...)`.
2. Parsed input becomes `[placeholder, functionName, arrayIndexes, arguments, useEntropy, extraEntropy]`.
3. Go handler dispatch is attempted first (`executeGoPlaceholderFunction(...)`).
4. If no Go handler exists, legacy Lua execution is used.

## Syntax And Parser Constraints

General syntax:

```text
{{Function.Name[optionalArrayIndexes](arg1, arg2, ...)}(useEntropyFromTestCaseExecutionUuid, extraEntropy)}
```

Constraints from the current parser implementation:

- Function arguments are split on commas.
- No quoting/escaping support for commas inside one argument.
- Dot notation in function names is normalized to underscore names internally.

## Supported Functions

### 1) `Fenix.TodayShiftDay`

Contract:

- Exactly one integer argument: `(shiftDays)`.
- Array indexes are not supported.
- Output format is `YYYY-MM-DD` in local time.

Examples:

```text
{{Fenix.TodayShiftDay(0)}}
{{Fenix.TodayShiftDay(-1)}}
{{Fenix.TodayShiftDay(1)}}
```

### 2) `Fenix.ControlledUniqueId`

Contract:

- Exactly three function arguments:
  - `textToProcess`
  - `useEntropyFromExecutionUUID` (`true`/`false`)
  - `extraEntropy` (integer)
- Optional array index: default `1`, max one index.

Important behavior:

- Date/time tokens are replaced using local current time.
- Random Jira tokens are deterministic from array index + entropy.
- Legacy non-Jira random formats are not replaced.
- Entropy for this function is derived from function arguments 2 and 3.

Supported date/time tokens:

- `%YYYY-MM-DD%`
- `%YYYYMMDD%`
- `%YYMMDD%`
- `%hh:mm:ss%`
- `%hh.mm.ss%`
- `%hhmmss%`
- `%hhmm%`
- `%mmss%`
- `%ms%`
- `%us%`
- `%ns%`

Supported random Jira tokens:

- `%n(length)%`
- `%a(length)%`
- `%A(length)%`
- `%aA(length)%`
- `%an(length)%`
- `%An(length)%`
- `%aAn(length)%`

Examples (parser-safe, no comma inside first argument):

```text
{{Fenix.ControlledUniqueId(%YYYY-MM-DD%, true, 0)}}
{{Fenix.ControlledUniqueId[2](ID-%n(5)%-%a(4)%-%A(4)%, true, 5)}}
{{Fenix.ControlledUniqueId(Year=YYYY-Month=MM-Day=DD, false, 1)}}
```

### 3) `Fenix.RandomPositiveDecimalValue`

Contract:

- Optional single array index, default `1`.
- Exactly five function arguments:
  - `IntegerPrecision`
  - `FractionPrecision`
  - `IntegerFieldWidth`
  - `FractionFieldWidth`
  - `DecimalPointCharacter` (single character)

Behavior:

- Deterministic random generation using array index + dispatcher entropy.
- Integer and fraction padding applied from field widths.
- Decimal separator replaced with `DecimalPointCharacter`.

Examples:

```text
{{Fenix.RandomPositiveDecimalValue(2, 3, 2, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue[2](2, 3, 2, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue(2, 3, 4, 4, ",")}}
```

### 4) `Fenix.RandomPositiveDecimalValue.Sum`

Contract:

- Array index list supports one or more integers; negatives subtract.
- Default index list is `[1]` when empty.
- Exactly same five function arguments as `Fenix.RandomPositiveDecimalValue`.

Behavior:

- Generates per-index deterministic values.
- Positive indexes add, negative indexes subtract.
- Applies same padding and decimal-point replacement as value variant.
- Negative sum formatting keeps leading zeros (example test expectation: `-044.613`).

Examples:

```text
{{Fenix.RandomPositiveDecimalValue.Sum[1](2, 3, 2, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue.Sum[-1,2](2, 3, 3, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue.Sum[1,2](2, 3, 4, 4, ",")}}
```

### 5) `HappyLuaTime` (Lua Placeholder)

Contract:

- No array indexes.
- No function arguments.
- Returns a string in the format:
  - `My name is Lua and the time is HH:MM:SS`

Example:

```text
{{HappyLuaTime()}}
```

## TestData Placeholder Handling

`ParseAndFormatPlaceholders(...)` supports:

- Preferred: `{{TestData.Context.ColumnName}}`
- Legacy: `{{Context.TestData.ColumnName}}`

Lookup key is always the final segment (`ColumnName`).

## Validation Coverage

Validation and deterministic behavior are covered in:

- `scriptEngine/go_placeholder_dispatcher_test.go`
- `scriptEngine/go_placeholder_fenix_today_shift_day_test.go`
- `scriptEngine/go_placeholder_fenix_controlled_unique_id_test.go`
- `scriptEngine/go_placeholder_fenix_random_positive_decimal_value_test.go`
- `scriptEngine/go_placeholder_fenix_random_positive_decimal_value_sum_test.go`
- `placeholderReplacementEngine/placeholderReplacementEngine_test.go`

## Per-Placeholder Example Files

- `Fenix_TodayShiftDay_Examples.md`
- `Fenix_ControlledUniqueId_Examples.md`
- `Fenix_RandomPositiveDecimalValue_Examples.md`
- `Fenix_RandomPositiveDecimalValue_Sum_Examples.md`
- `HappyLuaTime_Examples.md`
