package main

import (
	"fmt"
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
		return ValidationResult{IsOk: false}, err
	}
	dataPathAbs, err := filepath.Abs(params.DataPath)
	if err != nil {
		return ValidationResult{IsOk: false}, err
	}

	schemaLoader := gojsonschema.NewSchemaLoader()
	err = schemaLoader.AddSchemas(
		gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", schemaPathAbs)),
	)
	if err != nil {
		return ValidationResult{IsOk: false}, err
	}

	schema, err := schemaLoader.Compile(gojsonschema.NewStringLoader(fmt.Sprintf(`{ "$ref": "#/definitions/%s" }`, params.DefinitionName)))
	if err != nil {
		return ValidationResult{IsOk: false}, err
	}

	documentLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", dataPathAbs))
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return ValidationResult{IsOk: false}, err
	}
	isOk := result.Valid()
	errors := result.Errors()
	error := fmt.Errorf("Found errors in data: %+v", errors)

	// guide: https://github.com/xeipuuv/gojsonschema
	return ValidationResult{IsOk: isOk, Error: error}, nil
}
