package schema

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type Radio struct {
	SchemaField
	starlarkOptions *starlark.List
}

func newRadio(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var (
		id      starlark.String
		name    starlark.String
		desc    starlark.String
		icon    starlark.String
		def     starlark.String
		options *starlark.List
	)

	if err := starlark.UnpackArgs(
		"Radio",
		args, kwargs,
		"id", &id,
		"name", &name,
		"desc", &desc,
		"icon", &icon,
		"default", &def,
		"options", &options,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Radio: %s", err)
	}

	s := &Radio{}
	s.SchemaField.Type = "radio"
	s.ID = id.GoString()
	s.Name = name.GoString()
	s.Description = desc.GoString()
	s.Icon = icon.GoString()
	s.Default = def.GoString()

	var optionVal starlark.Value
	optionIter := options.Iterate()
	defer optionIter.Done()
	for i := 0; optionIter.Next(&optionVal); {
		if _, isNone := optionVal.(starlark.NoneType); isNone {
			continue
		}

		o, ok := optionVal.(*Option)
		if !ok {
			return nil, fmt.Errorf(
				"expected options to be a list of Option but found: %s (at index %d)",
				optionVal.Type(),
				i,
			)
		}

		s.Options = append(s.Options, o.SchemaOption)
	}
	s.starlarkOptions = options

	return s, nil
}

func (s *Radio) AsSchemaField() SchemaField {
	return s.SchemaField
}

func (s *Radio) AttrNames() []string {
	return []string{
		"id", "name", "desc", "icon", "default", "options",
	}
}

func (s *Radio) Attr(name string) (starlark.Value, error) {
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

	case "options":
		return s.starlarkOptions, nil

	default:
		return nil, nil
	}
}

func (s *Radio) String() string       { return "Radio(...)" }
func (s *Radio) Type() string         { return "Radio" }
func (s *Radio) Freeze()              {}
func (s *Radio) Truth() starlark.Bool { return true }

func (s *Radio) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
