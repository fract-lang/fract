# Abs
Returns absolute of value.

## Define
```
protected func Abs(x)
```

## Parameters
+ ``x`` <br>
Numeric to get absolute.

## Examples
```
open std.math

print(math.Abs(5))      # 5
print(math.Abs(-544))   # 544
```

# Ceil
Returns the smallest (closest to negative infinity).

## Define
```
protected func Ceil(x)
```

## Parameters
+ ``x`` <br>
Numeric.

## Examples
```
open std.math

print(math.Ceil(-23.11)) #  -23
print(math.Ceil(300.16)) #  301
print(math.Ceil(300.72)) #  301
```

# Degrees
Returns radians to degrees.

## Define
```
protected func Degrees(radians)
```

## Parameters
+ ``radians`` <br>
Radians to parse.

## Examples
```
open std.math

print(math.Degrees(3))  # 171.8873
print(math.Degrees(-3)) # -171.8873
print(math.Degrees(0))  # 0.0
```

# Fact
Returns factorial of number.

## Define
```
protected func Fact(x)
```

## Parameters
+ ``x`` <br>
Numeric.

## Examples
```
open std.math

print(math.Fact(3))   # 6
print(math.Fact(-3))  # NaN
print(math.Fact(35))  # 1.033314e+40
```

# Floor
Returns the floor of x as an integral.

## Define
```
protected func Floor(x)
```

## Parameters
+ ``x`` <br>
Numeric.

## Examples
```
open std.math

print(math.Floor(-23.11)) #  -24
print(math.Floor(300.16)) #  300
print(math.Floor(300.72)) #  300
```

# IsNaN
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

# IsNumeric
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

# Max
Returns maximum number.

## Define
```
protected func Max(x, y)
```

## Parameters
+ ``x`` <br>
First numeric.
+ ``y`` <br>
Second numeric.

## Examples
```
open std.math

print(math.Max(4, 1)) #  4
print(math.Max(1, 5)) #  5
```

# Min
Returns minimum number.

## Define
```
protected func Min(x, y)
```

## Parameters
+ ``x`` <br>
First numeric.
+ ``y`` <br>
Second numeric.

## Examples
```
open std.math

print(math.Min(4, 1)) #  1
print(math.Min(3, 5)) #  3
```

# Pow
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

# Radians
Returns degrees to radians.

## Define
```
protected func Radians(degrees)
```

## Parameters
+ ``degrees`` <br>
Degrees to parse.

## Examples
```
open std.math

print(math.Radians(3))  # 0.05235988
print(math.Radians(-3)) # -0.05235988
print(math.Radians(0))  # 0.0
```

# Round
Returns balue rounded to the nearest integer.

## Define
```
protected func Round(x)
```

## Parameters
+ ``x`` <br>
Numeric.

## Examples
```
open std.math

print(math.Round(4.43)) #  4
print(math.Round(1.63)) #  2
```

# Sqrt
Returns the square root.

## Define
```
protected func Sqrt(x)
```

## Parameters
+ ``x`` <br>
Numeric.

## Examples
```
open std.math

print(math.Sqrt(4.43)) #  2.104757
print(math.Sqrt(5))    #  2.236069
```

# ToNegative
Reverse positive number to negative number.

## Define
```
protected func ToNegative(x)
```

## Parameters
+ ``x`` <br>
Numeric.

## Examples
```
open std.math

print(math.ToNegative(4))   #  -4
print(math.ToNegative(-5))  #  -5
```

# ToPositive
Reverse negative number to positive number.

## Define
```
protected func ToPositive(x)
```

## Parameters
+ ``x`` <br>
Numeric.

## Examples
```
open std.math

print(math.ToPositive(4))   #  4
print(math.ToPositive(-5))  #  5
```
