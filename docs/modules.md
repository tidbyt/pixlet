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
load("render.star", r = "render")
```

## Starlib modules

Pixlet offers a subset of the modules provided by the [Starlib
project](https://github.com/qri-io/starlib). For documentation of the
individual modules, please refer to the Starlib documentation.

| Module | Description |
| --- | --- |
| [`compress/gzip.star`](https://github.com/qri-io/starlib/blob/master/compress/gzip) | gzip decompressing |
| [`compress/zipfile.star`](https://github.com/qri-io/starlib/blob/master/zipfile) | zip decompressing |
| [`encoding/base64.star`](https://github.com/qri-io/starlib/tree/master/encoding/base64) | Base 64 encoding and decoding |
| [`encoding/csv.star`](https://github.com/qri-io/starlib/tree/master/encoding/csv) | CSV decoding |
| [`encoding/json.star`](https://github.com/qri-io/starlib/tree/master/encoding/json) | JSON encoding and decoding |
| [`hash.star`](https://github.com/qri-io/starlib/tree/master/hash) | MD5, SHA1, SHA256 hash generation  |
| [`html.star`](https://github.com/qri-io/starlib/tree/master/html) | jQuery-like functions for HTML  |
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

## Pixlet module: HMAC

This module implements the HMAC algorithm as described by [RFC 2104](https://datatracker.ietf.org/doc/html/rfc2104.html).

| Function | Description |
| --- | --- |
| `md5(key, string)` | Returns md5 hash of a string using the provided key |
| `sha1(key, string)` | Returns sha1 hash of a string using the provided key |
| `sha256(key, string)` | Returns sha256 hash of a string using the provided key |

Example:

```starlark
load("hmac.star", "hmac")

