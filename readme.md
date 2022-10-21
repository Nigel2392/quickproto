Fast golang protocol parser

Supports:
* Headers (Support listed values, IE: Key=[Value, Value, Value])
* Files (Supports multiple files)
* Body (Can be encoded/decoded to base64)
* Delimiters
  * No alphabetic characters from [A-Z a-z 0-9 =]

Data is split apart by the delimiter.

Say the delimiter is a `$`:

The body and header will be split by the delimiter * 4

Each key value pair will be split by the delimiter * 2

Then the key values will be split by the delimiter.

Example:
```
key1$value1&value2$$key2$value2$$$$BODYBODYBODY
```