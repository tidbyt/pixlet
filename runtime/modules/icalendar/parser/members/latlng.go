package members

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseLatLng(l string) (float64, float64, error) {
	token := strings.SplitN(l, ";", 2)
	if len(token) != 2 {
		return 0.0, 0.0, fmt.Errorf("could not parse geo coordinates: %s", l)
	}
	lat, laterr := strconv.ParseFloat(token[0], 64)
	if laterr != nil {
		return 0.0, 0.0, fmt.Errorf("could not parse geo latitude: %s", token[0])
	}
	long, longerr := strconv.ParseFloat(token[1], 64)
	if longerr != nil {
		return 0.0, 0.0, fmt.Errorf("could not parse geo longitude: %s", token[1])
	}

	return lat, long, nil
}
