# Placeholder Attributes

This file lists all user-entered parameters supported by each placeholder function in `scriptEngine`.

## Function Parameters

| Function | Array Index Part `[ ... ]` | Function Arguments `( ... )` | Optional Entropy Tail `(useEntropy, extraEntropy)` |
|---|---|---|---|
| `Fenix.TodayShiftDay` | Not supported (must be empty) | `()` or `(shiftDays)` where `shiftDays` is an integer (for example `-1`, `0`, `10`) | Optional. `useEntropy` = `true/false`, `extraEntropy` = non-negative integer |
| `Fenix.ControlledUniqueId` | `[]` or `[index]` where `index` is an integer | One text argument. Text may contain tokens like `%YYYYMMDD%`, `%nnnn%`, `%a(5; 11)%`, `%A(4; 10)%` | Optional. `useEntropy` = `true/false`, `extraEntropy` = non-negative integer |
| `Fenix.RandomPositiveDecimalValue` | `[]` or `[index]` where `index` is an integer | Exactly `2` or `4` integer arguments: `(maxIntegerPartSize, numberOfDecimals)` or `(maxIntegerPartSize, numberOfDecimals, integerSpace, fractionSpace)` | Optional. `useEntropy` = `true/false`, `extraEntropy` = non-negative integer |
| `Fenix.RandomPositiveDecimalValue.Sum` | `[i1,i2,...]` integers (can be negative). Empty means default `[1]` | Exactly `2` or `4` integer arguments: `(maxIntegerPartSize, numberOfDecimals)` or `(maxIntegerPartSize, numberOfDecimals, integerSpace, fractionSpace)` | Optional. `useEntropy` = `true/false`, `extraEntropy` = non-negative integer |

## TestData Placeholder

| Placeholder Type | Format | User-entered part |
|---|---|---|
| TestData lookup | `{{TestData.AnyPrefix.ColumnDataName}}` | `ColumnDataName` (final segment) |

## Input Constraints

- `useEntropy` must be lowercase `true` or `false`.
- `extraEntropy` must be digits only (`0`, `1`, `25`, ...).
- Function argument parsing is comma-based, so commas inside a single argument are not supported by the current parser.
- Function names in template syntax use dots (`Fenix.X`), and are mapped internally to underscore names (`Fenix_X`).
- In the attribute-visible examples below: mandatory attributes start with uppercase, optional attributes start with lowercase.

## Entropy Calculation

```text
Entropy = extraEntropy
if useEntropyFromExecutionUUID == true {
  Entropy = crc32(testCaseExecutionUUID) + extraEntropy
}
```

- If the entropy tail is omitted, default values are `useEntropyFromExecutionUUID=true` and `extraEntropy=0`.
- When entropy tail is `false`, the execution UUID is ignored.
- `{{...}(false, N)}` gives `Entropy = N`.
- `{{...}(false)}` gives `Entropy = 0` (default extra entropy).
- Implemented in:
  - `go_placeholder_dispatcher.go:140`
  - `go_placeholder_dispatcher.go:141`
- Effect:
  - Random-based placeholders (`ControlledUniqueId`, `RandomPositiveDecimalValue`, `.Sum`) become deterministic without UUID contribution, based only on array index + `extraEntropy`.
  - `TodayShiftDay` is not entropy-driven, so `false` has no practical effect there.

## Attribute Examples With Results

### `Fenix.TodayShiftDay`

| Attribute Set | Example |
|---|---|
| Array Index Part | _none_ (not supported) |
| Function Arguments | `()` or `(-1)` |
| Entropy Tail | `(false, 0)` (optional) |

Examples with results:

```text
{{Fenix.TodayShiftDay()}} -> 2026-02-24
{{Fenix.TodayShiftDay(-1)}} -> 2026-02-23
{{Fenix.TodayShiftDay(10)}(false, 0)} -> 2026-03-06
```

Attribute-visible view:

| Placeholder | Attribute-visible placeholder | Result |
|---|---|---|
| `{{Fenix.TodayShiftDay()}}` | `{{Fenix.TodayShiftDay()}}` | `2026-02-24` |
| `{{Fenix.TodayShiftDay(-1)}}` | `{{Fenix.TodayShiftDay(daysShift=-1)}}` | `2026-02-23` |
| `{{Fenix.TodayShiftDay(10)}(false, 0)}` | `{{Fenix.TodayShiftDay(daysShift=10)}(useEntropy=false, extraEntropy=0)}` | `2026-03-06` |

### `Fenix.ControlledUniqueId`

| Attribute Set | Example |
|---|---|
| Array Index Part | `[]` or `[2]` |
| Function Arguments | `(Date-%YYYYMMDD%-%nnnn%)` |
| Entropy Tail | `(true, 5)` (optional) |

Examples with results:

```text
{{Fenix.ControlledUniqueId(Date-%YYYYMMDD%-%nnnn%)}} -> Date-20260224-5391
{{Fenix.ControlledUniqueId[2](ID-%a(6; 11)%-%A(4; 10)%)}(true, 5)} -> ID-gbrmar-IMPV
```

Attribute-visible view:

