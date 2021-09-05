package main

import (
	"fmt"
	"testing"
	"time"
)

func TestReturnsFailureOnValidationFailed(t *testing.T) {
	conf := defaultConfig()
	conf.FnValidator = func(ValidationParams) (ValidationResult, error) {
		return ValidationResult{IsOk: false}, nil
	}
	conf.FnChecker = func(SchemaVerCompatibilityArgs) (SchemaVerCompatibilityResult, error) {
		return SchemaVerCompatibilityResult{IsValid: true}, nil
	}
	success, err := runApp(conf)

	if err != nil {
		t.Errorf("%+v", err)
	}

	if success {
		t.Errorf("Should fail")
	}
}

func TestReturnsFailureOnCheckFailed(t *testing.T) {
	conf := defaultConfig()
	conf.FnValidator = func(ValidationParams) (ValidationResult, error) {
		return ValidationResult{IsOk: true}, nil
	}
	conf.FnChecker = func(SchemaVerCompatibilityArgs) (SchemaVerCompatibilityResult, error) {
		return SchemaVerCompatibilityResult{IsValid: false}, nil
	}
	success, err := runApp(conf)

	if err != nil {
		t.Errorf("%+v", err)
	}

	if success {
		t.Errorf("Should fail")
	}
}

func TestReturnsFailureOnValidationError(t *testing.T) {
	conf := defaultConfig()
	conf.FnValidator = func(ValidationParams) (ValidationResult, error) {
		return ValidationResult{IsOk: false}, fmt.Errorf("Dummy error")
	}
	conf.FnChecker = func(SchemaVerCompatibilityArgs) (SchemaVerCompatibilityResult, error) {
		return SchemaVerCompatibilityResult{IsValid: true}, nil
	}
	success, err := runApp(conf)

	if err == nil {
		t.Errorf("Error should not be nil")
	}

	if success {
		t.Errorf("Should fail")
	}
}

func TestReturnsFailureOnBothFailed(t *testing.T) {
	conf := defaultConfig()
	conf.FnValidator = func(ValidationParams) (ValidationResult, error) {
		return ValidationResult{IsOk: false}, nil
	}
	conf.FnChecker = func(SchemaVerCompatibilityArgs) (SchemaVerCompatibilityResult, error) {
		return SchemaVerCompatibilityResult{IsValid: false}, nil
	}
	success, err := runApp(conf)

	if err != nil {
		t.Errorf("%+v", err)
	}

	if success {
		t.Errorf("Should fail")
	}
}

func TestReturnsSuccessOnBothSucceeded(t *testing.T) {
	conf := defaultConfig()
	conf.FnValidator = func(ValidationParams) (ValidationResult, error) {
		return ValidationResult{IsOk: true}, nil
	}
	conf.FnChecker = func(SchemaVerCompatibilityArgs) (SchemaVerCompatibilityResult, error) {
		return SchemaVerCompatibilityResult{IsValid: true}, nil
	}
	success, err := runApp(conf)

	if err != nil {
		t.Errorf("%+v", err)
	}

	if !success {
		t.Errorf("Should succeed")
	}
}

func TestByDefaultRunsBothValidatorAndChecker(t *testing.T) {
	conf := defaultConfig()
	validatorDone := false
	checkerDone := false
	conf.FnValidator = func(ValidationParams) (ValidationResult, error) {
		validatorDone = true
		return ValidationResult{IsOk: true}, nil
	}
	conf.FnChecker = func(SchemaVerCompatibilityArgs) (SchemaVerCompatibilityResult, error) {
		checkerDone = true
		return SchemaVerCompatibilityResult{IsValid: true}, nil
	}
	_, err := runApp(conf)

	if err != nil {
		t.Errorf("%+v", err)
	}

	if !validatorDone {
		t.Errorf("Should run FnValidator")
	}

	if !checkerDone {
		t.Errorf("Should run FnChecker")
	}
}

