package main

import (
	"fmt"
	"time"
)

type AppConfig struct {
	SkipSchemaCompatibilityCheck        bool
	DefinitionName                      string
	SchemaFilePath                      string
	DataFilePath                        string
	CompatibilityCheckEndpointURLFormat string
	CompatibilityCheckRetryCount        int
	CompatibilityCheckRetryInterval     time.Duration
	FnValidator                         func(ValidationParams) (ValidationResult, error)
	FnChecker                           func(SchemaVerCompatibilityArgs) (SchemaVerCompatibilityResult, error)
}

func runApp(config AppConfig) (bool, string, error) {
	out := "\n"
	validationResult, err := config.FnValidator(ValidationParams{
		SchemaPath:     config.SchemaFilePath,
		DataPath:       config.DataFilePath,
		DefinitionName: config.DefinitionName,
	})
	if err != nil {
		return false, "", fmt.Errorf("Cannot validate json data: %+v", err)
	}
	if validationResult.IsOk {
		out += fmt.Sprintf("JSON Schema validation success: data file %s is valid against definition %s\n", config.DataFilePath, config.DefinitionName)
	} else {
		out += fmt.Sprintf("JSON Schema validation failed: data file %s is INVALID against definition %s\n", config.DataFilePath, config.DefinitionName)
	}

	if config.SkipSchemaCompatibilityCheck {
		out += "SchemaVer compatibility check skipped.\n"
		return validationResult.IsOk, out, nil
	}
	checkResult := SchemaVerCompatibilityResult{}
	err = retry(config.CompatibilityCheckRetryCount, config.CompatibilityCheckRetryInterval, func() error {
		checkResultAtAttempt, err := config.FnChecker(SchemaVerCompatibilityArgs{
			EndpointURLFormat: config.CompatibilityCheckEndpointURLFormat,
			SchemaPath:        config.SchemaFilePath,
			DefinitionName:    config.DefinitionName,
		})
		checkResult = checkResultAtAttempt
		return err
	})
	if err != nil {
		return false, out, fmt.Errorf("Cannot check SchemaVer compatibility: %+v", err)
	}
	if checkResult.IsValid {
		out += fmt.Sprintf("JSON SchemaVer compatibility check success\n")
	} else {
		out += fmt.Sprintf("JSON SchemaVer compatibility check failed %s\n", checkResult.ErrorMsg)
	}
	return validationResult.IsOk && checkResult.IsValid, out, nil
}