sum = hmac.md5("secret", "hello world!")
print(sum)
# Output: 0a0461e10e89506d7c31a145663bed93
```

## Pixlet module: Humanize

The `humanize` module has formatters for units to human friendly sizes. 

| Function | Description |
| --- | --- |
| `time(date)` | Lets you take a `time.Time` and spit it out in relative terms. For example, `12 seconds ago` or `3 days from now`. |
| `relative_time(date1, date2, label1?, label2?)` | Formats a time into a relative string. It takes two `time.Time`s and two labels. In addition to the generic time delta string (e.g. 5 minutes), the labels are used applied so that the label corresponding to the smaller time is applied. |
| `time_format(format, date?)` | Takes a [Java SimpleDateFormat](https://docs.oracle.com/javase/7/docs/api/java/text/SimpleDateFormat.html) and returns a [Go layout string](https://programming.guide/go/format-parse-string-time-date-example.html). If you pass it a `date`, it will apply the format using the converted layout string and return the formatted date. |
| `day_of_week(date)` | Returns an integer corresponding to the day of the week, where 0 = Sunday, 6 = Saturday. |
| `bytes(size, iec?)` | Lets you take numbers like `82854982` and convert them to useful strings like, `83 MB`. You can optionally format using IEC sizes like, `83 MiB`. |
| `parse_bytes(formatted_size)` | Lets you take strings like `83 MB` and convert them to the number of bytes it represents like, `82854982`. |
| `comma(num)` | Lets you take numbers like `123456` or `123456.78` and convert them to comma-separated numbers like `123,456` or `123,456.78`. |
| `float(format, num)` | Returns a formatted number as string with options. Examples: given n = 12345.6789:  `#,###.##` => `12,345.67`, `#,###.` => `12,345`|
| `int(format, num)` | Returns a formatted number as string with options. Examples: given n = 12345: `#,###.` => `12,345`|
| `ordinal(num)` | Lets you take numbers like `1` or `2` and convert them to a rank/ordinal format strings like, `1st` or `2nd`. |
| `ftoa(num, digits?)` | Converts a float to a string with no trailing zeros. |
| `plural(quantity, singular, plural?)` | Formats an integer and a string into a single pluralized string. The simple English rules of regular pluralization will be used if the plural form is an empty string (i.e. not explicitly given).. |
| `plural_word(quantity, singular, plural?)` | Builds the plural form of an English word. The simple English rules of regular pluralization will be used if the plural form is an empty string (i.e. not explicitly given). |
| `word_series(words, conjunction)` | Converts a list of words into a word series in English. It returns a string containing all the given words separated by commas, the coordinating conjunction, and a serial comma, as appropriate. |
| `oxford_word_series(words, conjunction)` | Converts a list of words into a word series in English, using an [Oxford comma](https://en.wikipedia.org/wiki/Serial_comma). It returns a string containing all the given words separated by commas, the coordinating conjunction, and a serial comma, as appropriate. |
| `url_encode(str)` | Escapes the string so it can be safely placed inside a URL query. |
| `url_decode(str)` | The inverse of `url_encode`. Converts each 3-byte encoded substring of the form "%AB" into the hex-decoded byte 0xAB |

Example:

See [examples/humanize.star](../examples/humanize.star) for an example.

## Pixlet module: XPath

The xpath module lets you extract data from XML documents using
[XPath](https://en.wikipedia.org/wiki/XPath) queries.

| Function | Description |
| --- | --- |
| `loads(doc)` | Parses an XML document and returns an xpath object|

On an xpath object, the following methods are available:

| Method | Description |
| --- | --- |
| `query(path)` | Retrieves text of the first tag matching the path |
| `query_all(path)` | Retrieves text of all tags matching the path |
| `query_node(path)` | Retrieves the first tag matching the path as an xpath object |
| `query_all_nodes(path)` | Retrieves all tags matching the path as xpath objects |

The `query_node` and `query_all_nodes` methods allow you to recursively query the XML document, which can be useful if you need to query several tags that are nested underneath some parent tag.

Example:

```starlark
load("xpath.star", "xpath")

doc = """
<foo>
    <bar>bar</bar>
    <bar>baz</bar>
</foo>
"""

def get_bars():
    x = xpath.loads(doc)
    return x.query_all("/foo/bar")

def also_get_bars():
    x = xpath.loads(doc)
    foo = x.query_node("/foo")
    return foo.query_all("/bar")
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

## Pixlet module: Schema

The schema module provides configuration options for your app. See the [schema documentation](schema/schema.md) for more details.

Example:

See [examples/schema_hello_world.star](../examples/schema_hello_world.star) for an example.

## Pixlet module: Secret

The secret module can decrypt values that were encrypted with `pixlet encrypt`.

| Function | Description |
| --- | --- |
| `decrypt(value)` | Decrypts and returns the value when running in Tidbyt cloud. Returns `None` when running locally. Decryption will fail if the name of the app doesn't match the name that was passed to `pixlet encrypt`.  |

Example:
```starlark
load("secret.star", "secret")

ENCRYPTED_API_KEY = "AV6+..." . # from `pixlet encyrpt`

def main(config):
    api_key = secret.decrypt(ENCRYPTED_API_KEY) or config.get("dev_api_key")
```

## Pixlet module: Sunrise

The `sunrise` module calculates sunrise and sunset times for a given set of GPS coordinates and timestamp. 

| Function | Description |
| --- | --- |
| `sunrise(lat, lng, date)` | Calculates the sunrise time for a given location and date. |
| `sunset(lat, lng, date)` | Calculates the sunset time for a given location and date. |
| `elevation(lat, lng, time)` | Calculates the elevation of the sun above the horizon for a given location and point in time. |
| `elevation_time(lat, lng, elev, date)` | Calculates the two times at which the sun was at the given elevation above the horizon for a given location and date. Returns None if the sun never reached the given elevation. |

Example:

See [examples/sunrise.star](../examples/sunrise.star) for an example.

## Pixlet module: Random

The `random` module provides a pseudorandom number generator for pixlet. The generator is automatically seeded on each execution. The seed itself changes every 15 seconds, making apps deterministic over that same time window. This behavior enables more effective caching of execution results on Tidbyt servers. Developer can reseed via `random.seed` if needed.

| Function | Description |
| --- | --- |
| `seed(s)` | Seeds the generator.|
| `number(min, max)` | Returns a random number between the min and max. The min has to be 0 or greater. The min has to be less than the max. |

Example:
```starlark
load("random.star", "random")

def main(config):
    num = random.number(0, 100)
    if num > 50:
        print("You win!")
    else:
        print("Better luck next time!")
```

## Pixlet module: QRCode

The `qrcode` module provides a QR code generator for pixlet!

| Function | Description |
| --- | --- |
| `generate(url, size, color?, background?)` | Returns a QR code as an image that can be passed into the image widget. |

Sizing works as follows:
- `small`: 21x21 pixels
- `medium`: 25x25 pixels
- `large`: 29x29 pixels

Note: we're working with some of the smallest possible QR codes in this module, so the amount of data that can be used for the URL is extremely limited.

Example:
```starlark
load("cache.star", "cache")
load("encoding/base64.star", "base64")
load("render.star", "render")
load("qrcode.star", "qrcode")

def main(config):
    url = "https://tidbyt.com?utm_source=pixlet_example"

    data = cache.get(url)
    if data == None:
        code = qrcode.generate(
            url = url,
            size = "large",
            color = "#fff",
            background = "#000",
        )
        cache.set(url, base64.encode(code), ttl_seconds = 3600)
    else:
        code = base64.decode(data)

    return render.Root(
        child = render.Padding(
            child = render.Image(src = code),
            pad = 1,
        ),
    )
```
