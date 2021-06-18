# Print
Prints message as object string to cli.

## Define
```
protected func Print(object)
```

## Parameters
+ ``object`` <br>
Object to print.

## Examples
```
open std

std.Print([5])
std.Print(5)

# {[{5 0}] true}{[{5 0}] false}
```

# Println
Prints object as object string and new line to cli.

## Define
```
protected func Println(object)
```

## Parameters
+ ``object`` <br>
Object to print.

## Examples
```
open std

std.Println([5]) # {[{5 0}] true}
std.Println(5)   # {[{5 0}] false}
```
