# THE FRACT PROGRAMMING LANGUAGE

Say hello to Fract!

[Wiki](https://github.com/fract-lang/fract/wiki) <br>
[Quick Start](https://github.com/fract-lang/fract/blob/main/docs/Fract/quick_start.md)

## HOW TO COMPILE
Fract is written in Go. <br>
Run one of the scripts ``build.bat`` or ``build.ps1`` to compile.

## EXAMPLE CODE
```
for index in { 0, 1, 2, 3, 4.0 }:
  if index = 1 | index = 3:
    continue
  end
  var test float64 := index
  index
end

# Output: 0, 2, 4.000000
```

## GOALS

### Operators
- [x] #
- [x] +
- [x] -
- [x] *
- [x] /
- [x] -
- [x] ^
- [x] \
- [x] %
- [x] //
- [x] \\
- [x] &
- [x] |
- [x] >
- [x] <
- [x] >=
- [x] <=
- [x] =
- [x] <>
- [x] :=
- [x] ;

### Mutable Primitive Types
- [x] int8
- [x] int16
- [x] int32
- [x] int64
- [x] uint8
- [x] uint16
- [x] uint32
- [x] uint64
- [x] float32
- [x] float64
- [x] bool
- [ ] char
- [ ] string

### Keywords
- [x] var
- [x] const
- [x] del
- [x] end
- [ ] func
- [ ] ret
- [ ] use
- [x] if
- [x] elif
- [x] for
- [x] continue
- [x] break
- [ ] enum
- [ ] struct
- [ ] match

### Primitive Structs
- [x] Array
- [ ] Matrix

### Standard Library Objects
- [ ] vector
- [ ] fract
- [ ] pcent
- [ ] tangle
- [ ] rect
- [ ] pie
