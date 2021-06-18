# Contains
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

# Index
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

# Max
Returns maximum of array.

## Define
```
protected func Max(array)
```

## Parameters
+ ``array`` <br>
Array.

## Examples
```
open std.array

const array = [0, 9, 1, 2, 4, 2, 6]

print(array.Max(array))  # 9
```

# Mean
Returns mean of array.

## Define
```
protected func Mean(array)
```

## Parameters
+ ``array`` <br>
Array.

## Examples
```
open std.array

const array = [0, 9, 1, 2, 4, 2, 6]

print(array.Mean(array))  # 3.428571
```

# Min
Returns minimum of array.

## Define
```
protected func Min(array)
```

## Parameters
+ ``array`` <br>
Array.

## Examples
```
open std.array

const array = [0, 9, 1, 2, 4, 2, 6]

print(array.Min(array))  # 0
```

# Reverse
Returns reversed array.

## Define
```
protected func Reverse(array)
```

## Parameters
+ ``array`` <br>
Array.

## Examples
```
open std.array

const array = [0, 9, 1, 2, 4, 2, 6]

print(array.Reverse(array))  # [6 2 4 2 1 9 0]
```

# Sort
Returns sorted array.

## Define
```
protected func Sort(array)
```

## Parameters
+ ``array`` <br>
Array.

## Examples
```
open std.array

const array = [0, 9, 1, 2, 4, 2, 6]

print(array.Sort(array))  # [0 1 2 2 4 6 9]
```

# Sub
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

# Sum
Returns sum of all array elements.

## Define
```
protected func Sum(array)
```

## Parameters
+ ``array`` <br>
Array.

## Examples
```
open std.array

const array = [0, 9, 1, 2, 4, 2, 6]

print(array.Sum(array)) # 24
```


