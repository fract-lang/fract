# ``IsSpace`` function

## Description
Check if all the characters in the text are whitespaces.

## Define
```
protected func IsSpace(str)
```

## Parameters
+ ``str`` <br>
String.

## Examples
```
open std.strings

print(strings.IsSpace(" "))       # true
print(strings.IsSpace("\t"))      # true
print(strings.IsSpace("3g er"))   # false
print(strings.IsSpace("\v\r\f"))  # true
```
