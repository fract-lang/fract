# string function

## Description
Convert to string an object.

## Define
```
protected func string(object, type="object")
```

## Examples
```
string(3435)                  # {[{3435 0}] false}
string(3435, type="parse")    # 3435
```

## Flags
### "type" parameter
+ ``object`` <br>
  Parse object format.
+ ``parse`` <br>
  Parse value format.
