package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lingdor/stackerror"
	"github.com/stoewer/go-strcase"
	// "suntech.com.vn/skylib/skylog.git/skylog"
)

func recoverSetReflectValue(structName, fieldName string, field reflect.Value, value interface{}) {
	if r := recover(); r != nil {
		fmt.Println("recovered from ", r)
		// skylog.Errorf("Field: %v of (%v) - Struct type: %v - DB type: %v", fieldName, structName, field.Type(), reflect.TypeOf(value))
		fmt.Println("Field: %v of (%v) - Struct type: %v - DB type: %v", fieldName, structName, field.Type(), reflect.TypeOf(value))
		e := stackerror.New("=======stackerror=======")
		fmt.Println(e.Error())
	}
}

// SetReflectValue function
func SetReflectValue(structName, fieldName string, field reflect.Value, value interface{}) {
	defer recoverSetReflectValue(structName, fieldName, field, value)

	if !field.IsValid() {
		// skylog.Infof("Field [%v] does not exist", field)
		fmt.Println("Field [%v] does not exist", field)
		return
	}

	if value == nil {
		return
	}

	switch field.Interface().(type) {
	case float64:
		field.Set(reflect.ValueOf(ToF64(value)))
	case *float64:
		field.Set(reflect.ValueOf(ToF64Ptr(value)))
	case int64:
		field.Set(reflect.ValueOf(ToI64(value)))
	case *int64:
		field.Set(reflect.ValueOf(ToI64Ptr(value)))
	case int32:
		field.Set(reflect.ValueOf(ToI32(value)))
	case *int32:
		field.Set(reflect.ValueOf(ToI32Ptr(value)))
	case int16:
		field.Set(reflect.ValueOf(ToI16(value)))
	case *int16:
		field.Set(reflect.ValueOf(ToI16Ptr(value)))
	case string:
		field.Set(reflect.ValueOf(value.(string)))
	case *string:
		field.Set(reflect.ValueOf(AddrOfString(value.(string))))
	case bool:
		field.Set(reflect.ValueOf(value.(bool)))
	case *bool:
		field.Set(reflect.ValueOf(AddrOfBool(value.(bool))))
	case []byte:
		field.Set(reflect.ValueOf(value.([]byte)))
	default:

	}
}

// SetReflectField function
func SetReflectField(st, field reflect.Value, fieldName string, value interface{}) {
	if !field.IsValid() {
		// skylog.Infof("Field [%v] does not exist on struct [%v]", fieldName, st.Type().Name())
		fmt.Println("Field [%v] does not exist on struct [%v]", fieldName, st.Type().Name())
		return
	}
	SetReflectValue(GetStructNameFromValue(st), fieldName, field, value)
}

// GetFieldValueOfStruct function
func GetFieldValueOfStruct(input interface{}, fieldName string) interface{} {
	values := reflect.ValueOf(input)

	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}

	if values.Kind() != reflect.Struct {
		return nil
	}

	field := values.FieldByName(fieldName)

	if !field.IsValid() {
		return nil
	}

	if field.Kind() == reflect.Ptr {
		if !field.IsNil() {
			return reflect.Indirect(field).Interface()
		} else {
			return nil
		}

	} else {
		return field.Interface()
	}
}

// ResetSliceOrStruct function
func ResetSliceOrStruct(source interface{}) {
	if reflect.TypeOf(source).Kind() == reflect.Slice {
		v := reflect.ValueOf(source)
		v.Elem().Set(reflect.MakeSlice(v.Type().Elem(), 0, v.Elem().Cap()))
	} else {
		if reflect.ValueOf(source).Kind() != reflect.Ptr {
			// skylog.Error("Require a pointer parameter")
			fmt.Println("Require a pointer parameter")
			return
		}
		p := reflect.ValueOf(source).Elem()
		p.Set(reflect.Zero(p.Type()))
	}
}

// IsStructOrPtrToStruct function
func IsStructOrPtrToStruct(source interface{}) bool {
	values := reflect.ValueOf(source)
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}

	return values.Kind() == reflect.Struct
}

