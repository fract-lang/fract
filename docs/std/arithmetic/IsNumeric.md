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
open std.arithmetic

print(arithmetic.IsNumeric(5))      # True
print(arithmetic.IsNumeric("544"))  # True
print(arithmetic.IsNumeric("F"))    # False
```
