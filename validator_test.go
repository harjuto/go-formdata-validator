package form_validator

import (
	"testing"
)


func TestInputNotAStruct(t *testing.T) {
	//var formData ValidateableFormData
	//err := formData.Validate("test")
	//
	//if err.Error() != "Form validator can only validate structs." {
	//	t.Errorf("Struct type validation failed: %v", err.Error())
	//}
}

func TestValidateEmptyStruct(t *testing.T) {
	//var formData ValidateableFormData
	//type emptyStruct struct {}
	//err := formData.Validate(emptyStruct{})
	//
	//if err != nil {
	//	t.Errorf("Unexpected result when validating with empty schema struct")
	//}
}

type skill struct {
	Endorsements int
	Name string
}

type friend struct {
	Id int
	Name string
}

type contactDetailsSchema struct {
	Id int
	Name string
	Username string
	Email string
	Address struct {
		Street string
		Suite string
		City string
		Zipcode string
		Geo struct {
			Lat string
			Lng string
	    	}
	}
	Phone string
	Website string
	Company struct {
		Name string
		CatchPhrase string
		Bs string
	}
	Skills []skill
}

type AccountListSchema []AccountSchema

type AccountSchema struct {
	_Id string
	Guid string
	IsActive bool
	Balance string
	Tags []string
	Friends []friend

}

var contactDetailsJson []byte = []byte (
	`{
		"id": 1,
		"name": "Leanne Graham",
		"username": "Bret",
		"email": "Sincere@april.biz",
		"address": {
			"street": "Kulas Light",
			"suite": "Apt. 556",
			"city": "Gwenborough",
			"zipcode": "92998-3874",
			"geo": {
				"lat": "-37.3159",
				"lng": "81.1496"
			}
		},
		"phone": "1-770-736-8031 x56442",
		"website": "hildegard.org",
		"company": {
			"name": "Romaguera-Crona",
			"catchPhrase": "Multi-layered client-server neural-net",
			"bs": "harness real-time e-markets"
		},
		"skills": [
			{
				"endorsements": 99,
				"name": "Go"
			},
			{
				"endorsements": 50,
				"name": "JavaScript"
			},
			{
				"endorsements": 1,
				"name": "Programming"
			}
		]
	}`,
)

var accountListJson []byte = []byte (
	`[
		{
			"_id": "582c20e9ce2d4ba629fb341e",
			"guid": "9c25ca58-45f9-4296-af51-8800e7477d12",
			"isActive": true,
			"balance": "$3,125.66",
			"tags": [
				"quis",
				"incididunt",
				"aliquip",
				"sit",
				"culpa",
				"consectetur",
				"eiusmod"
			],
			"friends": [
				{
					"id": 0,
					"name": "Talley Roberts"
				},
				{
					"id": 1,
					"name": "Marva Barry"
				},
				{
					"id": 2,
					"name": "Millie Orr"
				}
			]
		},
		{
			"_id": "582c20e958b5df36811d14e3",
			"guid": "ef6bc6c4-58cc-4a3d-b201-225adc626d2c",
			"isActive": true,
			"balance": "$3,666.76",
			"tags": [
				"anim",
				"mollit",
				"culpa",
				"laboris",
				"culpa",
				"est",
				"veniam"
			],
			"friends": [
				{
					"id": 0,
					"name": "Stewart Stout"
				},
				{
					"id": 1,
					"name": "Jacobs Lowe"
				},
				{
					"id": 2,
					"name": "Swanson Randall"
				}
			]
		},
		{
			"_id": "582c20e9113059c726148e91",
			"guid": "6f9bcef8-7176-453e-9a61-d9334b69b3bb",
			"isActive": true,
			"balance": "$2,665.92",
			"tags": [
				"ullamco",
				"consequat",
				"consequat",
				"nostrud",
				"consequat",
				"enim",
				"dolore"
			],
			"friends": [
				{
					"id": 0,
					"name": "Letha Gentry"
				},
				{
					"id": 1,
					"name": "Hancock Stein"
				},
				{
					"id": 2,
					"name": "Cathleen Hale"
				}
			]
		}
	]`,
)





func TestSchemaValidation(t *testing.T) {
	//err := ValidateSchema(contactDetailsJson, contactDetailsSchema{})
	//if err != nil {
	//	t.Errorf("Failed: %v", err.Error())
	//}

	err := ValidateSchema(accountListJson, AccountListSchema{})

	if err != nil {
		t.Errorf("Failed: %v", err.Error())
	}
}



