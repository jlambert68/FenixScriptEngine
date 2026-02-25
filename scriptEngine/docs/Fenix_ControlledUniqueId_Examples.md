# Fenix_ControlledUniqueId - Supported Replacements and Processing

This file summarizes what `scriptEngine/src/Fenix_ControlledUniqueId.lua` supports.

## Supported Replacements

The function replaces these tokens in the input text (`inputTable[3][1]`):

| Token | Replacement format | Placeholder usage (without attribute names) | Placeholder usage (with attribute names) | Example result |
|---|---|---|---|---|
| `%YYYY-MM-DD%` | Current date (`YYYY-MM-DD`) | `{{Fenix.ControlledUniqueId(@Date=%YYYY-MM-DD%)}}` | `{{Fenix.ControlledUniqueId(@InputText=Date=%YYYY-MM-DD%)}}` | `Date=2026-02-24` |
| `%YYYYMMDD%` | Current date (`YYYYMMDD`) | `{{Fenix.ControlledUniqueId(@Date=%YYYYMMDD%)}}` | `{{Fenix.ControlledUniqueId(@InputText=Date=%YYYYMMDD%)}}` | `Date=20260224` |
| `%YYMMDD%` | Current date (`YYMMDD`) | `{{Fenix.ControlledUniqueId(@Date=%YYMMDD%)}}` | `{{Fenix.ControlledUniqueId(@InputText=Date=%YYMMDD%)}}` | `Date=260224` |
| `%hh:mm:ss%` | Current time (`HH:MM:SS`) | `{{Fenix.ControlledUniqueId(@Time=%hh:mm:ss%)}}` | `{{Fenix.ControlledUniqueId(@InputText=Time=%hh:mm:ss%)}}` | `Time=13:07:09` |
| `%hh.mm.ss%` | Current time (`HH.MM.SS`) | `{{Fenix.ControlledUniqueId(@Time=%hh.mm.ss%)}}` | `{{Fenix.ControlledUniqueId(@InputText=Time=%hh.mm.ss%)}}` | `Time=13.07.09` |
| `%hhmmss%` | Current time (`HHMMSS`) | `{{Fenix.ControlledUniqueId(@Time=%hhmmss%)}}` | `{{Fenix.ControlledUniqueId(@InputText=Time=%hhmmss%)}}` | `Time=130709` |
| `%hhmm%` | Current time (`HHMM`) | `{{Fenix.ControlledUniqueId(@Time=%hhmm%)}}` | `{{Fenix.ControlledUniqueId(@InputText=Time=%hhmm%)}}` | `Time=1307` |
| `%nnn...%` | Zero-padded random digits (length = count of `n`) | `{{Fenix.ControlledUniqueId(@Rand=%nnnnn%)}}` | `{{Fenix.ControlledUniqueId(@InputText=Rand=%nnnnn%)}}` | `Rand=84018` |
| `%a(length; seed)%` | Random lowercase string | `{{Fenix.ControlledUniqueId(@Rand=%a(5; 11)%)}}` | `{{Fenix.ControlledUniqueId(@InputText=Rand=%a(5; 11)%)}}` | `Rand=ynint` |
| `%A(length; seed)%` | Random uppercase string | `{{Fenix.ControlledUniqueId(@Rand=%A(5; 10)%)}}` | `{{Fenix.ControlledUniqueId(@InputText=Rand=%A(5; 10)%)}}` | `Rand=OPNEV` |

## Processing Behavior

| Behavior | What the Lua file does | Placeholder usage (without attribute names) | Placeholder usage (with attribute names) | Example result |
|---|---|---|---|---|
| Default array index | If `inputTable[2]` is empty, index defaults to `{1}` | `{{Fenix.ControlledUniqueId(@Rand=%nnn%)}}` | `{{Fenix.ControlledUniqueId(@InputText=Rand=%nnn%)}}` | `Rand=840` |
| Single array index only | More than one array index is rejected | `{{Fenix.ControlledUniqueId[@1,@2](@X)}}` | `{{Fenix.ControlledUniqueId[@arrayIndexes=[@1,@2]](@InputText=@X)}}` | Error: `[1,2]` |
| Entropy table type check | `inputTable[4]` must be a table | `{{Fenix.ControlledUniqueId(@X)}(@<invalid>)}` | `{{Fenix.ControlledUniqueId(@InputText=@X)}(@useEntropy=@<invalid>)}` | Error: `Error - entropy is not of type 'Table', but is of type 'string'.` |
| Entropy param 1 validation | `entropyTable[1]` must be boolean (or string convertable to bool) | `{{Fenix.ControlledUniqueId(@X)}(@nope,@0)}` | `{{Fenix.ControlledUniqueId(@InputText=@X)}(@useEntropy=@nope, @extraEntropy=@0)}` | Error: `Error - entropy parameter no. 1 must be of type 'Boolean'...` |
| Entropy param 2 validation | `entropyTable[2]` must be integer (or string convertable to integer) | `{{Fenix.ControlledUniqueId(@X)}(@true,@nope)}` | `{{Fenix.ControlledUniqueId(@InputText=@X)}(@useEntropy=@true, @extraEntropy=@nope)}` | Error: `Error - entropy parameters no. 2 must be of type 'Integer'...` |
| Text type check | `inputTable[3][1]` must be a string | `{{Fenix.ControlledUniqueId(@<non-string>)}}` | `{{Fenix.ControlledUniqueId(@InputText=@<non-string>)}}` | Error: `textToProcess must be a string, got number` |

## Seed/Entropy Details in This Lua File

- For `%nnn...%` random numbers, seed is based on:
  - `arrayPositionTable[1] + entropyTable[2]`
- For `%a(length; seed)%` and `%A(length; seed)%`, the token's own `seed` argument is used for random string generation.
- `entropyTable[1]` (boolean) is validated, but not used in random generation logic inside this Lua file.

## Full Example From Lua File Style

Input text:

```text
Date: %YYYY-MM-DD%, Date: %YYYYMMDD%, Date: %YYMMDD%, Time: %hh:mm:ss%, Time: %hhmmss%, Time: %hhmm%, Random Number: %nnnnn%, Random String: %a(5; 11)%, Random String Uppercase: %A(5; 10)%, Time: %hh:mm:ss%, Time: %hh.mm.ss%
```

Function input table:

```lua
{"Fenix_ControlledUniqueId", {0}, {"Date: %YYYY-MM-DD%, Date: %YYYYMMDD%, Date: %YYMMDD%, Time: %hh:mm:ss%, Time: %hhmmss%, Time: %hhmm%, Random Number: %nnnnn%, Random String: %a(5; 11)%, Random String Uppercase: %A(5; 10)%, Time: %hh:mm:ss%, Time: %hh.mm.ss% "}, {true, 0}}
```

Result:

```text
Date: 2026-02-24, Date: 20260224, Date: 260224, Time: 13:07:09, Time: 130709, Time: 1307, Random Number: 84018, Random String: ynint, Random String Uppercase: OPNEV, Time: 13:07:09, Time: 13.07.09
```

## Notes

- Date/time example outputs above were generated with fixed time `2026-02-24 13:07:09` for reproducibility.
- Random outputs shown are deterministic for the same input table values.
