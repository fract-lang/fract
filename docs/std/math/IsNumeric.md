# ``IsNumeric`` function

## Description
Returns true if object is numeric, returns false if not.

## Define
```
protected func IsNumeric(object)
```

## Parameters
+ ``object`` <br>
Object to check.

## Examples
```
open std.math

print(math.IsNumeric(5))      # True
print(math.IsNumeric("544"))  # True
print(math.IsNumeric("F"))    # False
```
