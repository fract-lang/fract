# Capitalize
Converts the first character to uppercase.

## Define
```
protected func Capitalize(str)
```

## Parameters
+ ``str`` <br>
String.

## Examples
```
open std.strings

print(strings.Capitalize("Fract language"))  # Fract language
print(strings.Capitalize("fract language"))  # Fract language
print(strings.Capitalize("*fract language")) # *fract language
```

# IsIdentifier
A string is considered a valid identifier if it only contains alphanumeric letters (a-z) and (0-9),
or underscores (_). A valid identifier cannot start with a number, or contain any spaces.

## Define
```
protected func IsIdentifier(str)
```

## Parameters
+ ``str`` <br>
String.

## Examples
```
open std.strings

print(strings.IsIdentifier("Fract"))       # true
print(strings.IsIdentifier("Fract lang"))  # false
print(strings.IsIdentifier("3ger"))        # false
print(strings.IsIdentifier("G4r_cia"))     # true
```

# IsLetter
Returns true if char is letter, false if not.

## Define
```
protected func IsLetter(char)
```

## Parameters
+ ``char`` <br>
Char.

## Examples
```
open std.strings

print(strings.IsLetter("F"))  # true
print(strings.IsLetter("f"))  # true
print(strings.IsLetter("*"))  # false
print(strings.IsLetter("4"))  # false
```

# IsLower
Returns true if string is lowercase, false if not.

## Define
```
protected func IsLower(str)
```

## Parameters
+ ``str`` <br>
String.

## Examples
```
open std.strings

print(strings.IsLower("a")) # true
print(strings.IsLower("A")) # false
print(strings.IsLower("*")) # false
```

# IsSpace
Check if all the characters in the text are whitespaces.

## Define
```
protected func IsSpace(str)
```

## Parameters
+ ``str`` <br>
String.

## Examples
```
open std.strings

print(strings.IsSpace(" "))       # true
print(strings.IsSpace("\t"))      # true
print(strings.IsSpace("3g er"))   # false
print(strings.IsSpace("\v\r\f"))  # true
```

# IsUpper
Returns true if string is uppercase, false if not.

## Define
```
protected func IsUpper(str)
```

## Parameters
+ ``str`` <br>
String.

## Examples
```
open std.strings

print(strings.IsUpper("a")) # false
print(strings.IsUpper("A")) # true
print(strings.IsUpper("*")) # false
```

# Lower
Returns strings as lowercase.

## Define
```
protected func Lower(str)
```

## Parameters
+ ``str`` <br>
String.

## Examples
```
open std.strings

print(strings.Lower("Fract"))               # fract
print(strings.Lower("Fract-languaGE"))      # fract-language
print(strings.Lower("3* ER"))               # 3* er
print(strings.Lower("FOOBAR"))              # foobar
```

# Substring
Returns subbed string.

## Define
```
protected func Substring(str, start, length)
```

## Parameters
+ ``str`` <br>
String.
+ ``start`` <br>
Start index of take.
+ ``length`` <br>
Length.

## Examples
```
open std.strings

print(strings.Substring("Fract", 2, 3))           # act
print(strings.Substring("Fract-languaGE", 3, 7))  # ct-lang
print(strings.Substring("3* Er", 0, 1))           # 3
print(strings.Substring("FOOBAR", 6, 3))          #
```

# SwapCase
Swaps cases, lowercase becomes uppercase and vice versa.

## Define
```
protected func SwapCase(str)
```

## Parameters
+ ``str`` <br>
String.

## Examples
```
open std.strings

print(strings.SwapCase("Fract"))           # fRACT
print(strings.SwapCase("Fract-languaGE"))  # fRACT-LANGUAge
print(strings.SwapCase("3* Er"))           # 3* eR
print(strings.SwapCase("FOOBAR"))          # foobar
```

# Upper
Returns strings as uppercase.

## Define
```
protected func Upper(str)
```

## Parameters
+ ``str`` <br>
String.

## Examples
```
open std.strings

print(strings.Upper("Fract"))               # FRACT
print(strings.Upper("Fract-languaGE"))      # FRACT-LANGUAGE
print(strings.Upper("3* Er"))               # 3* ER
print(strings.Upper("FOOBAR"))              # FOOBAR
```

