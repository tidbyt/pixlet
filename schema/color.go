package schema

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type Color struct {
	SchemaField
	starlarkPalette *starlark.List
}

func normalizeHexColor(hex string) (string, error) {
	hex = strings.TrimPrefix(strings.ToLower(hex), "#")

	if len(hex) != 3 && len(hex) != 6 {
		return "", fmt.Errorf(
			"expected 3 or 6 hex chars but found %d",
			len(hex),
		)
	}
	if _, err := strconv.ParseInt(hex, 16, 64); err != nil {
		return "", fmt.Errorf(
			"expected hex chars a-f,0-9 but found %s",
			hex,
		)
	}

	return fmt.Sprintf("#%s", hex), nil
}

func newColor(
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
		palette *starlark.List
	)

	var err error

	if err = starlark.UnpackArgs(
		"Color",
		args, kwargs,
		"id", &id,
		"name", &name,
		"desc", &desc,
		"icon", &icon,
		"default", &def,
		"palette?", &palette,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for Color: %s", err)
	}

	s := &Color{}
	s.SchemaField.Type = "color"
	s.ID = id.GoString()
	s.Name = name.GoString()
	s.Description = desc.GoString()
	s.Icon = icon.GoString()

	s.Default, err = normalizeHexColor(def.GoString())
	if err != nil {
		return nil, fmt.Errorf("malformed default color: %w", err)
	}

	if palette == nil {
		return s, nil
	}

	s.starlarkPalette = starlark.NewList([]starlark.Value{})
	s.Palette = []string{}

	var paletteVal starlark.Value
	paletteIter := palette.Iterate()
	defer paletteIter.Done()
	for i := 0; paletteIter.Next(&paletteVal); {
		if _, isNone := paletteVal.(starlark.NoneType); isNone {
			continue
		}

		col, ok := paletteVal.(starlark.String)
		if !ok {
			return nil, fmt.Errorf(
				"expected palette to be a list of string but found: %s (at index %d)",
				paletteVal.Type(),
				i,
			)
		}

		hex, err := normalizeHexColor(col.GoString())
		if err != nil {
			return nil, fmt.Errorf("malformed palette color at index %d: %w", i, err)
		}

		s.Palette = append(s.Palette, hex)
		s.starlarkPalette.Append(starlark.String(hex))
	}

	return s, nil
}

func (s *Color) AsSchemaField() SchemaField {
	return s.SchemaField
}

func (s *Color) AttrNames() []string {
	return []string{
		"id", "name", "desc", "icon", "default", "palette",
	}
}

func (s *Color) Attr(name string) (starlark.Value, error) {
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

	case "palette":
		return s.starlarkPalette, nil

	default:
		return nil, nil
	}
}

func (s *Color) String() string       { return "Color(...)" }
func (s *Color) Type() string         { return "Color" }
func (s *Color) Freeze()              {}
func (s *Color) Truth() starlark.Bool { return true }

func (s *Color) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
