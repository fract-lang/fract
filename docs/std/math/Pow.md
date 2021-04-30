# ``Pow`` function

## Description
Returns power of number.

## Define
```
protected func Pow(x, y, z=NaN)
```

## Parameters
+ ``x`` <br>
Base.
+ ``y`` <br>
Power.
+ ``z`` <br>
Get modulus. (Set NaN for not get modulus.)

## Examples
```
open std.math

print(math.Pow(5, 5))      # 3125
print(math.Pow(5, 5, 2))   # 1
```
