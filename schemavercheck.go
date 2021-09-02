package main

import "fmt"

type SchemaVerCompatibilityArgs struct {
	SchemaPath        string
	EndpointURLFormat string
	DefinitionName    string
}

type SchemaVerCompatibilityResult struct {
	IsValid  bool
	ErrorMsg string
}

func CheckSchemaVerCompatibility(args SchemaVerCompatibilityArgs) (SchemaVerCompatibilityResult, error) {
	// url := fmt.Sprintf(args.EndpointURLFormat, args.DefinitionName)
	return SchemaVerCompatibilityResult{IsValid: false}, fmt.Errorf("Not implemented yet")
}
