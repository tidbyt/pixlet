load("http.star", "http")
load("assert.star", "assert")

res_1 = http.get(test_server_url, params = {"a": "b", "c": "d"})
assert.eq(res_1.url, test_server_url + "?a=b&c=d")
assert.eq(res_1.status_code, 200)
assert.eq(res_1.body(), '{"hello":"world"}')
assert.eq(res_1.json(), {"hello": "world"})

assert.eq(res_1.headers, {"Date": "Mon, 01 Jun 2000 00:00:00 GMT", "Content-Length": "17", "Content-Type": "text/plain; charset=utf-8"})

res_2 = http.get(test_server_url)
assert.eq(res_2.json()["hello"], "world")

headers = {"foo": "bar"}
http.post(test_server_url, json_body = {"a": "b", "c": "d"}, headers = headers)
http.post(test_server_url, form_body = {"a": "b", "c": "d"})
