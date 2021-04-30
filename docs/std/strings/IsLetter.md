# ``IsLetter`` function

## Description
Returns true if char is letter, false if not.

## Define
```
protected func IsLetter(char)
```

## Parameters
+ ``char`` <br>
Char.

## Examples
```
open std.strings

print(strings.IsLetter("F"))  # true
print(strings.IsLetter("f"))  # true
print(strings.IsLetter("*"))  # false
print(strings.IsLetter("4"))  # false
```
