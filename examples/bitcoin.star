load("render.star", "render")
load("http.star", "http")
load("encoding/base64.star", "base64")
load("cache_redis.star", "cache_redis")

COINDESK_PRICE_URL = "https://api.coindesk.com/v1/bpi/currentprice.json"

BTC_ICON = base64.decode("""
iVBORw0KGgoAAAANSUhEUgAAABEAAAARCAYAAAA7bUf6AAAAlklEQVQ4T2NkwAH+H2T/jy7FaP+
TEZtyDEG4Zi0TTPXXzoDF0A1DMQRsADbN6MZdO4NiENwQbAbERh1lWLzMmgFGo5iFZBDYEFwuwG
sISCPUIKyGgDRjAyBXYXMNIz5XgDQga8TpLboYgux8DO/AwoUuLiEqTLBFMcmxQ7V0gssgklIsL
AYozjsoBoE45OZi5DRBSnkCAMLhlPBiQGHlAAAAAElFTkSuQmCC
""")

cache_redis.connect("redis-10416.c114.us-east-1-4.ec2.cloud.redislabs.com:10416", "default", "rwE3yDkORKKS2hmaPU3TFittJoQQyPqC", 11389794)

def main():
    rate_cached = cache_redis.get("btc_rate")
    if rate_cached != None:
        print("Hit! Displaying cached data.")
        rate = int(rate_cached)
    else:
        print("Miss! Calling CoinDesk API.")
        rep = http.get(COINDESK_PRICE_URL)
        if rep.status_code != 200:
            fail("Coindesk request failed with status %d", rep.status_code)
        rate = rep.json()["bpi"]["USD"]["rate_float"]
        cache_redis.set("btc_rate", str(int(rate)), ttl_seconds = 240)

    return render.Root(
        child = render.Box(
            render.Row(
                expanded = True,
                main_align = "space_evenly",
                cross_align = "center",
                children = [
                    render.Image(src = BTC_ICON),
                    render.Text("$%d" % rate),
                ],
            ),
        ),
    )
