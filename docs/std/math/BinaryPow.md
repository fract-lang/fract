# ``BinaryPow`` function

## Description
Returns power using binary exponentiation of number.

## Define
```
protected func BinaryPow(x, y, z=false)
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

print(math.BinaryPow(5, 5))      # 3125
print(math.BinaryPow(5, 5, 2))   # 1
```
