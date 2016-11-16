package form_validator

import (
	"reflect"
	"bytes"
	"errors"
	"fmt"
	"log"
)

type ValidateableFormData map[string]interface{}

type TypeErrors struct {
	errors []error
}

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

// Schema = The struct we validate against
// fieldValue = The attempted interface we try to assign to schema field
func (f ValidateableFormData) Validate(t interface{}) error {

	var typeErrors TypeErrors

	if reflect.Indirect(reflect.ValueOf(t)).Kind() != reflect.Struct {
		return errors.New("Form validator can only validate structs.")
	}

	var schema reflect.Value = reflect.Indirect(reflect.ValueOf(t))

	validateFields(f, schema, &typeErrors)

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
		log.Printf("Expected: %v - Candidate: %v", schemaField.Type().Name(), candidateField.Kind())

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




