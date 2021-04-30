# ``SwapCase`` function

## Description
Swaps cases, lowercase becomes uppercase and vice versa.

## Define
```
protected func SwapCase(str)
```

## Parameters
+ ``str`` <br>
String.

## Examples
```
open std.strings

print(strings.SwapCase("Fract"))           # fRACT
print(strings.SwapCase("Fract-languaGE"))  # fRACT-LANGUAge
print(strings.SwapCase("3* Er"))           # 3* eR
print(strings.SwapCase("FOOBAR"))          # foobar
```
