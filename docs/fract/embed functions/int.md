# int function

## Description
Convert to integer an object.

## Define
```
protected func int(object, type="parse")
```

## Parameters
+ ``object`` <br>
Object to parse int.
+ ``type`` <br>
Parse type.

## Flags
### "type" parameter
+ ``parse`` <br>
Parse value to integer.
+ ``strcode`` <br>
Parse one char string to char code.

## Examples
```
int("3435")               # 3435
int(34.35)                # 34
int("", type="strcode")   # -1
int("A", type="strcode")  # 65
```
