# ``Index`` function

## Description
Returns index if element is found in array, returns -1 if not.

## Define
```
protected func Index(array, element, start=0)
```

## Parameters
+ ``array`` <br>
Array.
+ ``element`` <br>
Element to find.
+ ``start`` <br>
Start index of search.

## Examples
```
open std.array

const array = [0, 9, 1, 2, 4, 2, 6]

print(array.Index(array, 5))  # -1
print(array.Index(array, 2))  # 3
```
