package animation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type NumberOrPercentage interface {
	Transform(v int) float64
}

type Number struct {
	Value float64
}

func (self Number) Transform(_ int) float64 {
	return self.Value
}

type Percentage struct {
	Value float64
}

func (self Percentage) Transform(v int) float64 {
	return self.Value * float64(v)
}

var percentageRe = regexp.MustCompile(
	`^(?P<percentage>[0-9]+)%$`)

func ParsePercentage(str string, mapping map[string]float64) (Percentage, error) {
	match := percentageRe.FindStringSubmatch(str)
	if match != nil {
		result := make(map[string]string)

		for i, name := range percentageRe.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		if p, err := strconv.ParseInt(result["percentage"], 10, 64); err == nil {
			return Percentage{float64(p) / 100.0}, nil
		} else {
			return Percentage{}, err
		}
	}

	if p, ok := mapping[str]; ok {
		return Percentage{p}, nil
	}

	keys := make([]string, 0, len(mapping))
	for k := range mapping {
		keys = append(keys, k)
	}

	return Percentage{}, fmt.Errorf("invalid string for percentage: %s (expected '%s' or '<number>%%')", str, strings.Join(keys, "', '"))
}
