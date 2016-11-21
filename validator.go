package form_validator

import (
	"reflect"
	"bytes"
	"errors"
	"fmt"
	"encoding/json"
)

type (
	ValidateableFormObject map[string]interface{}
	ValidateableFormArray []ValidateableFormObject
	Validateable interface {
		Validate(schema reflect.Value, typeErrors *TypeErrors)
	}
	TypeErrors struct {
		Errors []error `json:"errors"`
	}
	TypeError struct {
		Field string `json:"field"`
		Message string `json:"message"`
	}
)

func (t *TypeErrors) add(err error) {
	t.Errors = append(t.Errors, err)
}

func (t TypeErrors) Error() string {
	jsonString, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(jsonString)
}

func NewTypeError(field, message string) error {
	return &TypeError{Field: field, Message: message}
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("[%v] %v", e.Field, e.Message)
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
		return errors.New(fmt.Sprintf("Form validator can only validate structs and slices"))
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
				} else {
					return errors.New(fmt.Sprintf("Something went wrong parsing input: %v", err.Error()))
				}

			} else {
				return errors.New(fmt.Sprintf("Something went wrong parsing input: %v", err.Error()))
			}
		default:
			return errors.New(fmt.Sprintf("Trying to parse malformed JSON: %v", err.Error()))
		}
	} else {
		object.Validate(schemaValue, &typeErrors)
	}

	if len(typeErrors.Errors) == 0 {
		return nil
	}

	return typeErrors

}

func validateFields(fields map[string]interface{}, schema reflect.Value, typeErrors *TypeErrors) {

	for key, value := range fields {
		candidateField := reflect.ValueOf(value)
		var fieldName string
		schemaField := schema.FieldByNameFunc( func(field string) bool {
			if bytes.EqualFold([]byte(field), []byte(key)) {
				fieldName = field
				return true
			}
			return false
		})

		validationError := validateField(schemaField, candidateField, fieldName, typeErrors)
		if validationError != nil {
			typeErrors.add(validationError)
		}
	}

}

func validateField(schemaField reflect.Value, candidateField reflect.Value, fieldName string, typeErrors *TypeErrors) error {

	// Check if field existed in schema
	if schemaField.IsValid() {
		var k reflect.Kind = schemaField.Kind()
		switch true {
			// json.Unmarshal formats all numeric values to float64
			case 	k == reflect.Int,
				k == reflect.Int8,
				k == reflect.Int16,
				k == reflect.Int32,
				k == reflect.Int64,
				k == reflect.Float32,
				k == reflect.Float64: {
				if candidateField.Kind() != reflect.Float64 {
					return NewTypeError(fieldName, fmt.Sprintf("Unexpected type: %v - expected: %v", candidateField.Kind().String(), k))
				}
			}
			case k == reflect.Struct:
				validateFields(candidateField.Interface().(map[string]interface{}), schemaField, typeErrors)

			case k != candidateField.Kind(): {
				if candidateField.Kind() == reflect.Float64 {
					return NewTypeError(fieldName, fmt.Sprintf(fieldName + ":" + " Unexpected type: %v - expected: %v", "numeric", schemaField.Kind().String()))
				}
				return NewTypeError(fieldName, fmt.Sprintf(fieldName + ":" + " Unexpected type: %v - expected: %v", candidateField.Kind().String(), schemaField.Kind().String()))
			}


		}
	}
	return nil
}

