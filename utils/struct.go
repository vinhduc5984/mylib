package utils

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/stoewer/go-strcase"
)

// ProtoToStruct function
func ProtoToStruct(source interface{}, dest interface{}) error {
	jsonOut, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(ConvertKeys(jsonOut, strcase.LowerCamelCase), dest)
}

// StructToProto function
func StructToProto(source interface{}, dest interface{}) error {
	jsonOut, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(ConvertKeys(jsonOut, strcase.SnakeCase), dest)
}

// ConvertKeys function
func ConvertKeys(j json.RawMessage, caseConvert func(string) string) json.RawMessage {
	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal([]byte(j), &m); err != nil {
		var resultArray []map[string]interface{}
		if err := json.Unmarshal([]byte(j), &resultArray); err == nil {
			result := []byte{}
			result = append(result, ([]byte("["))...)
			for index, itm := range resultArray {
				jsonOut, err := json.Marshal(itm)
				if err != nil {
					return j
				}
				result = append(result, ConvertKeys(jsonOut, caseConvert)...)
				if index < len(resultArray)-1 {
					result = append(result, ([]byte(","))...)
				}
			}

			return append(result, ([]byte("]"))...)
		}

		return j
	}

	for k, v := range m {
		fixed := caseConvert(k)

		delete(m, k)
		m[fixed] = ConvertKeys(v, caseConvert)
	}

	b, err := json.Marshal(m)
	if err != nil {
		return j
	}

	return json.RawMessage(b)
}

// GetStructType function
func GetStructType(input interface{}) reflect.Type {
	outValues := reflect.ValueOf(input)

	if reflect.TypeOf(input).Kind() == reflect.Ptr && outValues.Kind() == reflect.Ptr {
		outValues = outValues.Elem()
		if outValues.Kind() == reflect.Ptr {
			return outValues.Type().Elem()
		}
	}

	if reflect.TypeOf(input).Kind() == reflect.Ptr && outValues.Kind() == reflect.Struct { //struct
		return outValues.Type()
	} else if reflect.TypeOf(input).Kind() == reflect.Ptr && outValues.Kind() == reflect.Slice { //array
		sliceType := outValues.Type()
		return sliceType.Elem()
	} else {
		for i := 0; i < outValues.Len(); i++ {
			ele := outValues.Index(i).Elem().Elem()
			if ele.Kind() == reflect.Slice {
				sliceType := ele.Type()
				return sliceType.Elem()
			} else {
				return ele.Type()
			}
		}
	}

	return nil
}

func oneProtoStructConvert(source interface{}, dest interface{}) error {
	sourceValues := reflect.ValueOf(source)
	sourceTypes := reflect.TypeOf(source)
	destValues := reflect.ValueOf(dest)

	sourceStructName := GetStructNameFromValue(sourceValues)

	if sourceValues.Kind() == reflect.Ptr {
		sourceValues = sourceValues.Elem()
		sourceTypes = sourceValues.Type()
	}

	if destValues.Kind() == reflect.Ptr {
		destValues = destValues.Elem()
	}

	for i := 0; i < sourceValues.NumField(); i++ {
		sourceField := sourceValues.Field(i)
		if sourceField.CanInterface() {
			var sourceFieldValue interface{}
			if sourceField.Kind() == reflect.Ptr {
				if sourceField.IsNil() {
					sourceFieldValue = nil
				} else {
					sourceFieldValue = reflect.Indirect(sourceField).Interface()
				}
			} else {
				sourceFieldValue = sourceField.Interface()
			}

			destField := destValues.FieldByName(sourceTypes.Field(i).Name)
			if destField.IsValid() {
				if sourceFieldValue == int(0) || sourceFieldValue == int32(0) || sourceFieldValue == int64(0) || sourceFieldValue == "" {
					SetReflectValue(sourceStructName, sourceTypes.Field(i).Name, destField, nil)
				} else {
					SetReflectValue(sourceStructName, sourceTypes.Field(i).Name, destField, sourceFieldValue)
				}

			}
		}
	}

	return nil
}

func IsTransient(tag string) bool {
	return tag == "true"
}

