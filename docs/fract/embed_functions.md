# Embed Functions

Embed functions are can use in everywhere and declared by interpreter.

## exit
Exit process with exit code.

### Define
```
protected func exit(code=0)
```

### Parameters
+ ``code = 0`` <br>
Exit code.

### Examples
```
exit()
exit(1)
```

## float
Convert to float an object.

### Define
```
protected func float(object)
```

### Parameters
+ ``object`` <br>
Object to parse float.

### Examples
```
float("3435")     # 3435
float("34.35")    # 34.35
```

## input
Returns input from CLI as string.

### Define
```
protected func input(message="")
```

### Parameters
+ ``message=""`` <br>
Input message.

### Examples

```
var name = input()
var name = input("Your name: ")
```

## int
Convert to integer an object.

### Define
```
protected func int(object, type="parse")
```

### Parameters
+ ``object`` <br>
Object to parse int.
+ ``type`` <br>
Parse type.

### Flags
#### ``type`` parameter
+ ``parse`` <br>
Parse value to integer.
+ ``strcode`` <br>
Parse one char string to char code.

### Examples
```
int("3435")               # 3435
int(34.35)                # 34
int("", type="strcode")   # -1
int("A", type="strcode")  # 65
```

## len
Calculate length of object.

### Define
```
protected func len(object)
```

### Parameters
+ ``object`` <br>
Object to calculate length.

### Examples
```
len(["This", "is", "array"]) # Length of 3
len("String")                # Length of 6
```

## make
Return array has integer elements by size. Minimum value is: ``0``.

### Define
```
protected func make(size)
```

### Parameters
+ ``size`` <br>
Size of array.

### Examples
```
make(0)         # []
make(2)         # [0 0]
make(10)        # [0 0 0 0 0 0 0 0 0 0]
```

## print
Print values to CLI.

### Define
```
protected func print(value="", fin="\n")
```

### Parameters
+ ``value=""`` <br>
Value to print.
+ ``fin="\n"`` <br>
Value to print after "value" parameter.

### Examples
```
print("Hello World")

# OUTPUT
Hello World
```
```
print("Hello", fin=" ")
print("World")

# OUTPUT
Hello World
```

## range
Return range as array.

### Define
```
protected func range(start, to, step=1)
```

### Parameters
+ ``start`` <br>
Start value.
+ ``to`` <br>
Limit value.
+ ``step=1`` <br>
Increase/decrease rate for a step.

### Examples
```
range(1, 10)         # [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]
range(1, 10, step=2) # [ 1, 3, 5, 7, 9 ]
range(10, 1)         # [ 10, 9, 8, 7, 6, 5, 4, 3, 2, 1 ]
```

## string
Convert to string an object.

### Define
```
protected func string(object, type="object")
```

### Parameters
+ ``object`` <br>
Object to parse int.
+ ``type`` <br>
Parse type.

### Flags
#### "type" parameter
+ ``object`` <br>
Parse object to string data.
+ ``parse`` <br>
Parse value to string.
+ ``bytecode`` <br>
Parse string from byte or byte array.

### Examples
```
string(3435)                  # {[{3435 0}] false}
string(3435, type="parse")    # 3435
string(65, type="bytecode")   # A
```

## append
Append source values to destination array.

### Define
```
protected func append(dest, ...src)
```

### Parameters
+ ``dest`` <br>
Destination array.
+ ``src`` <br>
Source values.

### Examples
```
const binary = [0, 1]
var decimal = append(binary, 2, 3, 4, 5, 6, 7, 8, 9)
print(decimal)

# [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
```
