# range function

## Description
Return range as array.

## Define
```
protected func range(start, to, step=1)
```

## Parameters
+ ``start`` <br>
Start value.
+ ``to`` <br>
Limit value.
+ ``step=1`` <br>
Increase/decrease rate for a step.

## Examples
```
range(1, 10)         # [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
range(1, 10, step=2) # [ 1, 3, 5, 7, 9 ]
range(10, 1)         # [ 10, 9, 8, 7, 6, 5, 4, 3, 2, 1 ]
```
