# Imports

There are two types of imports in Fract.

## Standard Library Import
Standard libraries can be imported with standard import. <br>

### Examples
```
open std
```
```
open std.constants
```

## Local Import
You can import your other directories with local import. <br>

### Your directory
```
myproject/
├── bar
|   ├── foo.fract
└── main.fract
```
### Examples
```
open "bar"
```

## Naming
The last name of the directory is assumed to be the package name.

### Examaples
```
open std

print(std.Foo)
```
```
open "foo/bar"

print(bar.HelloMessage)
```

## Aliases
Are the names too long?

### Examples
```
open s std

print(s.Foo)
```
```
open b "foo/bar"

print(b.HelloMessage)
```

## Information
+ Files in the same directory as your startup file are assumed to have been imported.
+ All fract files imports from directory if import a directory.
+ If you want it to be imported, you must start the define name with a capital letter, definitions that do not start with a capital letter are not imported!
