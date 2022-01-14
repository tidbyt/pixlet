package schema

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type Dropdown struct {
	SchemaField
	starlarkOptions *starlark.List
}

func newDropdown(
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
		"Dropdown",
		args, kwargs,
		"id", &id,
		"name", &name,
		"desc", &desc,
		"icon", &icon,
		"default", &def,
		"options", &options,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Dropdown: %s", err)
	}

	s := &Dropdown{}
	s.SchemaField.Type = "dropdown"
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

func (s *Dropdown) AsSchemaField() SchemaField {
	return s.SchemaField
}

func (s *Dropdown) AttrNames() []string {
	return []string{
		"id", "name", "desc", "icon", "default", "options",
	}
}

func (s *Dropdown) Attr(name string) (starlark.Value, error) {
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

func (s *Dropdown) String() string       { return "Dropdown(...)" }
func (s *Dropdown) Type() string         { return "Dropdown" }
func (s *Dropdown) Freeze()              {}
func (s *Dropdown) Truth() starlark.Bool { return true }

func (s *Dropdown) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
