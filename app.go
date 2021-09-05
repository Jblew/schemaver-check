package main

import (
	"fmt"
	"log"
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

func runApp(config AppConfig) (bool, error) {
	validationResult, err := config.FnValidator(ValidationParams{
		SchemaPath:     config.SchemaFilePath,
		DataPath:       config.DataFilePath,
		DefinitionName: config.DefinitionName,
	})
	if err != nil {
		return false, fmt.Errorf("Cannot validate json data: %+v", err)
	}
	if validationResult.IsOk {
		log.Printf("JSON Schema validation success: data file %s is valid against definition %s", config.DataFilePath, config.DefinitionName)
	} else {
		log.Printf("JSON Schema validation failed: data file %s is INVALID against definition %s", config.DataFilePath, config.DefinitionName)
	}

	if config.SkipSchemaCompatibilityCheck {
		log.Println("SchemaVer compatibility check skipped.")
		return validationResult.IsOk, nil
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
		return false, fmt.Errorf("Cannot check SchemaVer compatibility: %+v", err)
	}
	if checkResult.IsValid {
		log.Println("JSON SchemaVer compatibility check success")
	} else {
		log.Printf("JSON SchemaVer compatibility check failed %s", checkResult.ErrorMsg)
	}
	return validationResult.IsOk && checkResult.IsValid, nil
}
