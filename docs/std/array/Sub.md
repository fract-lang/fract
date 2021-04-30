# ``Sub`` function

## Description
Returns subbed array.

## Define
```
protected func Sub(array, start, length)
```

## Parameters
+ ``array`` <br>
Array.
+ ``start`` <br>
Start index of take.
+ ``length`` <br>
Length.

## Examples
```
open std.array

const array = [0, 9, 1, 2, 4, 2, 6]

print(array.Sub(array, 2, 5))   # [1 2 4 2 6]
print(array.Sub(array, 5, 524)) # [2 6]
```
