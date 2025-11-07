package utils

// ToMap function
func ToMap(defaultValue interface{}, slice ...string) map[string]interface{} {
	retMap := make(map[string]interface{})

	for _, v := range slice {
		retMap[v] = defaultValue
	}
	return retMap
}

func FirstIndexOf(sourceSlice []interface{}, testValue interface{}) int {
	for index, v := range sourceSlice {
		if v == testValue {
			return index
		}
	}

	return -1
}
