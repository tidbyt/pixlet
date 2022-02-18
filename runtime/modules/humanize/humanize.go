package humanize

import (
	"fmt"
	"math"
	"sync"
	"time"

	gohumanize "github.com/dustin/go-humanize"
	gohumanizeEnglish "github.com/dustin/go-humanize/english"
	startime "go.starlark.net/lib/time"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	ModuleName = "humanize"
)

var (
	once   sync.Once
	module starlark.StringDict
)

func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		module = starlark.StringDict{
			ModuleName: &starlarkstruct.Module{
				Name: ModuleName,
				Members: starlark.StringDict{
					"time":               starlark.NewBuiltin("time", times),
					"relative_time":      starlark.NewBuiltin("relative_time", relativeTime),
					"bytes":              starlark.NewBuiltin("bytes", bytes),
					"parse_bytes":        starlark.NewBuiltin("parse_bytes", parseBytes),
					"comma":              starlark.NewBuiltin("comma", comma),
					"ordinal":            starlark.NewBuiltin("ordinal", ordinal),
					"ftoa":               starlark.NewBuiltin("ftoa", ftoa),
					"plural":             starlark.NewBuiltin("plural", plural),
					"plural_word":        starlark.NewBuiltin("plural_word", pluralWord),
					"word_series":        starlark.NewBuiltin("word_series", wordSeries),
					"oxford_word_series": starlark.NewBuiltin("oxford_word_series", oxfordWordSeries),
				},
			},
		}
	})

	return module, nil
}

func times(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starDate startime.Time
	)

	if err := starlark.UnpackArgs(
		"time",
		args, kwargs,
		"date", &starDate,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for time: %s", err)
	}

	date := time.Time(starDate)
	val := gohumanize.Time(date)

	return starlark.String(val), nil
}

func relativeTime(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starDateA  startime.Time
		starDateB  startime.Time
		starLabelA starlark.String
		starLabelB starlark.String
	)

	if err := starlark.UnpackArgs(
		"relative_time",
		args, kwargs,
		"date_a", &starDateA,
		"date_b", &starDateB,
		"label_a?", &starLabelA,
		"label_b?", &starLabelB,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for time: %s", err)
	}

	dateA := time.Time(starDateA)
	dateB := time.Time(starDateB)
	val := gohumanize.RelTime(dateA, dateB, starLabelA.GoString(), starLabelB.GoString())
	return starlark.String(val), nil
}

func bytes(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starBytes starlark.Int
		starIEC   starlark.Bool
	)

	if err := starlark.UnpackArgs(
		"bytes",
		args, kwargs,
		"size", &starBytes,
		"iec?", &starIEC,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for bytes: %s", err)
	}

	bytes := uint64(starBytes.BigInt().Uint64())
	iec := bool(starIEC)

	var val string
	if iec {
		val = gohumanize.IBytes(bytes)
	} else {
		val = gohumanize.Bytes(bytes)
	}
	return starlark.String(val), nil
}

func parseBytes(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starSize starlark.String
	)

	if err := starlark.UnpackArgs(
		"parse_bytes",
		args, kwargs,
		"size", &starSize,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for bytes: %s", err)
	}

	formatted, err := gohumanize.ParseBytes(starSize.GoString())

	if err != nil {
		return nil, fmt.Errorf("unable to parse bytes: %s: %s", starSize.GoString(), err)
	}

	return starlark.MakeUint64(formatted), nil
}

func comma(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starNum starlark.Value
	)

	if err := starlark.UnpackArgs(
		"comma",
		args, kwargs,
		"num", &starNum,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for comma: %s", err)
	}

	switch starNum := starNum.(type) {
	case starlark.Int:
		i := int64(starNum.BigInt().Int64())
		val := gohumanize.Comma(i)
		return starlark.String(val), nil
	case starlark.Float:
		f := float64(starNum)
		if math.IsInf(f, 0) {
			return nil, fmt.Errorf("cannot convert float infinity to integer")
		} else if math.IsNaN(f) {
			return nil, fmt.Errorf("cannot convert float NaN to integer")
		}
		val := gohumanize.Commaf(f)
		return starlark.String(val), nil
	}
	return nil, fmt.Errorf("cannot convert %s to int or float", starNum.Type())
}

