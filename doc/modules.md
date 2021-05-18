# Module reference

Pixlet scripts have access to a couple of modules that can be very
helpful in developing applets. To make use of a module, use the `load`
statement.

The example below would load the `render.star` Starlark module,
extract the `render` identifier and make it available to the script
as the symbol `render`.

```starlark
load("render.star", "render")
```

It's also possible to assign the identifiers to a different symbol. In
this example, the `render` module is made available to the script as
the symbol `r` instead of `render`:

```starlark
load("render.star", r="render")
```

## Starlib modules

Pixlet offers a subset of the modules provided by the [Starlib
project](https://github.com/qri-io/starlib). For documentation of the
individual modules, please refer to the Starlib documentation.

| Module | Description |
| --- | --- |
| [`encoding/base64.star`](https://github.com/qri-io/starlib/tree/master/encoding/base64) | Base 64 encoding and decoding |
| [`encoding/json.star`](https://github.com/qri-io/starlib/tree/master/encoding/json) | JSON encoding and decoding |
| [`http.star`](https://github.com/qri-io/starlib/tree/master/http) | HTTP client |
| [`math.star`](https://github.com/qri-io/starlib/tree/master/math) | Mathematical functions and constants |
| [`re.star`](https://github.com/qri-io/starlib/tree/master/re) | Regular expressions |
| [`time.star`](https://github.com/qri-io/starlib/tree/master/time) | Time operations |

## Pixlet module: Cache

In addition to the Starlib modules, Pixlet offers a cache module.

| Function | Description |
| --- | --- |
| `set(key, value, ttl_seconds=60)` | Writes a key-value pair to the cache, with expiration as a TTL. |
| `get(key)` | Retrieves a value by its key. Returns `None` if `key` doesn't exist or has expired. |

Keys and values must all be string. Serialization of non-string data
is the developer's responsibility.

Example:

```starlark
load("cache.star", "cache")
def get_counter():
    i = cache.get("counter")
    if i == None:
        i = 0
    cache.set("counter", str(i + 1), ttl_seconds=3600)
    return i + 1
...
```

## Pixlet module: XPath

The xpath module lets you extract data from XML documents using
[XPath](https://en.wikipedia.org/wiki/XPath) queries.

| Function | Description |
| --- | --- |
| `loads(doc) | Parses an XML document and returns an xpath object|

On an xpath object, the following methods are available:

| Method | Description |
| --- | --- |
| `query(path)` | Retrieves text of the first tag matching the path |
| `query_all(path)` | Retrieves text of all tags matching the path |

Example:

```starlark
load("xpath.star", "xpath")
doc = "<foo><bar>bar</bar><bar>baz</bar></foo>"

def get_bars():
    x = xpath.loads(doc)
    return x.query_all("/foo/bar")
...
```


## Pixlet module: Render

The `render.star` module is where Pixlet's Widgets live. All of them
are documented in a fair bit of detail in the [widget
documentation](widgets.md).

Example:

```starlark
load("render.star", r="render")
def main():
    return r.Root(child=r.Box(width=12, height=14, color="#ff0"))
```
