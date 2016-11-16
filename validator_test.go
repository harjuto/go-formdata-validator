package form_validator

import (
	"testing"
	"encoding/json"
	//"reflect"
)

type person struct {
	Name string
	Age int
	Weight float64
	Admin bool
}

type contactDetails struct {
	Name string
	Age int
	Weight float64
	Admin bool
	Address address
}

type address struct {
	Street string
	Number int
}

type types struct {
	integerType int
	stringType string
	float32Type float32
	float64Type float64
	boolType bool
}

func TestInputNotAStruct(t *testing.T) {
	var formData ValidateableFormData
	err := formData.Validate("test")

	if err.Error() != "Form validator can only validate structs." {
		t.Errorf("Struct type validation failed: %v", err.Error())

	}
}

func TestFieldValidation(t *testing.T) {
	//var typeStruct reflect.Value = reflect.Indirect(reflect.ValueOf(&types{}))
	//var field reflect.Value
	//var kind reflect.Kind
	//
	//// Pass
	//field = typeStruct.FieldByName("integerType")
	//kind = reflect.Float64
	//if err := validateField(field, kind, ) ; err != nil {
	//	t.Errorf("Failed to validate: %v as %v", field.Kind().String(), kind.String())
	//}
	//// Pass
	//field = typeStruct.FieldByName("integerType")
	//kind = reflect.Float64
	//if err := validateField(field, kind) ; err != nil {
	//	t.Errorf("Failed to validate: %v as %v", field.Kind().String(), kind.String())
	//}
	//// Pass
	//field = typeStruct.FieldByName("stringType")
	//kind = reflect.String
	//if err := validateField(field, kind) ; err != nil {
	//	t.Errorf("Failed to validate: %v as %v", field.Kind().String(), kind.String())
	//}
	//// Fail
	//field = typeStruct.FieldByName("integerType")
	//kind = reflect.String
	//if err := validateField(field, kind) ; err == nil {
	//	t.Errorf("Failed to validate: %v as %v", field.Kind().String(), kind.String())
	//}
	// Fail
	// Fail
}

func TestValidationCorrectData(t *testing.T) {
	var requestBody []byte = []byte(`
		{
			"Name":"John",
			"age":1,
			"weight":80,
			"admin":true,
			"kissa": "meow!",
			"address": {
				"street": "Sunset Ave",
				"number": 1
			}
		}
	`)


	var formData ValidateableFormData
	err := json.Unmarshal(requestBody, &formData)
	if err != nil {
		t.Errorf("Unable to parse %v", err.Error())
	}
	err = formData.Validate(contactDetails{})

	if err != nil {
		t.Errorf("Result: %v", err)
	}

}
