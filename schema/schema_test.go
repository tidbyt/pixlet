package schema_test

import (
	"context"
	"encoding/json"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tidbyt.dev/pixlet/runtime"
	"tidbyt.dev/pixlet/schema"
)

func loadApp(code string) (*runtime.Applet, error) {
	vfs := fstest.MapFS{
		"test.star": &fstest.MapFile{Data: []byte(code)},
		"ding.mp3":  &fstest.MapFile{Data: []byte("ding data")},
	}
	return runtime.NewAppletFromFS("test", vfs)
}

func TestSchemaAllTypesSuccess(t *testing.T) {
	code := `
load("schema.star", "schema")
load("ding.mp3", ding = "file")

# these won't be called unless GetSchemaHandler() is
def locationbasedhandler():
    return None

def generatedhandler():
    return None

def typeaheadhandler():
    return ":)"

def oauth2handler():
    return "a-refresh-token"

def get_schema():
    return schema.Schema(
        version = "1",

        notifications = [
            schema.Notification(
                id = "notificationid",
                name = "Notification",
                desc = "A Notification",
                icon = "notification",
                sounds = [
                    schema.Sound(
                        id = "ding",
                        title = "Ding!",
                        file = ding,
                    ),
                ],
            ),
        ],

        fields = [
            schema.Location(
                id = "locationid",
                name = "Location",
                desc = "A Location",
                icon = "locationDot",
            ),
            schema.LocationBased(
                id = "locationbasedid",
                name = "Locationbased",
                desc = "A Locationbased",
                icon = "locationDot",
                handler = locationbasedhandler,
            ),
            schema.Toggle(
                id = "onoffid",
                name = "On or off",
                desc = "An Onoff",
                icon = "schedule",
                default = False,
            ),
            schema.Text(
                id = "textid",
                name = "Text",
                desc = "A Text",
                icon = "gear",
                default = "Default text",
            ),
            schema.Dropdown(
                id = "dropdownid",
                name = "Dropdown",
                desc = "A Dropdown",
                icon = "iconthatdoesntexist",
                options = [
                    schema.Option(
                        display = "dt1",
                        value = "dv1",
                    ),
                    schema.Option(
                        display = "dt2",
                        value = "dv2",
                    ),
                ],
                default = "dv2",
            ),
            schema.Typeahead(
                id = "typeaheadid",
                name = "Typeahead",
                desc = "A Typeahead",
                icon = "train",
                handler = typeaheadhandler,
            ),
            schema.OAuth2(
                id = "oauth2id",
                name = "OAuth2",
                desc = "Authentication",
                icon = "train",
                handler = oauth2handler,
                client_id = "oauth2_clientid",
                authorization_endpoint = "https://example.com/auth",
                scopes = [
                    "foo",
                    "bar",
                ],
            ),
            schema.PhotoSelect(
            	id = "pngid",
            	name = "Photo",
            	desc = "Picture",
            	icon = "photo_camera",
            ),
            schema.Color(
                id = "colorid",
                name = "Color",
                desc = "A Color",
                icon = "brush",
                default = "ffaa66",
                palette = ["#ffaa66", "#bbb"],
            ),
        ],
    )

def main():
    return None
`

	app, err := loadApp(code)
	assert.NoError(t, err)

	var s schema.Schema
	json.Unmarshal(app.SchemaJSON, &s)

	assert.Equal(t, schema.Schema{
		Version: "1",

		Notifications: []schema.SchemaField{
			{
				Type:        "notification",
				ID:          "notificationid",
				Name:        "Notification",
				Description: "A Notification",
				Icon:        "notification",
				Sounds: []schema.SchemaSound{
					{
						ID:    "ding",
						Title: "Ding!",
						Path:  "ding.mp3",
					},
				},
			},
		},

		Fields: []schema.SchemaField{
			{
				Type:        "location",
				ID:          "locationid",
				Name:        "Location",
				Description: "A Location",
				Icon:        "locationDot",
			},
			{
				Type:        "locationbased",
				ID:          "locationbasedid",
				Name:        "Locationbased",
				Description: "A Locationbased",
				Handler:     "locationbasedhandler",
				Icon:        "locationDot",
			},
			{
				Type:        "onoff",
				ID:          "onoffid",
				Name:        "On or off",
				Description: "An Onoff",
				Default:     "false",
				Icon:        "schedule",
			},
			{
				Type:        "text",
				ID:          "textid",
				Name:        "Text",
				Description: "A Text",
				Icon:        "gear",
				Default:     "Default text",
			},
			{
				Type:        "dropdown",
				ID:          "dropdownid",
				Name:        "Dropdown",
				Description: "A Dropdown",
				Options: []schema.SchemaOption{
					{
						Display: "dt1",
						Text:    "dt1",
						Value:   "dv1",
					},
					{
						Display: "dt2",
						Text:    "dt2",
						Value:   "dv2",
					},
				},
				Default: "dv2",
				Icon:    "iconthatdoesntexist",
			},
			//{
			//	Type:        "text",
			//	ID:          "invisibletext",
			//	Name:        "Invisible Text",
			//	Description: "Conditionally visible text",
			//	Visibility: &schema.SchemaVisibility{
			//		Type:      "invisible",
			//		Condition: "equal",
			//		Variable:  "radio",
			//		Value:     "rv2",
			//	},
			//},
			//{
			//	Type:        "text",
			//	ID:          "invisibletext2",
			//	Name:        "Invisible Text",
			//	Description: "Conditionally visible text",
			//	Visibility: &schema.SchemaVisibility{
			//		Type:      "invisible",
			//		Condition: "not_equal",
			//		Variable:  "radio",
			//		Value:     "rv2",
			//	},
			//},
			//{
			//	Type:    "generated",
			//	ID:      "generatedid",
			//	Handler: "generatedhandler",
			//	Source:  "radioid",
			//},
			{
				Type:        "typeahead",
				ID:          "typeaheadid",
				Name:        "Typeahead",
				Description: "A Typeahead",
				Handler:     "typeaheadhandler",
				Icon:        "train",
			},
			{
				Type:                  "oauth2",
				ID:                    "oauth2id",
				Name:                  "OAuth2",
				Description:           "Authentication",
				Handler:               "oauth2handler",
				Icon:                  "train",
				ClientID:              "oauth2_clientid",
				AuthorizationEndpoint: "https://example.com/auth",
				Scopes:                []string{"foo", "bar"},
			},
			{
				Type:        "png",
				ID:          "pngid",
				Name:        "Photo",
				Description: "Picture",
				Icon:        "photo_camera",
			},
			{
				Type:        "color",
				ID:          "colorid",
				Name:        "Color",
				Description: "A Color",
				Icon:        "brush",
				Default:     "#ffaa66",
				Palette:     []string{"#ffaa66", "#bbb"},
			},
		},
	}, s)
}

