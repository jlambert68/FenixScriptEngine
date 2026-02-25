# ScriptEngine Placeholder Functions

This document describes the placeholders currently supported by `scriptEngine` and how to use them in templates.

## Placeholder Syntax

General format:

```text
{{Function.Name[optionalArrayIndexes](arg1, arg2, ...)}(useEntropyFromTestCaseExecutionUuid, extraEntropy)}
```

Notes:

- Dots in function names are converted to underscores internally.
- `useEntropyFromTestCaseExecutionUuid` is optional (default: `true`).
- `extraEntropy` is optional (default: `0`).
- Arguments are split on commas, so avoid commas inside a single argument value.
- Output from each function is inserted as plain text into the template.

## Supported Functions

### 1) `Fenix.TodayShiftDay`

Description:

- Returns today's date shifted by a number of days.
- Format: `YYYY-MM-DD`.

Arguments:

- No args: shift `0` days.
- One integer arg: number of days to shift (can be negative).

Examples:

```text
{{Fenix.TodayShiftDay()}}
{{Fenix.TodayShiftDay(-1)}}
{{Fenix.TodayShiftDay(10)}(false, 0)}
```

Sample result:

```text
{{Fenix.TodayShiftDay()}} -> 2026-02-24
{{Fenix.TodayShiftDay(-1)}} -> 2026-02-23
{{Fenix.TodayShiftDay(10)}(false, 0)} -> 2026-03-06
```

### 2) `Fenix.ControlledUniqueId`

Description:

- Replaces supported date/time/random tokens in an input string.
- Deterministic output based on array index + entropy.

Arguments:

- 1 string argument (the template text to transform).

Supported tokens inside the argument:

- `%YYYY-MM-DD%`
- `%YYYYMMDD%`
- `%YYMMDD%`
- `%hh:mm:ss%`
- `%hh.mm.ss%`
- `%hhmmss%`
- `%hhmm%`
- `%nnn%` (`n` repeated = number of digits)
- `%a(length; seed)%` lowercase random string
- `%A(length; seed)%` uppercase random string

Examples:

```text
{{Fenix.ControlledUniqueId(Date-%YYYYMMDD%-%nnnn%)}}
{{Fenix.ControlledUniqueId[2](ID-%a(6; 11)%-%A(4; 10)%)}(true, 5)}
```

Sample result:

```text
{{Fenix.ControlledUniqueId(Date-%YYYYMMDD%-%nnnn%)}} -> Date-20260224-5391
{{Fenix.ControlledUniqueId[2](ID-%a(6; 11)%-%A(4; 10)%)}(true, 5)} -> ID-gbrmar-IMPV
```

### 3) `Fenix.RandomPositiveDecimalValue`

Description:

- Generates a deterministic positive decimal value.

Arguments:

- `2 args`: `(maxIntegerPartSize, numberOfDecimals)`
- `4 args`: `(maxIntegerPartSize, numberOfDecimals, integerSpace, fractionSpace)` for zero-padding.

Array index:

- Optional single array index.
- If omitted, index `1` is used.

Examples:

```text
{{Fenix.RandomPositiveDecimalValue(2, 3)}}
{{Fenix.RandomPositiveDecimalValue[2](2, 3)}}
{{Fenix.RandomPositiveDecimalValue(1, 2, 3, 4)}}
{{Fenix.RandomPositiveDecimalValue[3](2, 3)}(true, 1)}
```

Sample result:

```text
{{Fenix.RandomPositiveDecimalValue(2, 3)}} -> 90.713
{{Fenix.RandomPositiveDecimalValue[2](2, 3)}} -> 46.100
{{Fenix.RandomPositiveDecimalValue(1, 2, 3, 4)}} -> 009.7100
{{Fenix.RandomPositiveDecimalValue[3](2, 3)}(true, 1)} -> 67.843
```

### 4) `Fenix.RandomPositiveDecimalValue.Sum`

Description:

- Generates values for each provided array index and sums/subtracts them.
- Positive index adds value, negative index subtracts value.

Arguments:

- Same as `Fenix.RandomPositiveDecimalValue`:
  - `(maxIntegerPartSize, numberOfDecimals)`
  - or `(maxIntegerPartSize, numberOfDecimals, integerSpace, fractionSpace)`

Array indexes:

- One or more indexes are supported.
- If omitted, `[1]` is used.

Examples:

```text
{{Fenix.RandomPositiveDecimalValue.Sum[1,2,3](2, 3)}}
{{Fenix.RandomPositiveDecimalValue.Sum[1,-2,3](2, 2)}}
{{Fenix.RandomPositiveDecimalValue.Sum[1,2](2, 3, 4, 4)}(true, 7)}
```

Sample result:

```text
{{Fenix.RandomPositiveDecimalValue.Sum[1,2,3](2, 3)}} -> 140.238
{{Fenix.RandomPositiveDecimalValue.Sum[1,-2,3](2, 2)}} -> 48.029999999999994
{{Fenix.RandomPositiveDecimalValue.Sum[1,2](2, 3, 4, 4)}(true, 7)} -> 0131.3110
```

## TestData Placeholder

Description:

- Pulls values from `testDataPointValues` map in `placeholderReplacementEngine`.

Format:

```text
{{TestData.AnyPrefix.ColumnDataName}}
```

The final segment is used as map key (for example `TestData.Customer.FirstName` uses key `FirstName`).

Example:

```text
{{TestData.Customer.FirstName}}
{{TestData.Customer.LastName}}
{{TestData.Order.OrderId}}
```

Sample result:

```text
{{TestData.Customer.FirstName}} -> Alice
{{TestData.Customer.LastName}} -> Doe
{{TestData.Order.OrderId}} -> ORD-10045
```

## Template Example

```text
Hello {{TestData.Customer.FirstName}} {{TestData.Customer.LastName}},

RunDate: {{Fenix.TodayShiftDay(0)}}
CorrelationId: {{Fenix.ControlledUniqueId(CORR-%YYYYMMDD%-%nnnn%-%A(4; 10)%)}}
Price: {{Fenix.RandomPositiveDecimalValue(2, 2)}}
Net: {{Fenix.RandomPositiveDecimalValue.Sum[1,-2,3](2, 2)}}
```

Rendered sample:

```text
Hello Alice Doe,

RunDate: 2026-02-24
CorrelationId: CORR-20260224-5391-IMPV
Price: 90.71
Net: 48.029999999999994
```

Notes for sample results:

- Date/time-based results change with current local time.
- Random-related placeholders are deterministic for the same input + execution UUID + entropy.
- TestData sample results depend on your loaded `testDataPointValues` map.
- Values shown above were generated on February 24, 2026 using testCaseExecutionUuid `f8c06f7e-0a8a-4d75-9f25-5e5fb8d2a6d3`.
