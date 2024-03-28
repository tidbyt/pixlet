load("http.star", "http")
load("render.star", "render")
load("schema.star", "schema")
load("secret.star", "secret")
load("encoding/json.star", "json")

OAUTH2_CLIENT_SECRET = secret.decrypt("your-client-secret")

EXAMPLE_PARAMS = """
{"code": "your-code", "grant_type": "authorization_code", "client_id": "your-client-id", "redirect_uri": "https://appauth.tidbyt.com/your-app-id"}
"""

def main(config):
    token = config.get("auth")

    if token:
        msg = "Authenticated"
    else:
        msg = "Unauthenticated"

    return render.Root(
        child = render.Marquee(
            width = 64,
            child = render.Text(msg),
        ),
    )

def oauth_handler(params):
    # deserialize oauth2 parameters, see example above.
    params = json.decode(params)

    # exchange parameters and client secret for an access token
    res = http.post(
        url = "https://github.com/login/oauth/access_token",
        headers = {
            "Accept": "application/json",
        },
        form_body = dict(
            params,
            client_secret = OAUTH2_CLIENT_SECRET,
        ),
        form_encoding = "application/x-www-form-urlencoded",
    )
    if res.status_code != 200:
        fail("token request failed with status code: %d - %s" %
             (res.status_code, res.body()))

    token_params = res.json()
    access_token = token_params["access_token"]

    return access_token

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.OAuth2(
                id = "auth",
                name = "GitHub",
                desc = "Connect your GitHub account.",
                icon = "github",
                handler = oauth_handler,
                client_id = "your-client-id",
                authorization_endpoint = "https://github.com/login/oauth/authorize",
                scopes = [
                    "read:user",
                ],
            ),
        ],
    )