func TestSchemaWithNotificationInFields(t *testing.T) {
	code := `
load("schema.star", "schema")
load("ding.mp3", ding = "file")

def get_schema():
    return schema.Schema(
        version = "1",

        fields = [
            schema.Notification(
                id = "notificationid",
                name = "Notification",
                desc = "A Notification",
                icon = "notification",
                sounds = [
                    schema.Sound(
                        id = "ding",
                        title = "Ding!",
                        file = ding,
                    ),
                ],
            ),
        ],
    )

def main():
    return None
`

	_, err := loadApp(code)
	assert.ErrorContains(t, err, "expected fields")
}

func TestSchemaWithFieldsInNotifications(t *testing.T) {
	code := `
load("schema.star", "schema")
load("ding.mp3", ding = "file")

def get_schema():
    return schema.Schema(
        version = "1",

        notifications = [
            schema.Color(
                id = "colorid",
                name = "Color",
                desc = "A Color",
                icon = "brush",
                default = "ffaa66",
            ),
        ],
    )

def main():
    return None
`

	_, err := loadApp(code)
	assert.ErrorContains(t, err, "expected notifications")
}

// test with all available config types and flags.
func TestSchemaAllTypesSuccessLegacy(t *testing.T) {
	code := `
def get_schema():
    return [
        {"type": "location",
         "name": "Location",
         "description": "A Location",
         "id": "locationid",
         "icon": "locationDot",
        },
        {"type": "locationbased",
         "id": "locationbasedid",
         "name": "Locationbased",
         "description": "A Locationbased",
         "handler": "locationbasedhandler",
         "icon": "locationDot",
        },
        {"type": "onoff",
          "id": "onoffid",
          "name": "On or off",
          "description":"An Onoff",
          "default": "false",
          "icon": "schedule",
        },
        {"type": "text",
         "id": "textid",
         "name": "Text",
         "description": "A Text",
         "default": "Default text",
        },
        {"type": "dropdown",
         "id": "dropdownid",
         "name": "Dropdown",
         "icon": "iconthatdoesntexist",
         "description": "A Dropdown",
         "options": [{"text": "dt1", "value": "dv1"},
                     {"text": "dt2", "value": "dv2"}],
         "default": "dv2",
        },
        {"type": "radio",
         "id": "radioid",
         "name": "Radio",
         "description": "A Radio",
         "options": [{"text": "rt1", "value": "rv1"},
                     {"text": "rt2", "value": "rv2"}],
         "default": "rv1",
        },
        {"type": "text",
         "id": "invisibletext",
         "name": "Invisible Text",
         "description": "Conditionally visible text",
         "visibility": {"type": "invisible",
                        "condition": "equal",
                        "variable": "radio",
                        "value": "rv2",
                       },
        },
        {"type": "text",
         "id": "invisibletext2",
         "name": "Invisible Text",
         "description": "Conditionally visible text",
         "visibility": {"type": "invisible",
                        "condition": "not_equal",
                        "variable": "radio",
                        "value": "rv2",
                       },
        },
        {"type": "generated",
         "id": "generatedid",
         "source": "radioid",
         "handler": "generatedhandler",
        },
        {"type": "typeahead",
         "id": "typeaheadid",
         "icon": "train",
         "name": "Typeahead",
         "description": "A Typeahead",
         "handler": "typeaheadhandler",
        },
        {"type": "oauth2",
         "id": "oauth2id",
         "icon": "train",
         "name": "OAuth2",
         "description": "Authentication",
         "handler": "oauth2handler",
         "client_id": "oauth2_clientid",
         "authorization_endpoint": "https://example.com/auth",
         "scopes": ["foo", "bar"],
        },
        {"type": "png",
         "id": "pngid",
         "icon": "photo_camera",
         "name": "Photo",
         "description": "Picture",
        },
        {"type": "color",
         "id": "colorid",
         "icon": "brush",
         "name": "Color",
         "description": "A Color",
         "default": "ffaa66",
         "palette": ["#ffaa66", "bbb"],
        },
    ]

# these won't be called unless GetSchemaHandler() is
def locationbasedhandler():
    return None

def generatedhandler():
    return None

def typeaheadhandler():
    return ":)"

def oauth2handler():
    return "a-refresh-token"

def main():
    return None
`

	app, err := loadApp(code)
	assert.NoError(t, err)

	var s schema.Schema
	json.Unmarshal(app.SchemaJSON, &s)

	assert.Equal(t, schema.Schema{
		Version: "1",
		Fields: []schema.SchemaField{
			{
				Type:        "location",
				ID:          "locationid",
				Name:        "Location",
				Description: "A Location",
				Icon:        "locationDot",
			},
			{
				Type:        "locationbased",
				ID:          "locationbasedid",
				Name:        "Locationbased",
				Description: "A Locationbased",
				Handler:     "locationbasedhandler",
				Icon:        "locationDot",
			},
			{
				Type:        "onoff",
				ID:          "onoffid",
				Name:        "On or off",
				Description: "An Onoff",
				Default:     "false",
				Icon:        "schedule",
			},
			{
				Type:        "text",
				ID:          "textid",
				Name:        "Text",
				Description: "A Text",
				Default:     "Default text",
			},
			{
				Type:        "dropdown",
				ID:          "dropdownid",
				Name:        "Dropdown",
				Description: "A Dropdown",
				Options: []schema.SchemaOption{
					{
						Text:  "dt1",
						Value: "dv1",
					},
					{
						Text:  "dt2",
						Value: "dv2",
					},
				},
				Default: "dv2",
				Icon:    "iconthatdoesntexist",
			},
			{
				Type:        "radio",
				ID:          "radioid",
				Name:        "Radio",
				Description: "A Radio",
				Options: []schema.SchemaOption{
					{
						Text:  "rt1",
						Value: "rv1",
					},
					{
						Text:  "rt2",
						Value: "rv2",
					},
				},
				Default: "rv1",
			},
			{
				Type:        "text",
				ID:          "invisibletext",
				Name:        "Invisible Text",
				Description: "Conditionally visible text",
				Visibility: &schema.SchemaVisibility{
					Type:      "invisible",
					Condition: "equal",
					Variable:  "radio",
					Value:     "rv2",
				},
			},
			{
				Type:        "text",
				ID:          "invisibletext2",
				Name:        "Invisible Text",
				Description: "Conditionally visible text",
				Visibility: &schema.SchemaVisibility{
					Type:      "invisible",
					Condition: "not_equal",
					Variable:  "radio",
					Value:     "rv2",
				},
			},
			{
				Type:    "generated",
				ID:      "generatedid",
				Handler: "generatedhandler",
				Source:  "radioid",
			},
			{
				Type:        "typeahead",
				ID:          "typeaheadid",
				Name:        "Typeahead",
				Description: "A Typeahead",
				Handler:     "typeaheadhandler",
				Icon:        "train",
			},
			{
				Type:                  "oauth2",
				ID:                    "oauth2id",
				Name:                  "OAuth2",
				Description:           "Authentication",
				Handler:               "oauth2handler",
				Icon:                  "train",
				ClientID:              "oauth2_clientid",
				AuthorizationEndpoint: "https://example.com/auth",
				Scopes:                []string{"foo", "bar"},
			},
			{
				Type:        "png",
				ID:          "pngid",
				Name:        "Photo",
				Description: "Picture",
				Icon:        "photo_camera",
			},
			{
				Type:        "color",
				ID:          "colorid",
				Name:        "Color",
				Description: "A Color",
				Icon:        "brush",
				Default:     "ffaa66",
				Palette:     []string{"#ffaa66", "bbb"},
			},
		},
	}, s)
}

