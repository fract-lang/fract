# string function

## Description
Convert to string an object.

## Define
```
protected func string(object, type="object")
```

## Parameters
+ ``object`` <br>
Object to parse int.
+ ``type`` <br>
Parse type.

## Flags
### "type" parameter
+ ``object`` <br>
Parse object to string data.
+ ``parse`` <br>
Parse value to string.

## Examples
```
string(3435)                  # {[{3435 0}] false}
string(3435, type="parse")    # 3435
```
