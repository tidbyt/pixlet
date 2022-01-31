package schema_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"tidbyt.dev/pixlet/runtime"
	"tidbyt.dev/pixlet/schema"
)

func loadApp(code string) (*runtime.Applet, error) {
	app := &runtime.Applet{}
	err := app.Load("test.star", []byte(code), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return app, nil
}

// appruntime.Schema test with all available config types and flags.
func TestSchemaAllTypesSuccess(t *testing.T) {
	code := `
def get_schema():
    return [
        {"type": "location",
         "name": "Location",
         "description": "A Location",
         "id": "locationid",
         "icon": "place",
        },
        {"type": "locationbased",
         "id": "locationbasedid",
         "name": "Locationbased",
         "description": "A Locationbased",
         "handler": "locationbasedhandler",
         "icon": "place",
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

	jsonSchema := app.GetSchema()

	var s schema.Schema
	json.Unmarshal([]byte(jsonSchema), &s)

	assert.Equal(t, schema.Schema{
		Version: "1",
		Fields: []schema.SchemaField{
			{
				Type:        "location",
				ID:          "locationid",
				Name:        "Location",
				Description: "A Location",
				Icon:        "place",
			},
			{
				Type:        "locationbased",
				ID:          "locationbasedid",
				Name:        "Locationbased",
				Description: "A Locationbased",
				Handler:     "locationbasedhandler",
				Icon:        "place",
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

	jsonSchema, err := app.CallSchemaHandler(context.Background(), "generatedid", "foobar")
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

// Verifies that appruntime.Schemas returned by a generated field's handler is
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
    return [{"type": "text",
            # this missing field is required:  "id": "generatedid",
             "name": "Generated Text",
             "description": "This Text is generated",
             "default": "-%s-" % param,
            }]

def main():
    return None
`

	app, err := loadApp(code)
	assert.NoError(t, err)

	_, err = app.CallSchemaHandler(context.Background(), "generatedid", "foobar")
	assert.Error(t, err)
}

// Verifies that appruntime.Schemas returned by a generated field's handler is
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

	stringValue, err := app.CallSchemaHandler(context.Background(), "locationbasedid", "fart")
	assert.NoError(t, err)
	assert.Equal(t, "[{\"text\":\"Your only option is\",\"value\":\"fart\"}]", stringValue)
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

	_, err = app.CallSchemaHandler(context.Background(), "locationbasedid", "fart")
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

	stringValue, err := app.CallSchemaHandler(context.Background(), "typeaheadid", "farts")
	assert.NoError(t, err)
	assert.Equal(t, "[{\"text\":\"You searched for\",\"value\":\"farts\"}]", stringValue)
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

	_, err = app.CallSchemaHandler(context.Background(), "typeaheadid", "fart")
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

	stringValue, err := app.CallSchemaHandler(context.Background(), "oauth2id", "farts")
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

	_, err = app.CallSchemaHandler(context.Background(), "oauth2id", "farts")
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
