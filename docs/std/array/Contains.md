# ``Contains`` function

## Description
Returns true if element is contained in array, returns false if not.

## Define
```
protected func Contains(array, element, start=0)
```

## Parameters
+ ``array`` <br>
Array.
+ ``element`` <br>
Element to check.
+ ``start`` <br>
Start index of search.

## Examples
```
open std.array

const array = [0, 9, 1, 2, 4, 2, 6]

print(array.Contains(array, 5))  # false
print(array.Contains(array, 2))  # true
```
