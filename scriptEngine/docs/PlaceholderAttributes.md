# Placeholder Attributes

This file lists user-entered parameters supported by each placeholder function in `scriptEngine`.

## Jira Story Mapping

The current placeholder set maps to the stories in `docs/Jira Epics, Stories and Subtasks.txt`:

| Jira Story | Runtime Placeholder | Go File |
|---|---|---|
| `TemplateEngine.TodayShiftDay(arg)` | `Fenix.TodayShiftDay` | `go_placeholder_fenix_today_shift_day.go` |
| `TemplateEngine.ControlledUniqueId(args[])` | `Fenix.ControlledUniqueId` | `go_placeholder_fenix_controlled_unique_id.go` |
| `TemplateEngine.RandomPositiveDecimalValue(args[])` | `Fenix.RandomPositiveDecimalValue` | `go_placeholder_fenix_random_positive_decimal_value.go` |
| `TemplateEngine.RandomPositiveDecimalValue.Sum(args[])` | `Fenix.RandomPositiveDecimalValue.Sum` | `go_placeholder_fenix_random_positive_decimal_value_sum.go` |

## Function Parameters

| Function | Array Index Part `[ ... ]` | Function Arguments `( ... )` | Entropy Tail `(useEntropy, extraEntropy)` |
|---|---|---|---|
| `Fenix.TodayShiftDay` | Not supported (must be empty) | Exactly one integer: `(shiftDays)` | Parsed by dispatcher, not used by function behavior |
| `Fenix.ControlledUniqueId` | `[]` or `[index]` where `index` is an integer | Exactly three arguments: `(textToProcess, useEntropyFromExecutionUUID, extraEntropy)` | Parsed by dispatcher; Jira entropy values are passed as function arguments |
| `Fenix.RandomPositiveDecimalValue` | `[]` or `[index]` where `index` is an integer | Exactly five arguments: `(IntegerPrecision, FractionPrecision, IntegerFieldWidth, FractionFieldWidth, DecimalPointCharacter)` | Optional and used by dispatcher entropy calculation |
| `Fenix.RandomPositiveDecimalValue.Sum` | `[i1,i2,...]` integers (can be negative). Empty means default `[1]` | Exactly five arguments: `(IntegerPrecision, FractionPrecision, IntegerFieldWidth, FractionFieldWidth, DecimalPointCharacter)` | Optional and used by dispatcher entropy calculation |

## ControlledUniqueId Token Set (Jira Format)

Date/time tokens:

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

Random tokens:

- `%n(length)%`
- `%a(length)%`
- `%A(length)%`
- `%aA(length)%`
- `%an(length)%`
- `%An(length)%`
- `%aAn(length)%`

## TestData Placeholder

| Placeholder Type | Format | User-entered part |
|---|---|---|
| TestData lookup | `{{TestData.AnyPrefix.ColumnDataName}}` | `ColumnDataName` (final segment) |

## Input Constraints

- `DecimalPointCharacter` must be exactly one character.
- `Fenix.ControlledUniqueId` requires exactly three function arguments.
- `Fenix.RandomPositiveDecimalValue` requires exactly five function arguments.
- `Fenix.RandomPositiveDecimalValue.Sum` requires exactly five function arguments.
- Function names in template syntax use dots (`Fenix.X`) and are mapped internally to underscore names (`Fenix_X`).

## Entropy Calculation (Dispatcher)

```text
Entropy = extraEntropy
if useEntropyFromExecutionUUID == true {
  Entropy = crc32(testCaseExecutionUUID) + extraEntropy
}
```

- If entropy tail is omitted, defaults are `useEntropyFromExecutionUUID=true` and `extraEntropy=0`.
- `{{...}(false, N)}` gives `Entropy = N`.
- `{{...}(false)}` gives `Entropy = 0`.

## Attribute Examples

### `Fenix.TodayShiftDay`

```text
{{Fenix.TodayShiftDay(0)}}
{{Fenix.TodayShiftDay(-1)}}
{{Fenix.TodayShiftDay(10)}}
```

### `Fenix.ControlledUniqueId`

```text
{{Fenix.ControlledUniqueId(%YYYY-MM-DD%, true, 0)}}
{{Fenix.ControlledUniqueId[2](ID-%n(5)%-%a(4)%-%A(4)%, true, 5)}}
{{Fenix.ControlledUniqueId(%Year: YYYY, Month: MM, Day: DD%, false, 1)}}
```

### `Fenix.RandomPositiveDecimalValue`

```text
{{Fenix.RandomPositiveDecimalValue(2, 3, 2, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue[2](2, 3, 2, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue(2, 3, 4, 4, ",")}}
```

### `Fenix.RandomPositiveDecimalValue.Sum`

```text
{{Fenix.RandomPositiveDecimalValue.Sum[1](2, 3, 2, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue.Sum[-1,2](2, 3, 3, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue.Sum[1,2](2, 3, 4, 4, ",")}}
```

Notes:

- Sample values depend on current date/time, array index, execution UUID, and entropy.
- Legacy non-Jira random token forms are not supported.