func TestSchemaWithGeneratedFieldSuccess(t *testing.T) {
	code := `
def get_schema():
    return [
        {"type": "location",
         "id": "loc",
         "name": "Location",
         "description": "A Location",
        },
        {"id": "generatedid",
         "type": "generated",
         "source": "loc",
         "handler": "generate_schema",
        },
    ]

def generate_schema(param):
    return [{"type": "text",
             "id": "generatedid",
             "name": "Generated Text",
             "description": "This Text is generated",
             "default": "-%s-" % param,
            }]

def main():
    return None
`

	app, err := loadApp(code)
	assert.NoError(t, err)

	// Handlers are not identified by ID
	_, err = app.CallSchemaHandler(context.Background(), "generatedid", "foobar")
	assert.Error(t, err)

	// They're identified by function name
	jsonSchema, err := app.CallSchemaHandler(context.Background(), "generate_schema", "foobar")
	assert.NoError(t, err)

	var s schema.Schema
	json.Unmarshal([]byte(jsonSchema), &s)

	assert.Equal(t, schema.Schema{
		Version: "1",
		Fields: []schema.SchemaField{
			{
				Type:        "text",
				ID:          "generatedid",
				Name:        "Generated Text",
				Description: "This Text is generated",
				Default:     "-foobar-",
			},
		},
	}, s)
}

