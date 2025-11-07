package utils

import (
	"fmt"
	"strconv"
)

// ToBool convert interface{} to bool.
func ToBool(source interface{}) (bool, error) {
	if source == nil {
		return false, nil
	}

	str := fmt.Sprintf("%v", source)

	res, err := strconv.ParseBool(str)
	if err != nil {
		return false, err
	}
	return res, nil
}