// IsPtrToStruct function
func IsPtrToStruct(source interface{}) bool {
	values := reflect.ValueOf(source)
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	} else {
		return false
	}

	return values.Kind() == reflect.Struct
}

// IsPtrToStructOrArrayOfStruct functionIsPtrToStructOrArrayOfStruct
func IsPtrToStructOrArrayOfStruct(source interface{}) bool {
	values := reflect.ValueOf(source)
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
		if values.Kind() == reflect.Ptr {
			kind := values.Type().Elem().Kind()
			return kind == reflect.Struct || kind == reflect.Slice
		} else {
			return values.Kind() == reflect.Struct || values.Kind() == reflect.Slice
		}
	} else {
		return false
	}
}

// IsStructOrArrayOfStruct function
func IsStructOrArrayOfStruct(source interface{}) bool {
	values := reflect.ValueOf(source)
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
		return values.Kind() == reflect.Struct || values.Kind() == reflect.Slice
	} else {
		return values.Kind() == reflect.Struct || values.Kind() == reflect.Slice
	}
}

// GetStructNameInSnakeCase function
func GetStructNameInSnakeCase(input interface{}) (string, error) {
	inputValue := reflect.ValueOf(input)

	if inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
	}

	var sliceType reflect.Type
	if inputValue.Kind() == reflect.Slice {
		sliceType = inputValue.Type()
	}

	structType := inputValue.Type()
	if sliceType != nil {
		structType = sliceType.Elem()
	}

	names := strings.Split(structType.String(), ".")
	tableName := names[len(names)-1]
	return strcase.SnakeCase(tableName), nil
}

// GetStructNameFromValue function
func GetStructNameFromValue(inputValue reflect.Value) string {
	if inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
	}

	var sliceType reflect.Type
	if inputValue.Kind() == reflect.Slice {
		sliceType = inputValue.Type()
	}

	structType := inputValue.Type()
	if sliceType != nil {
		structType = sliceType.Elem()
	}

	names := strings.Split(structType.String(), ".")
	return names[len(names)-1]
}

// IsPtrToStructOrArrayOfStruct function
func IsPtrToArrayOfStruct(source interface{}) bool {
	values := reflect.ValueOf(source)
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
		if values.Kind() == reflect.Ptr {
			kind := values.Type().Elem().Kind()
			return kind == reflect.Slice
		} else {
			return values.Kind() == reflect.Slice
		}
	} else {
		return false
	}
}

// GetFieldTagValueOfStruct function
func GetFieldTagValueOfStruct(input interface{}, fieldName, tag string) string {
	// values := reflect.ValueOf(input)
	types := reflect.TypeOf(input)

	if types.Kind() == reflect.Ptr {
		types = types.Elem()
	}

	// if values.Kind() == reflect.Ptr {
	// 	values = values.Elem()
	// }

	if types.Kind() == reflect.Slice {
		types = types.Elem()
	}

	// if values.Kind() == reflect.Slice {
	// 	if values.Len() > 0 {
	// 		values = values.Index(0)
	// 	} else {
	// 		return ""
	// 	}
	// }

	if types.Kind() == reflect.Ptr {
		types = types.Elem()
	}

	// if values.Kind() == reflect.Ptr {
	// 	values = values.Elem()
	// }

	field, found := types.FieldByName(fieldName)
	if !found {
		return ""
	}

	// field, found := values.Type().FieldByName(fieldName)
	// if !found {
	// 	return ""
	// }

	return field.Tag.Get(tag)
}

// Xóa khoảng trắng từ tất cả các trường kiểu string trong struct
func TrimSpaces(s interface{}) {
	valueOfS := reflect.ValueOf(s).Elem() // Lấy giá trị của con trỏ đến struct

	for i := 0; i < valueOfS.NumField(); i++ {
		field := valueOfS.Field(i)
		if field.Kind() == reflect.String {
			// Xóa khoảng trắng từ giá trị của trường
			trimmedValue := strings.TrimSpace(field.String())
			field.SetString(trimmedValue)
		}
	}
}