func ordinal(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starNum starlark.Int
	)

	if err := starlark.UnpackArgs(
		"ordinal",
		args, kwargs,
		"num", &starNum,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for ordinal: %s", err)
	}

	num := int(starNum.BigInt().Int64())
	val := gohumanize.Ordinal(num)
	return starlark.String(val), nil
}

func ftoa(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starNum    starlark.Float
		starDigits starlark.Value
	)

	if err := starlark.UnpackArgs(
		"ftoa",
		args, kwargs,
		"num", &starNum,
		"digits?", &starDigits,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for ftoa: %s", err)
	}

	var val string
	num := float64(starNum)

	switch starDigits := starDigits.(type) {
	case starlark.Int:
		digits := int(starDigits.BigInt().Int64())
		val = gohumanize.FtoaWithDigits(num, digits)
	case starlark.Float:
		digits := int(starDigits)
		val = gohumanize.FtoaWithDigits(num, digits)
	}

	if val == "" {
		val = gohumanize.Ftoa(num)
	}
	return starlark.String(val), nil
}

func plural(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starQuantity starlark.Int
		starSingular starlark.String
		starPlural   starlark.String
	)

	if err := starlark.UnpackArgs(
		"plural",
		args, kwargs,
		"quantity", &starQuantity,
		"singular", &starSingular,
		"plural?", &starPlural,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for plural: %s", err)
	}

	val := gohumanizeEnglish.Plural(int(starQuantity.BigInt().Int64()), starSingular.GoString(), starPlural.GoString())
	return starlark.String(val), nil
}

func pluralWord(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starQuantity starlark.Int
		starSingular starlark.String
		starPlural   starlark.String
	)

	if err := starlark.UnpackArgs(
		"plural_word",
		args, kwargs,
		"quantity", &starQuantity,
		"singular", &starSingular,
		"plural?", &starPlural,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for pluralWord: %s", err)
	}

	val := gohumanizeEnglish.PluralWord(int(starQuantity.BigInt().Int64()), starSingular.GoString(), starPlural.GoString())
	return starlark.String(val), nil
}

func getWordList(words *starlark.List) ([]string, error) {
	goList := make([]string, 0, words.Len())
	iter := words.Iterate()
	defer iter.Done()

	var listVal starlark.Value
	for i := 0; iter.Next(&listVal); i++ {
		word, ok := listVal.(starlark.String)
		if !ok {
			return nil, fmt.Errorf(
				"expected data to be a list of String but found %s (at index %d)",
				listVal.Type(),
				i,
			)
		}
		goList = append(goList, word.GoString())
	}
	return goList, nil
}

func wordSeries(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starWords       *starlark.List
		starConjunction starlark.String
	)

	if err := starlark.UnpackArgs(
		"word_series",
		args, kwargs,
		"words", &starWords,
		"conjunction", &starConjunction,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for wordSeries: %s", err)
	}

	words, err := getWordList(starWords)

	if err != nil {
		return nil, fmt.Errorf("failed to get word list: %s", err)
	}

	val := gohumanizeEnglish.WordSeries(words, starConjunction.GoString())
	return starlark.String(val), nil
}

func oxfordWordSeries(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		starWords       *starlark.List
		starConjunction starlark.String
	)

	if err := starlark.UnpackArgs(
		"oxford_word_series",
		args, kwargs,
		"words", &starWords,
		"conjunction", &starConjunction,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for oxfordWordSeries: %s", err)
	}

	words, err := getWordList(starWords)

	if err != nil {
		return nil, fmt.Errorf("failed to get word list: %s", err)
	}

	val := gohumanizeEnglish.OxfordWordSeries(words, starConjunction.GoString())
	return starlark.String(val), nil
}
