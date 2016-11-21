# go-formdata-validator
Small library which uses structs to validate json data sent to server. Runs validation only, does not marshal json into given schema struct.

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/harjuto/go-formdata-validator) 

# Usage

For setup and function calls, See example folder. 
Produces below json as a result
``` json
    {
        "errors": [
            {
                "field":"Phone",
                "message":"Phone: Unexpected type: numeric - expected: string"
            },
            {
                "field":"Skills",
                "message":"Skills: Unexpected type: map - expected: slice"
            },
            {
                "field":"Id",
                "message":"Unexpected type: string - expected: numeric"
            }
        ]
    }
```

# 
# TODO

1. Marshalling valid json to given struct.
2. Run some benchmarks
3. Figure how to use this as a middleware

# Issues
1. Does not validate ints and floats properly. Floats are truncated to ints without a error / hint message.
