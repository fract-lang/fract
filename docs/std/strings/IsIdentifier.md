# ``IsIdentifier`` function

## Description
A string is considered a valid identifier if it only contains alphanumeric letters (a-z) and (0-9),
or underscores (_). A valid identifier cannot start with a number, or contain any spaces.

## Define
```
protected func IsIdentifier(str)
```

## Parameters
+ ``str`` <br>
String.

## Examples
```
open std.strings

print(strings.IsIdentifier("Fract"))       # true
print(strings.IsIdentifier("Fract lang"))  # false
print(strings.IsIdentifier("3ger"))        # false
print(strings.IsIdentifier("G4r_cia"))     # true
```
