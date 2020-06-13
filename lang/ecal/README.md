ECAL - Event Condition Action Language
--
ECAL is a language to create a rule based system which reacts to events provided that a defined condition holds:

Event -> Condition -> Action

The condition and action part are defined by rules called event sinks which are the core constructs of ECAL.

Notation
--
Source code is Unicode text encoded in UTF-8. Single language statements are separated by a semicolon or a newline.

Constant values are usually enclosed in double quotes "" or single quotes '', both supporting escape sequences. Constant values can also be provided as raw strings prefixing a single or double quote with an 'r'. A raw string can contain any character including newlines and does not contain escape sequences.

Blocks are denoted with curley brackets. Most language constructs (conditions, loops, etc.) are very similar to other languages.

Event Sinks
--
Sinks are should have unique names which identify them and the following attributes:

Attribute | Description
-|-
kindmatch  | Matching condition for event kind e.g. db.op.TableInsert. A list of strings in dot notation which describes event kinds. May contain `*` characters as wildcards.
scopematch | Matching condition for event cascade scope e.g. db.dbRead db.dbWrite. A list of strings in dot notation which describe the scopes which are required for this sink to trigger.
statematch | Match on event state: A simple map of required key / value states in the event state. `NULL` values can be used as wildcards (i.e. match is only on key).
priority | Priority of the sink. Sinks of higher priority are executed first. The higher the number the lower the priority - 0 is the highest priority.
suppresses | A list of sink names which should be suppressed if this sink is executed.

Example:
```
sink "mysink"
    r"
    A comment describing the sink.
    "
    kindmatch [ foo.bar.* ],
    scopematch [ "data.read", "data.write" ],
    statematch { a : 1, b : NULL },
    priority 0,
    suppresses [ "myothersink" ]
    {
      <ECAL Code>
    }
```

Events which match
...

Events which don't match
...

Function
--
Functions define reusable pieces of code dedicated to perform a particular task based on a set of given input values. In ECAL functions are first-class citizens in that they can be assigned to variables, passed as arguments, immediately invoked or deferred for last execution. Each parameter can have a default value which is by default NULL.

Example:
```
func myfunc(a, b, c=1) {
  <ECAL Code>
}
```

Comments
--
Comments are defined with `#` as single line comments and `/*` `*/` for multiline comments.
Single line comments will comment all characters after the `#` until the next newline.
```
/*
  Multi line comment
  Some comment text
*/

# Single line comment

a := 1 # Single line comment after a statement
```

Constant Values
--
Constant values are used to initialize variables or as operands in expressions.

Numbers can be expressed in all common notations:
Formatting|Description
-|-
123|Normal integer
123.456|With decimal point
1.234560e+02|Scientific notation

Strings can be normal quoted stings which interpret backslash escape characters:
```
\a → U+0007 alert or bell
\b → U+0008 backspace
\f → U+000C form feed
\n → U+000A line feed or newline
\r → U+000D carriage return
\t → U+0009 horizontal tab
\v → U+000b vertical tab
\\ → U+005c backslash
\" → U+0022 double quote
\uhhhh → a Unicode character whose codepoint can be expressed in 4 hexadecimal digits. (pad 0 in front)
```

Normal quoted strings also interpret inline expressions escaped with `{}`:
```
"Foo bar {1+2}"
```
Inline expression may also specify number formatting:
```
"Foo bar {1+2}.2f"
```
Formatting|Description
-|-
{}.f|With decimal point full precision
{}.3f|Decimal point with precision 3
{}.5w3f|5 Width with decimal point with precision 3
{}.e|Scientific notation

Strings can also be expressed in raw form which will not interpret any escape characters.
```
r"Foo bar {1+2}"
```

Expression|Value
-|-
`"foo'bar"`| `foo'bar`
`'foo"bar'`| `foo"bar`
`'foo\u0028bar'`| `foo(bar`
`"foo\u0028bar"`| `foo(bar`
`"Foo bar {1+2}"`| `Foo bar 3`
`r"Foo bar {1+2}"`| `Foo bar {1+2}`

Variable Assignments
--
A variable is a storage location for holding a value. Variables can hold single values (strings and numbers) or structures like an array or a map. Variables names can only contain [a-zA-Z] and [a-zA-Z0-9] from the second character.

A variable is assigned with the assign operator ':='
```
a := 1
b := "test"
c := [1,2,3]
d := {1:2,3:4}
```
Multi-assignments are possible using lists:
```
[a, b] := [1, 2]
```

Expressions
--
Variables and constants can be combined with operators to form expressions. Boolean expressions can also be formed with variables:
```
a := 1 + 2 * 5
b := a > 10
c := a == 11
d := false or c
```

Operators
--
The following operators are available:

Boolean: `and`, `or`, `not`, `>`, `>=`, `<`, `<=`, `==`, `!=`

Arithmetic: `+`, `-`, `*`, `/`, `//` (integer division), `%` (integer modulo)

String:
Operator|Description|Example
-|-|-
like|Regex match|`"Hans" like "H??s"`
hasPrefix|prefix match|`"Hans" hasPrefix "Ha"`
hasSuffix|suffix match|`"Hans" hasSuffix "ns"`

List:
Operator|Description|Example
-|-|-
in|Item is in list|`6 in [1, 6, 7]`
notin|Item is not in list|`6 notin [1, 6, 7]`

Composition structures access
--
Composition structures like lists and maps can be accessed with access operators:

Structure|Accessor|Description
-|-|-
List|variable[index]|Access the n-th element starting from 0.
Map|variable[field]|Access a map
Map|variable.field|Access a map (field name can only contain [a-zA-Z] and [a-zA-Z0-9] from the second character)
```
a := [1, 2, 3]
b := a[1] # B has the value 2

c := { "foo" : 2 }
d := c["foo"]
e := c.foo
```

Loop statements
---------------
All loops are defined as a 'for' block statement. Counting loops are defined with the 'range' function. The following code iterates from 2 until 10 in steps of 2:
```
for a in range(2, 10, 2) {
	<ECAL Code>
}
```

Conditional loops are using a condition after the for statement:
```
for a > 0 {
  <ECAL Code>
}
```

It is possible to loop over lists and even have multiple assignments:
```
for [a, b] in [[1, 1], [2, 2], [3, 3]] {

}
```
or
```
x := { "c" : 0, "a" : 2, "b" : 4}
for [a, b] in x {
  <ECAL Code>
}
```

Conditional statements
----------------------
The "if" statement specifies the conditional execution of multiple branches based on defined conditions:
```
if a == 1 {
    a := a + 1
} elif a == 2 {
    a := a + 2
} else {
    a := 99
}
```
