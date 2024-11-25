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
		fmt.Println("[models] Object Field is:", objField.Kind())

		if objField.CanSet() {
			fmt.Println("[models] Can Set Object Field")
			if objField.Kind() == reflect.Ptr {
				switch inputType := val.(type) {
				case string:
					typedVal, _ := val.(string)
					objField.Set(reflect.ValueOf(&typedVal))

				case int, int8, int16, int32, int64:
					typedVal, _ := val.(int64)
					objField.Set(reflect.ValueOf(&typedVal))

				case float32, float64:
					f64, _ := val.(float64)
					f32 := float32(f64)
					objField.Set(reflect.ValueOf(&f32))

				case bool:
					typedVal, _ := val.(bool)
					objField.Set(reflect.ValueOf(&typedVal))

				case nil:
					objField.Set(reflect.Zero(objField.Type()))

				default:
					fmt.Println("[models] Unknown Pointer Type:", inputType)
				}

			} else {
				switch inputType := val.(type) {
				case string:
					objField.SetString(val.(string))

				case int, int8, int16, int32, int64:
					objField.SetInt(val.(int64))

				case float32, float64:
					objField.SetFloat(val.(float64))

				case bool:
					objField.SetBool(val.(bool))

				default:
					fmt.Println("[models] Unknown Type:", inputType)
				}
			}
		} else {
			fmt.Println("[models] Cannot Set Object Field")
			return object, utils.NewApiError(utils.ErrorType_UpdateError, "Update Error")
		}
	}

	return object, nil
}
