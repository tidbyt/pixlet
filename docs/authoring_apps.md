# Authoring Apps
If you haven't already, check out the [tutorial](tutorial.md) for how to get started writing apps. This guide picks up where the tutorial leaves off and provides practices and philosophy on how to build apps using pixlet.

## Architecture
Pixlet is heavily influenced by the way Tidbyt devices work. The Tidbyt displays tiny 64x32 animations and images. It keeps a local cache of the apps that are installed on it.

Each Tidbyt regularly sends heartbeats to the Tidbyt cloud, announcing what it has cached. The Tidbyt cloud decides which app needs to be rendered next, and executes the appropriate Starlark script. It encodes the result as a [WebP](https://developers.google.com/speed/webp) image, and sends it back to the device.

To mimic how we host apps internally, `pixlet render` executes the Starlark script and `pixlet push` pushes the resulting WebP to your Tidbyt.

## Config
When running an app, Pixlet passes a `config` object to the app's `main()`:

```starlark
def main(config):
    who = config.get("who")
    print("Hello, %s" % who)
```

The `config` object contains values that are useful for your app. You can set the actual values by:

1. Passing URL query parameters when using `pixlet serve`.
2. Setting command-line arguments via `pixlet render`.

When apps that are published to the [Tidbyt Community repo][3], users can install and configure them with the Tidbyt smartphone app. [Define a schema for your app][4] to enable this.

Your app should always be able to render, even if a config value isn't provided. Provide defaults for every config value, or check if the value is `None`. This will ensure the app behaves as expected even if config was not provided.

For example, the following ensures there will always be a value for `who`:

```starlark
DEFAULT_WHO = "world"

def main(config):
    who = config.get("who") or DEFAULT_WHO
    print("Hello, %s" % who)
```

The `config` object also has helpers to convert config values into specific types:

```starlark
config.str("foo") # returns a string, or None if not found
config.bool("foo") # returns a boolean (True or False), or None if not found
```

## Cache
Use the `cache` module to cache results from API requests or other data that's needed between renders. We require sensible caching for apps in the [Tidbyt Community repo](https://github.com/tidbyt/community). Caching cuts down on API requests, and can make your app more reliable.

**Make sure to create cache keys that are unique for the type of information you are caching.** Each app has its own cache, but that cache is shared among every user. Two installations of the same app by two different users will share the same cache.

A good strategy is to create cache keys based on the config parameters or information being requested.

## Secrets

Many apps need secret values like API keys. When publishing your app to the [Tidbyt community repo][3], encrypt sensitive values so that only the Tidbyt cloud servers can decrypt them.

To encrypt values, use the `pixlet encrypt` command. For example:

```shell
# replace "googletraffic" with the folder name of your app in the community repo
$ pixlet encrypt googletraffic top_secret_google_api_key_123456
"AV6+...."  # encrypted value
```

Use the `secret.decrypt()` function in your app to decrypt this value:

```starlark
load("secret.star", "secret")

def main(config):
    api_key = secret.decrypt("AV6+...") or config.get("dev_api_key")
```

When you run `pixlet` locally, `secret.decrypt` will always return `None`. When your app runs in the Tidbyt cloud, `secret.decrypt` will return the string that you passed to `pixlet encrypt`.


## Fail
The [`fail()`][1] function will immediately end the execution of your app and return an error. It should be used incredibly sparingly, and only in cases that are _permanent_ failures. 

For example, if your app receives an error from an external API, try these options before `fail()`:

1. Return a cached response.
2. Display a useful message or fallback data.
3. [`print()`][2] an error message.
3. Handle the error in a way that makes sense for your app.

[1]: https://github.com/bazelbuild/starlark/blob/master/spec.md#fail
[2]: https://github.com/bazelbuild/starlark/blob/master/spec.md#print
[3]: https://github.com/tidbyt/community
[4]: schema/schema.md
