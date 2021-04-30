# ``Substring`` function

## Description
Returns subbed string.

## Define
```
protected func Substring(str, start, length)
```

## Parameters
+ ``str`` <br>
String.
+ ``start`` <br>
Start index of take.
+ ``length`` <br>
Length.

## Examples
```
open std.strings

print(strings.Substring("Fract", 2, 3))           # act
print(strings.Substring("Fract-languaGE", 3, 7))  # ct-lang
print(strings.Substring("3* Er", 0, 1))           # 3
print(strings.Substring("FOOBAR", 6, 3))          # 
```
