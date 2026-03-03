# HappyLuaTime - Examples

This file contains examples for the Lua placeholder `HappyLuaTime`.

## Signature

```text
{{HappyLuaTime()}}
```

Rules:

- Array indexes are not supported.
- No function arguments are supported.

## Valid Example

```text
{{HappyLuaTime()}}
```

Possible output:

```text
My name is Lua and the time is 14:37:05
```

Output pattern:

```text
^My name is Lua and the time is [0-9]{2}:[0-9]{2}:[0-9]{2}$
```

## Invalid Examples

```text
{{HappyLuaTime(1)}}      // no function arguments allowed
{{HappyLuaTime[1]()}}    // array indexes not supported
```

