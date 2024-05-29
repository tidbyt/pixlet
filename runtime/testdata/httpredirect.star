load("assert.star", "assert")
load("http.star", "http")
load("render.star", "render")

def main():
    res_1 = http.post(test_server_url + "/login")
    assert.eq(res_1.status_code, 200)
    assert.eq(res_1.body(), '{"hello":"world"}')
    assert.eq(res_1.json(), {"hello": "world"})
    return render.Root(child = render.Text("pass"))
