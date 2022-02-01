package schema

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type OAuth2 struct {
	SchemaField
	starlarkHandler *starlark.Function
	starlarkScopes  *starlark.List
}

func newOAuth2(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var (
		id           starlark.String
		name         starlark.String
		desc         starlark.String
		icon         starlark.String
		handler      *starlark.Function
		clientID     starlark.String
		authEndpoint starlark.String
		scopes       *starlark.List
	)

	if err := starlark.UnpackArgs(
		"OAuth2",
		args, kwargs,
		"id", &id,
		"name", &name,
		"desc", &desc,
		"icon", &icon,
		"handler", &handler,
		"client_id", &clientID,
		"authorization_endpoint", &authEndpoint,
		"scopes", &scopes,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for OAuth2: %s", err)
	}

	s := &OAuth2{}
	s.SchemaField.Type = "oauth2"
	s.ID = id.GoString()
	s.Name = name.GoString()
	s.Description = desc.GoString()
	s.Icon = icon.GoString()
	s.Handler = handler.Name()
	s.starlarkHandler = handler
	s.ClientID = clientID.GoString()
	s.AuthorizationEndpoint = authEndpoint.GoString()
	s.starlarkScopes = scopes

	if s.starlarkScopes != nil {
		scopesIter := s.starlarkScopes.Iterate()
		defer scopesIter.Done()

		var scopeVal starlark.Value
		for i := 0; scopesIter.Next(&scopeVal); {
			if _, isNone := scopeVal.(starlark.NoneType); isNone {
				continue
			}

			scope, ok := scopeVal.(starlark.String)
			if !ok {
				return nil, fmt.Errorf(
					"expected fields to be a list of string but found: %s (at index %d)",
					scopeVal.Type(),
					i,
				)
			}

			s.Scopes = append(s.Scopes, scope.GoString())
		}
	}

	return s, nil
}

func (s *OAuth2) AsSchemaField() SchemaField {
	return s.SchemaField
}

func (s *OAuth2) AttrNames() []string {
	return []string{
		"id", "name", "desc", "icon", "handler", "client_id", "authorization_endpoint", "scopes",
	}
}

func (s *OAuth2) Attr(name string) (starlark.Value, error) {
	switch name {

	case "id":
		return starlark.String(s.ID), nil

	case "name":
		return starlark.String(s.Name), nil

	case "desc":
		return starlark.String(s.Description), nil

	case "icon":
		return starlark.String(s.Icon), nil

	case "handler":
		return s.starlarkHandler, nil

	case "client_id":
		return starlark.String(s.ClientID), nil

	case "authorization_endpoint":
		return starlark.String(s.AuthorizationEndpoint), nil

	case "scopes":
		return s.starlarkScopes, nil

	default:
		return nil, nil
	}
}

func (s *OAuth2) String() string       { return "OAuth2(...)" }
func (s *OAuth2) Type() string         { return "OAuth2" }
func (s *OAuth2) Freeze()              {}
func (s *OAuth2) Truth() starlark.Bool { return true }

func (s *OAuth2) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
