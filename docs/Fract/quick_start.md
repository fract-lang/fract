# QUICK START TO FRACT

+ [Keywords](https://github.com/fract-lang/fract/blob/main/docs/Fract/keywords.md)
+ [DataTypes](https://github.com/fract-lang/fract/blob/main/docs/Fract/data_types.md)
+ [Operators](https://github.com/fract-lang/fract/blob/main/docs/Fract/operators.md)

## Comments
Fract does not support multiline comments. ``#`` Is used for comments.
### Examples
```
555 + 5 # Comment
```

## Process Priority
Fract adheres to transaction priority!
## Examples
```
5 + 2 * 2     # 9
(5 + 2) * 2   # 14
```

## Print
Fract uses nothing to print. Just write the value and print it out.

> If the statement starts with parentheses, it's "definitely" considered a print statement!

### Examples
```
5555 + 1 # Print 5556
```

```
var x int32 := 5
x # Print 5
```

## Variables
### Definition

> A value must be given when defining a variable!

> Variable names are must comply to [naming conventions](https://github.com/fract-lang/fract/blob/main/docs/Fract/naming_conventions.md).

> Cannot change values of const variables!

#### Syntax
```
var [NAME] [DATA_TYPE] := [VALUE]
```
```
const [NAME] [DATA_TYPE] := [VALUE]
```
#### Examples
```
const Pi float32 := 3.14
```
```
var Fibonacci.First int32 := 1
```

### Set Defined
#### Syntax
```
[NAME] := [VALUE]
```
#### Examples
```
var a int32 := 45    # Value is 45
a := 1               # Value is 1
```

### Deletion Defines
You can free space and customize usage by deleting definitions from memory.

> Fract does not allow direct memory management! You can contribute to usage by deleting only memorized definitions.

#### Syntax
```
del [NAME] [NAME] [NAME]...
```
#### Examples
```
var a int32 := 4
var a int32 := 5       # Error, defined name 'a'

------------------

var a int32 := 4
del a                  # Remove 'a' from memory
var a int32 := 5       # No error, a is 5
```
```
var a int32 := 0
var b int32 := 0
del a b                # Remove 'a' and 'b'
                       # No defined variables
```
