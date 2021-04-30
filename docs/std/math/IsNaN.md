# ``IsNaN`` function

## Description
Returns true if number is NaN, false if not.

## Define
```
protected func IsNaN(x)
```

## Parameters
+ ``x`` <br>
Numeric.

## Examples
```
open std.math

print(math.IsNaN(NaN)) #  true
print(math.IsNaN(300.16)) #  false
```
