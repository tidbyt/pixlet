"""
Applet: Test App
Summary: For Testing
Description: It's an app for testing.
Author: Test Dev
"""

load("render.star", "render")
load("http.star", "http")
load("assert.star", "assert")

def main(config):
    resp = http.get(
        url = "https://example.com",
        ttl_seconds = 60,
    )
    assert.eq(resp.headers.get("Tidbyt-Cache-Status"), "MISS")

    resp = http.get(
        url = "https://example.com",
        ttl_seconds = 60,
    )
    assert.eq(resp.headers.get("Tidbyt-Cache-Status"), "HIT")

    resp = http.post(
        url = "https://example.com",
        ttl_seconds = 0,
    )
    assert.eq(resp.headers.get("Tidbyt-Cache-Status"), "MISS")

    resp = http.post(
        url = "https://example.com",
        ttl_seconds = 60,
    )
    assert.eq(resp.headers.get("Tidbyt-Cache-Status"), "MISS")

    resp = http.post(
        url = "https://example.com",
        ttl_seconds = 60,
    )
    assert.eq(resp.headers.get("Tidbyt-Cache-Status"), "HIT")

    return render.Root(
        child = render.Box(
            width = 64,
            height = 32,
        ),
    )
