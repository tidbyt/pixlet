# Authoring Apps
If you haven't already, check out the [tutorial](tutorial.md) for how to get started writing apps. This guide picks up where the tutorial leaves off and provides practices and philosophy on how to build apps using pixlet.

## Architecture
The way pixlet works is heavily influenced by the way Tidbyt works. The Tidbyt displays really tiny 64x32 images and keeps a local cache of images for the installations on the device. On a regular interval, the Tidbyt sends heartbeats to the Tidbyt backend letting the backend know what installations it has inside its cache. The backend then determines what apps need rendered, executes the appropriate Starlark script, and encodes the result as a [WebP](https://developers.google.com/speed/webp) image to send back to the device.

To mimic how we host apps internally, `pixlet render` executes the Starlark script and `pixlet push` pushes the resulting WebP to your Tidbyt.

## Cache
As part of pixlet, we offer the `cache` module to cache results from API requests or other data that's needed between renders. Inside of `pixlet`, it's not super useful given the `pixlet` binary has a short lived execution. When publishing apps using the Tidbyt [Community](https://github.com/tidbyt/community) repo, it can cut down on API requests and add some reliability to your app.

The main thing to keep in mind when working with `cache` is that it's scoped per app. This means that two installations by two different users will have the same cache. Anything that's cached should not be specific to an installation and instead specific to the parameters being requested. Make sure you create a cache key that is unique for the type of information you are caching.

## Config
Config is used to configure an app. It's a key/value pair of config values and is passed into `main`. It's exposed through query parameters in the url of `pixlet serve` or through command line args through `pixlet render`. When publishing apps to the Tidbyt [Community](https://github.com/tidbyt/community) repo, the Tidbyt backend will populate the config values from values provided in the mobile app.

The important thing to remember here is that your app should always be able to render even if a config value isn't provided. Providing default values for every config value or checking it for `None` will ensure the app behaves as expected even if config was not provided. For example, the following ensures there will always be a value for `foo`:

```starlark
DEFAULT_FOO = "bar"

def main(config):
    foo = config.get("foo", DEFAULT_FOO)
```

The `config` object also has helpers to convert config values into specific types:

```starlark
config.str("foo") # returns a string, or None if not found
config.bool("foo") # returns a boolean (True or False), or None if not found
```


## Fail
Using a `fail()` inside of your app should be used incredibly sparingly. It kills the entire execution in an unrecoverable fashion. If your app depends on an external API, it cannot cache the response, and has no means of recovering with a failed response - then a `fail()` is appropriate. Otherwise, using a `print()` and handling the error appropriately is a better approach.