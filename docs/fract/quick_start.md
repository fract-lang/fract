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
### Examples
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

## Statement Terminator
With the Statement terminator, you can perform multiple operations on the same line without moving to a new line.

### Examples
```
5; 2   # Print 5 and 2
```

## Range Decomposition
Until the brackets are closed, they are tokenized.

> For example, if you have endless long conditions in conditional statements, you can use parentheses to use the bottom lines!

### Examples
```
(4 +
4)        # Tokenizer Result: (4 + 4)
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
del [NAME], [NAME], [NAME],...
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
del a, b               # Remove 'a' and 'b'
                       # No defined variables
```

### Conditional Expressions
You can let the algorithm flow with the conditions. Fract offers the If-Else If-Else structure like most programming languages.
"If" is the main condition, the alternative conditions that will come later are shown as "Else If".
When one condition is fulfilled, other conditions do not. Therefore, "If" must be rewritten each time to create a different order of conditions.
#### Syntax
```
if [CONDITION]:
end
```
```
if [CONDITION]:
  # ...
elif [CONDITION]:
  # ...
elif [CONDITION]:
  # ...
end
```

The term "Else" has not been identified as a direct keyword by Fract. For "Else", a condition set to "true" should always be written at the bottom of the condition row.
```
if [CONDITION]:
  # ...
elif [CONDITION]:
  # ...
elif true:
  # ...
end
```

A condition can be given any kind of value, but it only works with true(1) and false(0).
Unlike most languages, you won't get an error even if you only enter an integer value in the condition. It looks at the value and if it is 1 it fulfills the condition.
```
var example int32 := 0
if example:
  # ...
end
```
