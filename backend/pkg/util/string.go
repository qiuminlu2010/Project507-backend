package util

import "strconv"

func StringsToInts(strings []string) ([]int, error) {
	var ints []int
	for _, s := range strings {
		d, err := strconv.Atoi(s)
		if err != nil {
			// panic(err)
			return nil, err
		}
		ints = append(ints, d)
	}
	return ints, nil
}