func TestSkipsCheckWhenTold(t *testing.T) {
	conf := defaultConfig()
	conf.SkipSchemaCompatibilityCheck = true
	validatorDone := false
	checkerDone := false
	conf.FnValidator = func(ValidationParams) (ValidationResult, error) {
		validatorDone = true
		return ValidationResult{IsOk: true}, nil
	}
	conf.FnChecker = func(SchemaVerCompatibilityArgs) (SchemaVerCompatibilityResult, error) {
		checkerDone = true
		return SchemaVerCompatibilityResult{IsValid: true}, nil
	}
	_, err := runApp(conf)

	if err != nil {
		t.Errorf("%+v", err)
	}

	if !validatorDone {
		t.Errorf("Should run FnValidator")
	}

	if checkerDone {
		t.Errorf("Should skip runing FnChecker")
	}
}

func TestReturnsSuccessWhenValidationSucceededAndCheckSkipped(t *testing.T) {
	conf := defaultConfig()
	conf.SkipSchemaCompatibilityCheck = true
	validatorDone := false
	checkerDone := false
	conf.FnValidator = func(ValidationParams) (ValidationResult, error) {
		validatorDone = true
		return ValidationResult{IsOk: true}, nil
	}
	conf.FnChecker = func(SchemaVerCompatibilityArgs) (SchemaVerCompatibilityResult, error) {
		checkerDone = true
		return SchemaVerCompatibilityResult{IsValid: false}, nil
	}
	success, err := runApp(conf)

	if err != nil {
		t.Errorf("%+v", err)
	}

	if !validatorDone {
		t.Errorf("Should run FnValidator")
	}

	if checkerDone {
		t.Errorf("Should skip runing FnChecker")
	}

	if !success {
		t.Errorf("Should succeed")
	}
}

func TestRetriesCheckAndReturnsSuccessOnCorrectAttempt(t *testing.T) {
	conf := defaultConfig()
	conf.CompatibilityCheckRetryCount = 5
	checkerCount := 0
	conf.FnValidator = func(ValidationParams) (ValidationResult, error) {
		return ValidationResult{IsOk: true}, nil
	}
	conf.FnChecker = func(SchemaVerCompatibilityArgs) (SchemaVerCompatibilityResult, error) {
		checkerCount++
		var err error
		if checkerCount < 3 {
			err = fmt.Errorf("Cannot connect yet")
		} else {
			err = nil
		}
		return SchemaVerCompatibilityResult{IsValid: true}, err
	}
	success, err := runApp(conf)

	if err != nil {
		t.Errorf("%+v", err)
	}

	if checkerCount != 3 {
		t.Errorf("Should run FnChecker 3 times (%d instead)", checkerCount)
	}

	if !success {
		t.Errorf("Should succeed")
	}
}

func TestRetriesCheckAndFailsAfterMaxAttempts(t *testing.T) {
	conf := defaultConfig()
	conf.CompatibilityCheckRetryCount = 5
	checkerCount := 0
	conf.FnValidator = func(ValidationParams) (ValidationResult, error) {
		return ValidationResult{IsOk: true}, nil
	}
	conf.FnChecker = func(SchemaVerCompatibilityArgs) (SchemaVerCompatibilityResult, error) {
		checkerCount++
		var err error
		if checkerCount < 6 {
			err = fmt.Errorf("Cannot connect yet")
		} else {
			err = nil
		}
		return SchemaVerCompatibilityResult{IsValid: false}, err
	}
	success, err := runApp(conf)

	if err == nil {
		t.Errorf("Should return error")
	}

	if checkerCount != 5 {
		t.Errorf("Should run FnChecker 5 times")
	}

	if success {
		t.Errorf("Should fail")
	}
}

func defaultConfig() AppConfig {
	return AppConfig{
		SkipSchemaCompatibilityCheck:    false,
		CompatibilityCheckRetryCount:    1,
		CompatibilityCheckRetryInterval: 5 * time.Millisecond,
	}
}
