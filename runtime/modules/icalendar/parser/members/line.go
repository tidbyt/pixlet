package members

import "strings"

func ParseRecurrenceParams(p string) (string, map[string]string) {
	tokens := strings.Split(p, ";")

	parameters := make(map[string]string)
	for _, p = range tokens {
		t := strings.Split(p, "=")
		if len(t) != 2 {
			continue
		}
		parameters[t[0]] = t[1]
	}

	return tokens[0], parameters
}

func ParseParameters(p string) (string, map[string]string) {
	tokens := strings.Split(p, ";")

	parameters := make(map[string]string)

	for _, p = range tokens[1:] {
		t := strings.Split(p, "=")
		if len(t) != 2 {
			continue
		}

		parameters[t[0]] = t[1]
	}

	return tokens[0], parameters
}

func UnescapeString(l string) string {
	l = strings.Replace(l, `\\`, `\`, -1)
	l = strings.Replace(l, `\;`, `;`, -1)
	l = strings.Replace(l, `\,`, `,`, -1)

	return l
}
