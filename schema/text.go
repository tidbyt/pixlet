package schema

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type Text struct {
	SchemaField
}

func newText(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var (
		id   starlark.String
		name starlark.String
		desc starlark.String
		icon starlark.String
		def  starlark.String
	)

	if err := starlark.UnpackArgs(
		"Text",
		args, kwargs,
		"id", &id,
		"name", &name,
		"desc", &desc,
		"icon", &icon,
		"default?", &def,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Text: %s", err)
	}

	s := &Text{}
	s.SchemaField.Type = "text"
	s.ID = id.GoString()
	s.Name = name.GoString()
	s.Description = desc.GoString()
	s.Icon = icon.GoString()
	s.Default = def.GoString()

	return s, nil
}

func (s *Text) AsSchemaField() SchemaField {
	return s.SchemaField
}

func (s *Text) AttrNames() []string {
	return []string{
		"id", "name", "desc", "icon", "default",
	}
}

func (s *Text) Attr(name string) (starlark.Value, error) {
	switch name {

	case "id":
		return starlark.String(s.ID), nil

	case "name":
		return starlark.String(s.Name), nil

	case "desc":
		return starlark.String(s.Description), nil

	case "icon":
		return starlark.String(s.Icon), nil

	case "default":
		return starlark.String(s.Default), nil

	default:
		return nil, nil
	}
}

func (s *Text) String() string       { return "Text(...)" }
func (s *Text) Type() string         { return "Text" }
func (s *Text) Freeze()              {}
func (s *Text) Truth() starlark.Bool { return true }

func (s *Text) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