| Placeholder | Attribute-visible placeholder | Result |
|---|---|---|
| `{{Fenix.ControlledUniqueId(Date-%YYYYMMDD%-%nnnn%)}}` | `{{Fenix.ControlledUniqueId(InputText=Date-%YYYYMMDD%-%nnnn%)}}` | `Date-20260224-5391` |
| `{{Fenix.ControlledUniqueId[2](ID-%a(6; 11)%-%A(4; 10)%)}(true, 5)}` | `{{Fenix.ControlledUniqueId[arrayIndex=2](InputText=ID-%a(6; 11)%-%A(4; 10)%)}(useEntropy=true, extraEntropy=5)}` | `ID-gbrmar-IMPV` |

### `Fenix.RandomPositiveDecimalValue`

| Attribute Set | Example |
|---|---|
| Array Index Part | `[]`, `[2]`, `[3]` |
| Function Arguments | `(2, 3)` or `(1, 2, 3, 4)` |
| Entropy Tail | `(true, 1)` (optional) |

Examples with results:

```text
{{Fenix.RandomPositiveDecimalValue(2, 3)}} -> 90.713
{{Fenix.RandomPositiveDecimalValue[2](2, 3)}} -> 46.100
{{Fenix.RandomPositiveDecimalValue(1, 2, 3, 4)}} -> 009.7100
{{Fenix.RandomPositiveDecimalValue[3](2, 3)}(true, 1)} -> 67.843
```

Attribute-visible view:

| Placeholder | Attribute-visible placeholder | Result |
|---|---|---|
| `{{Fenix.RandomPositiveDecimalValue(2, 3)}}` | `{{Fenix.RandomPositiveDecimalValue(MaxIntegerPartSize=2, NumberOfDecimals=3)}}` | `90.713` |
| `{{Fenix.RandomPositiveDecimalValue[2](2, 3)}}` | `{{Fenix.RandomPositiveDecimalValue[arrayIndex=2](MaxIntegerPartSize=2, NumberOfDecimals=3)}}` | `46.100` |
| `{{Fenix.RandomPositiveDecimalValue(1, 2, 3, 4)}}` | `{{Fenix.RandomPositiveDecimalValue(MaxIntegerPartSize=1, NumberOfDecimals=2, integerSpace=3, fractionSpace=4)}}` | `009.7100` |
| `{{Fenix.RandomPositiveDecimalValue[3](2, 3)}(true, 1)}` | `{{Fenix.RandomPositiveDecimalValue[arrayIndex=3](MaxIntegerPartSize=2, NumberOfDecimals=3)}(useEntropy=true, extraEntropy=1)}` | `67.843` |

### `Fenix.RandomPositiveDecimalValue.Sum`

| Attribute Set | Example |
|---|---|
| Array Index Part | `[1,2,3]`, `[1,-2,3]` |
| Function Arguments | `(2, 3)` or `(2, 3, 4, 4)` |
| Entropy Tail | `(true, 7)` (optional) |

Examples with results:

```text
{{Fenix.RandomPositiveDecimalValue.Sum[1,2,3](2, 3)}} -> 140.238
{{Fenix.RandomPositiveDecimalValue.Sum[1,-2,3](2, 2)}} -> 48.029999999999994
{{Fenix.RandomPositiveDecimalValue.Sum[1,2](2, 3, 4, 4)}(true, 7)} -> 0131.3110
```

Attribute-visible view:

| Placeholder | Attribute-visible placeholder | Result |
|---|---|---|
| `{{Fenix.RandomPositiveDecimalValue.Sum[1,2,3](2, 3)}}` | `{{Fenix.RandomPositiveDecimalValue.Sum[arrayIndexes=[1,2,3]](MaxIntegerPartSize=2, NumberOfDecimals=3)}}` | `140.238` |
| `{{Fenix.RandomPositiveDecimalValue.Sum[1,-2,3](2, 2)}}` | `{{Fenix.RandomPositiveDecimalValue.Sum[arrayIndexes=[1,-2,3]](MaxIntegerPartSize=2, NumberOfDecimals=2)}}` | `48.029999999999994` |
| `{{Fenix.RandomPositiveDecimalValue.Sum[1,2](2, 3, 4, 4)}(true, 7)}` | `{{Fenix.RandomPositiveDecimalValue.Sum[arrayIndexes=[1,2]](MaxIntegerPartSize=2, NumberOfDecimals=3, integerSpace=4, fractionSpace=4)}(useEntropy=true, extraEntropy=7)}` | `0131.3110` |

### `TestData` Placeholder

| Attribute Set | Example |
|---|---|
| Prefix after `TestData` | `Customer` or `Order` |
| Final key segment | `FirstName`, `OrderId` |

Examples with results:

```text
{{TestData.Customer.FirstName}} -> Alice
{{TestData.Customer.LastName}} -> Doe
{{TestData.Order.OrderId}} -> ORD-10045
```

Attribute-visible view:

| Placeholder | Attribute-visible placeholder | Result |
|---|---|---|
| `{{TestData.Customer.FirstName}}` | `{{TestData.Context=Customer.ColumnDataName=FirstName}}` | `Alice` |
| `{{TestData.Customer.LastName}}` | `{{TestData.Context=Customer.ColumnDataName=LastName}}` | `Doe` |
| `{{TestData.Order.OrderId}}` | `{{TestData.Context=Order.ColumnDataName=OrderId}}` | `ORD-10045` |

Notes:

- Sample outputs above were generated with `testCaseExecutionUuid = f8c06f7e-0a8a-4d75-9f25-5e5fb8d2a6d3`.
- Date/time values vary by runtime date and local timezone.
