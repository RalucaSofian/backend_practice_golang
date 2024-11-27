package models

import (
	"app/utils"
	"fmt"
	"reflect"
)

type Model interface {
	IsModel()
}

func Update[T Model](object T, input map[string]interface{}) (T, error) {
	// pass object by reference
	objVal := reflect.ValueOf(&object).Elem()

	for key, val := range input {
		// convert key into field name
		inputFieldName := utils.CamelToPascalCase(key)
		objField := objVal.FieldByName(inputFieldName)
		fmt.Println("[models] Object Field is:", objField.Type())

		if objField.CanSet() {
			fmt.Println("[models] Can Set Object Field")
			switch inputType := val.(type) {
			case string:
				typedVal, _ := val.(string)
				if objField.Kind() == reflect.Ptr {
					objField.Set(reflect.ValueOf(&typedVal))
				} else {
					objField.SetString(typedVal)
				}

			case float32, float64:
				f64, _ := val.(float64)
				f32 := float32(f64)
				if objField.Kind() == reflect.Ptr {
					if objField.Type().String() == "*int" {
						intVal := int(f32)
						objField.Set(reflect.ValueOf(&intVal))
					} else {
						objField.Set(reflect.ValueOf(&f32))
					}
				} else {
					if objField.Type().String() == "int" {
						objField.SetInt(int64(f32))
					} else {
						objField.SetFloat(f64)
					}
				}

			case bool:
				typedVal, _ := val.(bool)
				if objField.Kind() == reflect.Ptr {
					objField.Set(reflect.ValueOf(&typedVal))
				} else {
					objField.SetBool(typedVal)
				}

			case nil:
				objField.Set(reflect.Zero(objField.Type()))

			default:
				fmt.Println("[models] Unknown Type:", inputType)
			}

		} else {
			fmt.Println("[models] Cannot Set Object Field")
			return object, utils.NewApiError(utils.ErrorType_UpdateError, "Update Error")
		}
	}

	return object, nil
}