// Verifies that schema returned by a generated field's handler is
// validated
func TestSchemaWithGeneratedHandlerMalformed(t *testing.T) {
	code := `
def get_schema():
    return [
        {"type": "location",
         "id": "loc",
         "name": "Location",
         "description": "A Location",
        },
        {"id": "generatedid",
         "type": "generated",
         "source": "loc",
         "handler": "generate_schema",
        },
    ]

def generate_schema(param):
    field = {"type": "text",
             "name": "Generated Text",
             "description": "This Text is generated",
             "default": "-%s-" % param,
            }
    if param == "win":
        # this missing field is required
        field["id"] = "generatedid"
    return [field]

def main():
    return None
`

	app, err := loadApp(code)
	assert.NoError(t, err)

	_, err = app.CallSchemaHandler(context.Background(), "generate_schema", "win")
	assert.NoError(t, err)

	_, err = app.CallSchemaHandler(context.Background(), "generate_schema", "fail")
	assert.Error(t, err)
}

// Verifies that schema returned by a generated field's handler is
// validated
func TestSchemaWithGeneratedHandlerMissing(t *testing.T) {
	code := `
def get_schema():
    return [
        {"type": "location",
         "id": "loc",
         "name": "Location",
         "description": "A Location",
        },
        {"id": "generatedid",
         "type": "generated",
         "source": "loc",
         "handler": "this_handler_doesnt_exist",
        },
    ]

def main():
    return None
`

	_, err := loadApp(code)
	assert.Error(t, err)
}

