# ScriptEngine Placeholder Functions

This document describes the placeholders currently supported by `scriptEngine` and how to use them in templates.

## Jira Story Mapping And Go File Layout

The placeholders come from `docs/Jira Epics, Stories and Subtasks.txt` and are implemented with one placeholder function per Go file:

| Jira Story | Runtime Placeholder | Go Function | Go File |
|---|---|---|---|
| `TemplateEngine.TodayShiftDay(arg)` | `Fenix.TodayShiftDay` | `goFenixTodayShiftDay` | `go_placeholder_fenix_today_shift_day.go` |
| `TemplateEngine.ControlledUniqueId(args[])` | `Fenix.ControlledUniqueId` | `goFenixControlledUniqueID` | `go_placeholder_fenix_controlled_unique_id.go` |
| `TemplateEngine.RandomPositiveDecimalValue(args[])` | `Fenix.RandomPositiveDecimalValue` | `goFenixRandomPositiveDecimalValue` | `go_placeholder_fenix_random_positive_decimal_value.go` |
| `TemplateEngine.RandomPositiveDecimalValue.Sum(args[])` | `Fenix.RandomPositiveDecimalValue.Sum` | `goFenixRandomPositiveDecimalValueSum` | `go_placeholder_fenix_random_positive_decimal_value_sum.go` |

Supporting/shared code:

- `go_placeholder_registration.go` registers all placeholder handlers.
- `go_placeholder_time_provider.go` defines the injectable time source for deterministic tests.
- `go_placeholder_fenix_random_positive_decimal_helpers.go` contains shared decimal helper logic used by value and sum variants.

## Placeholder Syntax

General format:

```text
{{Function.Name[optionalArrayIndexes](arg1, arg2, ...)}(useEntropyFromTestCaseExecutionUuid, extraEntropy)}
```

Notes:

- Dots in function names are converted to underscores internally.
- Entropy tail values are parsed by the shared dispatcher.
- For `Fenix.ControlledUniqueId`, Jira uses three function arguments: `(text, useEntropyFromTestCaseExecutionUuid, extraEntropy)`.

## Supported Functions

### 1) `Fenix.TodayShiftDay`

Description:

- Returns today's date shifted by a number of days.
- Format: `YYYY-MM-DD`.

Arguments:

- Exactly one integer argument: shift days.

Examples:

```text
{{Fenix.TodayShiftDay(0)}}
{{Fenix.TodayShiftDay(-1)}}
{{Fenix.TodayShiftDay(10)}}
```

### 2) `Fenix.ControlledUniqueId`

Description:

- Replaces supported date/time/random tokens in an input string.
- Deterministic output based on array index plus entropy.

Arguments:

- Exactly three arguments:
  - `textToProcess` (string)
  - `useEntropyFromTestCaseExecutionUuid` (`true`/`false`)
  - `extraEntropy` (integer)

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

Examples:

```text
{{Fenix.ControlledUniqueId(%YYYY-MM-DD%, true, 0)}}
{{Fenix.ControlledUniqueId[2](ID-%n(5)%-%a(4)%-%A(4)%, true, 5)}}
{{Fenix.ControlledUniqueId(%Year: YYYY, Month: MM, Day: DD%, false, 1)}}
```

### 3) `Fenix.RandomPositiveDecimalValue`

Description:

- Generates a deterministic positive decimal value.

Arguments:

- Exactly five arguments:
  - `IntegerPrecision`
  - `FractionPrecision`
  - `IntegerFieldWidth`
  - `FractionFieldWidth`
  - `DecimalPointCharacter`

Array index:

- Optional single array index.
- If omitted, index `1` is used.

Examples:

```text
{{Fenix.RandomPositiveDecimalValue(2, 3, 2, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue[2](2, 3, 2, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue(2, 3, 4, 4, ",")}}
```

### 4) `Fenix.RandomPositiveDecimalValue.Sum`

Description:

- Generates values for each provided array index and sums/subtracts them.
- Positive index adds value, negative index subtracts value.

Arguments:

- Exactly the same five arguments as `Fenix.RandomPositiveDecimalValue`.

Array indexes:

- One or more indexes are supported.
- If omitted, `[1]` is used.

Examples:

```text
{{Fenix.RandomPositiveDecimalValue.Sum[1](2, 3, 2, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue.Sum[-1,2](2, 3, 3, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue.Sum[1,2](2, 3, 4, 4, ",")}}
```

## TestData Placeholder

Description:

- Pulls values from `testDataPointValues` map in `placeholderReplacementEngine`.

Format:

```text
{{TestData.AnyPrefix.ColumnDataName}}
```

The final segment is used as map key (for example `TestData.Customer.FirstName` uses key `FirstName`).

## Template Example

```text
Hello {{TestData.Customer.FirstName}} {{TestData.Customer.LastName}},

RunDate: {{Fenix.TodayShiftDay(0)}}
CorrelationId: {{Fenix.ControlledUniqueId(CORR-%YYYYMMDD%-%n(4)%-%A(4)%, true, 0)}}
Price: {{Fenix.RandomPositiveDecimalValue(2, 2, 3, 2, ".")}}
Net: {{Fenix.RandomPositiveDecimalValue.Sum[1,-2,3](2, 2, 4, 2, ".")}}
```

Notes:

- Date/time-based results change with current local time.
- Random-related placeholders are deterministic for the same input plus execution UUID and entropy values.
- `Fenix.ControlledUniqueId` uses Jira token formats only.