func oneTransientStructConvert(source interface{}, dest interface{}) error {
	sourceValues := reflect.ValueOf(source)
	sourceTypes := reflect.TypeOf(source)
	destValues := reflect.ValueOf(dest)

	sourceStructName := GetStructNameFromValue(sourceValues)

	if sourceValues.Kind() == reflect.Ptr {
		sourceValues = sourceValues.Elem()
		sourceTypes = sourceValues.Type()
	}

	if destValues.Kind() == reflect.Ptr {
		destValues = destValues.Elem()
	}

	for i := 0; i < sourceValues.NumField(); i++ {
		sourceField := sourceValues.Field(i)

		if sourceField.CanInterface() {
			tags := GetFieldTagValueOfStruct(source, sourceTypes.Field(i).Name, "readonly")

			isTransient := IsTransient(tags)

			if isTransient {
				var sourceFieldValue interface{}
				if sourceField.Kind() == reflect.Ptr {
					if sourceField.IsNil() {
						sourceFieldValue = nil
					} else {
						sourceFieldValue = reflect.Indirect(sourceField).Interface()
					}
				} else {
					sourceFieldValue = sourceField.Interface()
				}

				destField := destValues.FieldByName(sourceTypes.Field(i).Name)
				if destField.IsValid() {
					if sourceFieldValue == int(0) || sourceFieldValue == int32(0) || sourceFieldValue == int64(0) || sourceFieldValue == "" {
						SetReflectValue(sourceStructName, sourceTypes.Field(i).Name, destField, nil)
					} else {
						SetReflectValue(sourceStructName, sourceTypes.Field(i).Name, destField, sourceFieldValue)
					}

				}
			}

		}
	}

	return nil
}

// ProtoStructConvert function
func ProtoStructConvert(source interface{}, dest interface{}) error {
	if !IsStructOrArrayOfStruct(source) {
		return errors.New("Source must be a struct or array of struct")
	}

	if !IsPtrToStructOrArrayOfStruct(dest) {
		return errors.New("Destination must be a pointer to a struct or array of struct")
	}

	sourceValues := reflect.ValueOf(source)
	if sourceValues.Kind() == reflect.Ptr {
		sourceValues = sourceValues.Elem()
		if sourceValues.Kind() == reflect.Ptr {
			sourceValues = sourceValues.Elem()
		}
	}

	destValues := reflect.ValueOf(dest)
	if destValues.Kind() == reflect.Ptr {
		destValues = destValues.Elem()
		if destValues.Kind() == reflect.Ptr {
			destValues = destValues.Elem()
		}

	}
	structType := GetStructType(dest)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()

	}

	if destValues.Kind() == reflect.Struct {
		newOneDest := reflect.New(structType)
		oneProtoStructConvert(source, newOneDest.Interface())
		if destValues.Type().Kind() == reflect.Ptr {
			destValues.Set(newOneDest)
		} else {
			destValues.Set(reflect.Indirect(newOneDest))
		}
	} else if destValues.Kind() == reflect.Slice {
		for i := 0; i < sourceValues.Len(); i++ {
			newOneDest := reflect.New(structType)
			oneProtoStructConvert(sourceValues.Index(i).Interface(), newOneDest.Interface())
			if destValues.Type().Elem().Kind() == reflect.Ptr {
				destValues.Set(reflect.Append(destValues, newOneDest))
			} else {
				destValues.Set(reflect.Append(destValues, reflect.Indirect(newOneDest)))
			}
		}
	}

	return nil
}

// TransientStructConvert function
func TransientStructConvert(source interface{}, dest interface{}) error {
	if !IsStructOrArrayOfStruct(source) {
		return errors.New("Source must be a struct or array of struct")
	}

	if !IsPtrToStructOrArrayOfStruct(dest) {
		return errors.New("Destination must be a pointer to a struct or array of struct")
	}

	sourceValues := reflect.ValueOf(source)
	if sourceValues.Kind() == reflect.Ptr {
		sourceValues = sourceValues.Elem()
		if sourceValues.Kind() == reflect.Ptr {
			sourceValues = sourceValues.Elem()
		}
	}

	destValues := reflect.ValueOf(dest)
	if destValues.Kind() == reflect.Ptr {
		destValues = destValues.Elem()
		if destValues.Kind() == reflect.Ptr {
			destValues = destValues.Elem()
		}

	}
	structType := GetStructType(dest)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	if sourceValues.Kind() == reflect.Struct && destValues.Kind() == reflect.Struct {
		oneTransientStructConvert(source, dest)
	} else if sourceValues.Kind() == reflect.Slice && destValues.Kind() == reflect.Slice {
		for i := 0; i < sourceValues.Len(); i++ {
			oneTransientStructConvert(sourceValues.Index(i).Interface(), destValues.Index(i).Interface())
			destValues.Set(reflect.Append(destValues, destValues.Index(i)))
		}
	}

	return nil
}
