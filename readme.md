Fast golang protocol parser

Supports:
* Headers
* Body (Can be encoded/decoded to base64)
* Delimiters

Data is split apart by the delimiter.

Say the delimiter is a `$`:

The body and header will be split by the delimiter * 4

Each key value pair will be split by the delimiter * 2

Then the key values will be split by the delimiter.

Example:
```
key1$value1$$key2$value2$$$$BODYBODYBODY
```