func TestSchemaWithGeneratedIconNotAllowed(t *testing.T) {
	code := `
def get_schema():
    return [
        {"type": "location",
         "id": "loc",
         "name": "Location",
         "description": "A Location",
        },
        {"id": "generatedid",
         "icon": "schedule",
         "type": "generated",
         "source": "loc",
         "handler": "generate_schema",
        },
    ]

def generate_schema(param):
    return [{"type": "text",
             "id": "generatedid",
             "name": "Generated Text",
             "description": "This Text is generated",
             "default": "-%s-" % param,
            }]

def main():
    return None
`

	_, err := loadApp(code)
	assert.Error(t, err)
}

func TestSchemaWithLocationBasedHandlerSuccess(t *testing.T) {
	code := `

def get_schema():
    return [
        {"type": "locationbased",
         "id": "locationbasedid",
         "name": "Locationbased",
         "description": "A Locationbased",
         "handler": "handle_location",
        },
    ]

def handle_location(location):
    return [{"text": "Your only option is", "value": location}]

def main():
    return None
`

	app, err := loadApp(code)
	assert.NoError(t, err)

	stringValue, err := app.CallSchemaHandler(context.Background(), "handle_location", "fart")
	assert.NoError(t, err)
	assert.Equal(t, "[{\"display\":\"\",\"text\":\"Your only option is\",\"value\":\"fart\"}]", stringValue)
}

func TestSchemaWithLocationBasedHandlerMalformed(t *testing.T) {
	code := `

def get_schema():
    return [
        {"type": "locationbased",
         "id": "locationbasedid",
         "name": "Locationbased",
         "description": "A Locationbased",
         "handler": "handle_location",
        },
    ]

def handle_location(location):
    return [{"text": "this option has no value"}]

def main():
    return None
`

	app, err := loadApp(code)
	assert.NoError(t, err)

	_, err = app.CallSchemaHandler(context.Background(), "handle_location", "fart")
	assert.Error(t, err)
}

