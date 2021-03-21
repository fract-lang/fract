# QUICK START TO FRACT

+ [Keywords](https://github.com/fract-lang/fract/blob/main/docs/Fract/keywords.md)
+ [Data Types](https://github.com/fract-lang/fract/blob/main/docs/Fract/data_types.md)
+ [Operators](https://github.com/fract-lang/fract/blob/main/docs/Fract/operators.md)

## Comments
``#`` Is used for singline comments.
### Examples
```
555 + 5 # Comment
```
``#> ... <#`` are used for multiline comments.
### Examples
```
#>
  Hello Function
  Desc: Print hello to screen.
  func Hello()
    ...
<#
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

### Syntax
```
[VALUE]
```

### Examples
```
5555 + 1 # Print 5556
```
```
var x int32 = 5
x # Print 5
```

## Exit Keyword
With the Exit keyword, you can end the execution with an exit code.

### Syntax
```
exit [CODE]
```

### Examples
```
exit 0
```

## Statement Terminator
With the Statement terminator, you can perform multiple operations on the same line without moving to a new line.

### Syntax
```
[STATEMENT]; [STATEMENT]; [STATEMENT];...
```

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

> Can not change values of const variables!

### Syntax
```
var [NAME] = [VALUE]
```
```
const [NAME] = [VALUE]
```
### Examples
```
const Pi = 3.14
```
```
var Fibonacci.First = 1
```

## Set Defined
### Syntax
```
[NAME] = [VALUE]
```
### Examples
```
var a = 45      # Value is 45
a = 1           # Value is 1
```

## Deletion Defines
You can free space and customize usage by deleting definitions from memory.

> Fract does not allow direct memory management! You can contribute to usage by deleting only memorized definitions.

### Syntax
```
del [NAME], [NAME], [NAME],...
```
### Examples
```
var a = 4
var a = 5        # Error, defined name 'a'

------------------

var a = 4
del a            # Remove 'a' variable from memory
var a = 5        # No error, a is 5
```
```
var a = 0
var b = 0
del a, b         # Remove 'a' and 'b' variables from memory
                 # No defined variables

------------------

# Function removing

del a()          # Remove 'a' function from memory
del a, a()       # Remove 'a' variable and function from memory
```

## Arrays
They are structures that can hold more than one value in arrays. An array has a limited size and this size is determined at creation time. <br>
Syntax for creating an array that characterizes the int32 data type with 4, 5, 6, 7 elements:
```
var array = { 4, 5, 6, 7 }  # Elements: 4, 5, 6, 7
```
The syntax for creating an array of a certain size without value:
```
var array = [5] # Elements: 0, 0, 0, 0, 0
```
The syntax for accessing an element of an array with index:
```
array[index]
```
The syntax for setting an element of an array with index:
```
array[index] = value
```

### How can you quickly use the data in an array for arithmetic operations?
```
var array = { 0, 4, 4, 2 }          # Elements: 0 4 4 2
array = array + 5                   # Elements: 5 9 9 7
```
```
var array = { 0, 4, 4, 2 }          # Elements: 0 4 4 2
var array2 = { 2, 2, 2, 2 }         # Elements: 2 2 2 2
array = array + array2              # Elements: 2 6 6 4
```

> An array can be manipulated with an arithmetic value. However, when executing with a different array, the array must have only one element or the same number of elements.

## Conditional Expressions
You can let the algorithm flow with the conditions. Fract offers the If-Else If-Else structure like most programming languages.
"If" is the main condition, the alternative conditions that will come later are shown as "Else If".
When one condition is fulfilled, other conditions do not. Therefore, "If" must be rewritten each time to create a different order of conditions.

> All kinds of data can be given, but conditioning only looks for 0 and 1. 0 is accepted as false, 1 true.

### Syntax
```
if [CONDITION]
end
```
```
if [CONDITION]
  # ...
elif [CONDITION]
  # ...
elif [CONDITION]
  # ...
else
 # ...
end
```

A condition can be given any kind of value, but it only works with true(1) and false(0).
Unlike most languages, you won't get an error even if you only enter an integer value in the condition. It looks at the value and if it is 1 it fulfills the condition.
```
var example = 0
if example
  # ...
end
```

## Loops

Repetitive operations can be done using loops.

### While Loop
The while is a loop that happens as long as the condition is met.

#### Syntax
```
for [CONDITION]
  # ...
end
```
#### Examples
```
var counter = 0
for counter <= 10
  counter
  counter = counter + 1
end
```

### Foreach Loop
You can rotate the elements of arrays one by one with the foreach loop.

#### Syntax
```
for [VARIABLE_NAME] in [VALUE]
  # ...
end
```
#### Examples
```
var t1 = { 0, 3, 2, 1, 90 }
for index in { 0, 1, 2, 3, 4 }
  t1[index]
end
```
```
var t1 = { 0, 3, 2, 1, 90 }
for item in t1
  item
end
```

### Break Keyword
With the keyword break, it is possible to override and terminate the entire loop even when the loop could still return.
#### Examples
```
var counter = 0
for counter <= 10
  counter = counter + 1
  if counter > 5
    break
  end
  counter
end

# Output: 0 1 2 3 4 5
```

### Continue Keyword
It can be used to pass the cycle to the next cycle step. If there is no next loop step, the loop is terminated.
```
for index in { 0, 1, 2, 3, 4.0 }
  if index == 1 | index == 3
    continue
  end
  var test = index
  index
end

# Output: 0, 2, 4.000000
```

## Functions
Functions are very useful for adding functionality to your code.

### Syntax
Define:
```
func [NAME]([PARAM], [PARAM], [PARAM],...)
  ...
end
```
Define with default values:
```
func [NAME]([PARAM], [PARAM]=[VALUE],[PARAM]=[VALUE],...)
  ...
end
```
Call:
```
[NAME]([PARAM], [PARAM],...)
```
Call with parameter setter:
```
[NAME]([PARAM_NAME]=[VALUE], [PARAM_NAME]=[VALUE],...)
```

### Examples
```
func range(start, to, step=1)
  var lst = [0]
  var index = 0
  for start < to
    lst = { lst, start }
    start = start + step
    index = index + 1
  end
  ret lst
end

func int.prime(x)
  if x < 2
    ret false
  end

  for y in range(to=x, start=2, step=1)
    if x % y == 0
      ret false
    end
  end
  ret true
end

int.prime(3)
```
```
func print.hello()
  "Hello"
end

print.hello()
```

### Ret Keyword
The keyword ret is used to return the value of the function.

#### Syntax
```
ret [VALUE]
```

#### Examples
```
func reverse(x)
  ret x * -1
end

reverse(-500) # Returns: 500
```
