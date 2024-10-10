package members

func ParseRecurrenceRule(v string) (map[string]string, error) {
	_, params := ParseRecurrenceParams(v)

	return params, nil
}
