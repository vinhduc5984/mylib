package utils

import (
	"fmt"
	"reflect"
)

// Print function
func Print(showDetails bool, sources ...interface{}) {
	values := reflect.ValueOf(sources)
	for i := 0; i < values.Len(); i++ {
		ele := values.Index(i).Elem()
		if ele.Kind() == reflect.Ptr {
			ele = ele.Elem()
		}

		if ele.Kind() == reflect.Slice {
			sliceValues := ele
			for i := 0; i < sliceValues.Len(); i++ {
				ele2 := sliceValues.Index(i)
				if ele2.Kind() == reflect.Struct {
					printStruct(showDetails, ele2)
				} else {
					Print(showDetails, ele2)
				}
			}
		} else if ele.Kind() == reflect.Struct {
			printStruct(showDetails, ele)
		} else {
			fmt.Print(ele)
		}
	}
}

func printStruct(showDetails bool, inStruct reflect.Value) {
	for i := 0; i < inStruct.NumField(); i++ {
		field := inStruct.Type().Field(i)
		var value interface{}
		if inStruct.Field(i).Kind() == reflect.Ptr {
			if !inStruct.Field(i).IsNil() {
				if inStruct.Field(i).CanInterface() {
					value = reflect.Indirect(inStruct.Field(i))
				}
			}
		} else {
			if inStruct.Field(i).CanInterface() {
				value = inStruct.Field(i).Interface()
			}

		}

		if showDetails {
			value = fmt.Sprintf("%v(%v: %v)", value, field.Name, field.Type)
		}
		if showDetails {
			fmt.Printf("%30v\t", value)
		} else {
			fmt.Printf("%20v\t", value)
		}
	}
	fmt.Println()
}

func PrintRequest(methodName string, req interface{}) {

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println("Request: %v", methodName)
	if len(fmt.Sprintf("%#v", req)) < 1024 {

		fmt.Println("Param: %v", req)
	}

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	// skylog.Infof(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	// skylog.Infof("Request: %v", methodName)

	// if len(fmt.Sprintf("%#v", req)) < 1024 {
	// 	skylog.Infof("Param: %v", req)
	// }

	// skylog.Infof(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

}
