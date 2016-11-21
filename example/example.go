package example

import (
	"github.com/harjuto/go-formdata-validator"
	"fmt"
)

type exampleSchema struct {
	Name string
	Age int
	Admin bool
}

func init() {
	jsonBody := []byte(`
		{"name": "John", "age":30, "admin": false}`,
	)

	err := form_validator.ValidateSchema(jsonBody, exampleSchema{})

	fmt.Print(err)
}
