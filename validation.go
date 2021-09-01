package main

import "fmt"

type ValidationParams struct {
	SchemaPath     string
	DataPath       string
	DefinitionName string
}

type ValidationResult struct {
	IsOk  bool
	Error error
}

func ValidateAgainstSpecificDefinition(params ValidationParams) (error, ValidationResult) {
	// guide: https://github.com/xeipuuv/gojsonschema
	return fmt.Errorf("Not implemented yet"), ValidationResult{IsOk: false, Error: nil}
}
