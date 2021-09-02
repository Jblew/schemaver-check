package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

type ValidationParams struct {
	SchemaPath     string
	DataPath       string
	DefinitionName string
}

type ValidationResult struct {
	IsOk  bool
	Error error
}

func ValidateAgainstSpecificDefinition(params ValidationParams) (ValidationResult, error) {
	schemaPathAbs, err := filepath.Abs(params.SchemaPath)
	if err != nil {
		return ValidationResult{IsOk: false}, fmt.Errorf("Cannot obtain absolute schema path: %+v", err)
	}
	dataPathAbs, err := filepath.Abs(params.DataPath)
	if err != nil {
		return ValidationResult{IsOk: false}, fmt.Errorf("Cannot obtain absolute data path: %+v", err)
	}

	schemaBytes, _ := ioutil.ReadFile(schemaPathAbs)
	schemaMap := make(map[string]interface{})
	json.Unmarshal(schemaBytes, &schemaMap)
	if err != nil {
		return ValidationResult{IsOk: false}, fmt.Errorf("Cannot unmarhall schema: %+v", err)
	}
	// insert $ref to specify the root type
	schemaMap["$ref"] = params.DefinitionName

	schemaLoader := gojsonschema.NewGoLoader(schemaMap)
	documentLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", dataPathAbs))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return ValidationResult{IsOk: false}, err
	}
	isOk := result.Valid()
	errors := result.Errors()
	error := fmt.Errorf("Found errors in data: %+v", errors)

	return ValidationResult{IsOk: isOk, Error: error}, nil
}