func TestSchemaWithTypeaheadHandlerSuccess(t *testing.T) {
	code := `

def get_schema():
    return [
        {"type": "typeahead",
         "id": "typeaheadid",
         "name": "Typeahead",
         "description": "A Typeahead",
         "handler": "handle_typeahead",
        },
    ]

def handle_typeahead(pattern):
    return [{"text": "You searched for", "value": pattern}]

def main():
    return None
`

	app, err := loadApp(code)
	assert.NoError(t, err)

	stringValue, err := app.CallSchemaHandler(context.Background(), "handle_typeahead", "farts")
	assert.NoError(t, err)
	assert.Equal(t, "[{\"display\":\"\",\"text\":\"You searched for\",\"value\":\"farts\"}]", stringValue)
}

func TestSchemaWithTypeaheadHandlerMalformed(t *testing.T) {
	code := `

def get_schema():
    return [
        {"type": "typeahead",
         "id": "typeaheadid",
         "name": "Typeahead",
         "description": "A Typeahead",
         "handler": "handle_typeahead",
        },
    ]

def handle_typeahead(pattern):
    return [{"value": "this option has not text field"}]

def main():
    return None
`

	app, err := loadApp(code)
	assert.NoError(t, err)

	_, err = app.CallSchemaHandler(context.Background(), "handle_typeahead", "fart")
	assert.Error(t, err)
}

func TestSchemaWithOAuth2HandlerSuccess(t *testing.T) {
	code := `

def get_schema():
    return [
        {"type": "oauth2",
         "id": "oauth2id",
         "icon": "train",
         "name": "OAuth2",
         "description": "Authentication",
         "handler": "oauth2handler",
         "client_id": "oauth2_clientid",
         "authorization_endpoint": "https://example.com/auth",
         "scopes": ["foo", "bar"],
        },
    ]

def oauth2handler(params):
    return "a-refresh-token"

def main():
    return None
`

	app, err := loadApp(code)
	assert.NoError(t, err)

	stringValue, err := app.CallSchemaHandler(context.Background(), "oauth2handler", "farts")
	assert.NoError(t, err)
	assert.Equal(t, "a-refresh-token", stringValue)
}

func TestSchemaWithOAuth2HandlerMalformed(t *testing.T) {
	code := `

def get_schema():
    return [
        {"type": "oauth2",
         "id": "oauth2id",
         "icon": "train",
         "name": "OAuth2",
         "description": "Authentication",
         "handler": "oauth2handler",
         "client_id": "oauth2_clientid",
         "authorization_endpoint": "https://example.com/auth",
         "scopes": ["foo", "bar"],
        },
    ]

def oauth2handler(params):
    return 123

def main():
    return None
`

	app, err := loadApp(code)
	assert.NoError(t, err)

	_, err = app.CallSchemaHandler(context.Background(), "oauth2handler", "farts")
	assert.Error(t, err)
}

func TestEmptySchemaSerialization(t *testing.T) {
	s := &schema.Schema{
		Version: "1",
	}

	ser, err := json.Marshal(s)
	require.NoError(t, err)
	assert.Equal(t, `{"version":"1","schema":[]}`, string(ser))
}

func TestSchemaExtraHandlers(t *testing.T) {
	code := `
load("schema.star", "schema")

def get_restaurants(param):
    return [schema.Option(display = "McDonalds", value = "mcd")]

def get_somethingelse(param):
    if param == "win":
        return [schema.Option(display = "hey", value = "ho")]
    else:
        return "this handler shouldn't return string"

def not_exposed(param):
    return get_somethingelse("win")

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.LocationBased(
                id = "restaurant",
                name = "Restaurant",
                desc = "Restaurant to track",
                icon = "food",
                handler = get_restaurants,
            ),
        ],
        handlers = [
            schema.Handler(
                handler = get_somethingelse,
                type = schema.HandlerType.Options,
            ),
        ],
    )

def main():
    return None
`

	app, err := loadApp(code)
	assert.NoError(t, err)
	assert.NotNil(t, app)

	data, err := app.CallSchemaHandler(context.Background(), "get_somethingelse", "win")
	assert.NoError(t, err)
	var options []map[string]string
	assert.NoError(t, json.Unmarshal([]byte(data), &options))
	assert.Equal(t, 1, len(options))
	assert.Equal(t, "hey", options[0]["display"])
	assert.Equal(t, "ho", options[0]["value"])

	_, err = app.CallSchemaHandler(context.Background(), "get_somethingelse", "fail")
	assert.Error(t, err)
}

