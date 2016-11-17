package form_validator

import (
	"reflect"
	"bytes"
	"errors"
	"fmt"
	"log"
	"encoding/json"
)

type (
	ValidateableFormObject map[string]interface{}
	ValidateableFormArray []ValidateableFormObject
	Validateable interface {
		Validate(schema reflect.Value, typeErrors *TypeErrors)
	}
	TypeErrors struct {
		errors []error
	}
)

func (t *TypeErrors) add(err error) {
	t.errors = append(t.errors, err)
}

func (t TypeErrors) Error() string {
	var errors string = ""
	for _,err := range t.errors {
		errors += err.Error()
	}
	return errors
}

func (v ValidateableFormObject) Validate(schema reflect.Value, typeErrors *TypeErrors) {
	validateFields(v, schema, typeErrors)

}
func (v ValidateableFormArray) Validate(schema reflect.Value, typeErrors *TypeErrors) {
	for _, value := range v {
		validateFields(value,  reflect.Indirect(reflect.New(schema.Type().Elem())), typeErrors)
	}
}

func ValidateSchema(input []byte, output interface{}) error {
	var typeErrors TypeErrors

	var schemaValue = reflect.Indirect(reflect.ValueOf(output))

	if schemaValue.Kind() != reflect.Struct && schemaValue.Kind() != reflect.Slice {
		return errors.New(fmt.Sprintf("Form validator can only validate structs. Attempted to validate: %v", schemaValue.Kind()))
	}

	var object ValidateableFormObject
	var array ValidateableFormArray

	err := json.Unmarshal(input, &object)

	if err != nil {
		switch err.(type) {
		case *json.UnmarshalTypeError:
			if marshalErrorStruct, ok := err.(*json.UnmarshalTypeError); ok {
				if marshalErrorStruct.Value == "array" {
					err := json.Unmarshal(input, &array)
					if err != nil {
						return err
					}
					array.Validate(schemaValue, &typeErrors)
				}
				return err
			}
			return err
		default:
			return err
		}
	} else {
		object.Validate(schemaValue, &typeErrors)
	}

	if len(typeErrors.errors) == 0 {
		return nil
	}

	return typeErrors

}






func validateFields(fields map[string]interface{}, schema reflect.Value, typeErrors *TypeErrors) {

	for key, value := range fields {
		candidateField := reflect.ValueOf(value)
		schemaField := schema.FieldByNameFunc( func(field string) bool {
			return bytes.EqualFold([]byte(field), []byte(key))
		})

		validationError := validateField(schemaField, candidateField, typeErrors)
		if validationError != nil {
			typeErrors.add(validationError)
		}
	}

}

func validateField(schemaField reflect.Value, candidateField reflect.Value, typeErrors *TypeErrors) error {

	if schemaField.IsValid() {
		log.Printf("Expected: %v - Candidate: %v", schemaField.Type().Kind(), candidateField.Kind())

		switch schemaField.Kind() {

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
			// If interface{} is used, Go parses numeric values to float64.
			if candidateField.Kind() != reflect.Float64 {
				return errorMessage(candidateField, schemaField)
			}

		case reflect.String:
			if candidateField.Kind() != reflect.String {
				return errorMessage(candidateField, schemaField)
			}

		case reflect.Bool:
			if candidateField.Kind() != reflect.Bool {
				return errorMessage(candidateField, schemaField)
			}

		case reflect.Slice: {
			if candidateField.Kind() != reflect.Slice {
				return errorMessage(candidateField, schemaField)
			}
		}
		case reflect.Struct:
			validateFields(candidateField.Interface().(map[string]interface{}), schemaField, typeErrors)

		default:
			return errorMessage(candidateField, schemaField)
		}

	}
	return nil
}

func errorMessage(actual reflect.Value, expected reflect.Value) error {
	return errors.New(fmt.Sprintf(expected.String() + ":" + "Unexpected type: %v - expected: %v", actual.Kind().String(), expected.Kind().String()))
}