func TestSchemaGeneratedV2OrWhatever(t *testing.T) {
	code := `
load("schema.star", "schema")

def build_boroughs(param):
    if param != "true":
        return []
    return [
        schema.Dropdown(
            id = "borough",
            name = "Borough",
            desc = "Pick a borough!",
            icon = "football",
            options = [
                schema.Option(display = "Brooklyn", value = "BK"),
                schema.Option(display = "Shaolin", value = "WU"),
            ],
            default = "BK",
        ),
    ]

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.Toggle(
                id = "with_borough",
                name = "Limit to borough",
                desc = "Optionally limit app to a certain borough",
                default = False,
                icon = "football",
            ),
            schema.Generated(
                id = "generated shouldnt need id, but still does",
                source = "with_borough",
                handler = build_boroughs,
            ),
        ],
    )

def main():
    return None
`

	app, err := loadApp(code)
	assert.NoError(t, err)
	assert.NotNil(t, app)

	data, err := app.CallSchemaHandler(context.Background(), "build_boroughs", "false")
	assert.NoError(t, err)
	var schema schema.Schema
	assert.NoError(t, json.Unmarshal([]byte(data), &schema))
	assert.Equal(t, "1", schema.Version)
	assert.Equal(t, 0, len(schema.Fields))

	data, err = app.CallSchemaHandler(context.Background(), "build_boroughs", "true")
	assert.NoError(t, err)
	assert.NoError(t, json.Unmarshal([]byte(data), &schema))
	assert.Equal(t, 1, len(schema.Fields))
}

func TestSchemaGeneratedFieldWithHandler(t *testing.T) {
	code := `
load("schema.star", "schema")

def get_station_selector(param):
    if param != "true":
        return []
    return [
        schema.LocationBased(
            id = "station",
            name = "Station",
            desc = "Pick a station!",
            icon = "train",
            handler = get_stations,
        ),
    ]

def get_stations(loc):
    return [
        schema.Option(display="Bedford (L)", value = "L08"),
        schema.Option(display="3rd Ave (L)", value = "3rd"),
    ]

def get_schema():
    return schema.Schema(
        version = "1",
        fields = [
            schema.Toggle(
                id = "select_station",
                name = "Select a station",
                desc = "Optionally select a station",
                default = False,
                icon = "football",
            ),
            schema.Generated(
                id = "generated shouldnt need id, but still does",
                source = "select_station",
                handler = get_station_selector,
            ),
        ],
        # since get_stations isn't referenced in the actual schema (only by
        # the generated field handler) it must be explicitly exported
        handlers = [
            schema.Handler(
                handler = get_stations,
                type = schema.HandlerType.Options,
            ),
        ],
    )

def main():
    return None
`

	app, err := loadApp(code)
	assert.NoError(t, err)
	assert.NotNil(t, app)

	data, err := app.CallSchemaHandler(context.Background(), "get_station_selector", "true")
	assert.NoError(t, err)
	var s schema.Schema
	assert.NoError(t, json.Unmarshal([]byte(data), &s))
	assert.Equal(t, "1", s.Version)
	assert.Equal(t, 1, len(s.Fields))
	assert.Equal(t, "locationbased", s.Fields[0].Type)
	assert.Equal(t, "get_stations", s.Fields[0].Handler)

	data, err = app.CallSchemaHandler(context.Background(), "get_stations", "locationdata")
	var options []schema.SchemaOption
	assert.NoError(t, err)
	assert.NoError(t, json.Unmarshal([]byte(data), &options))
	assert.Equal(t, 2, len(options))
	assert.Equal(t, "L08", options[0].Value)
	assert.Equal(t, "3rd", options[1].Value)
}
