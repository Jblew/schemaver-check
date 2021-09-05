package main

import (
	"fmt"
	"testing"
)

func TestReturnsFailureOnValidationFailed(t *testing.T) {
	conf := defaultConfig()
	conf.FnValidator = func(ValidationParams) (ValidationResult, error) {
		return ValidationResult{IsOk: false}, nil
	}
	conf.FnChecker = func(SchemaVerCompatibilityArgs) (SchemaVerCompatibilityResult, error) {
		return SchemaVerCompatibilityResult{IsValid: true}, nil
	}
	success, _, err := runApp(conf)

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
	success, _, err := runApp(conf)

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
	success, _, err := runApp(conf)

	if err == nil {
		t.Errorf("Error should not be nil")
	}

	if success {
		t.Errorf("Should fail")
	}
}

func TestReturnsFailureOnBothFailed(t *testing.T) {
	t.Errorf("Not specified yet")
}

func TestReturnsSuccessOnBothSucceeded(t *testing.T) {
	t.Errorf("Not specified yet")
}

func TestSkipsCheckWhenTold(t *testing.T) {
	t.Errorf("Not specified yet")
}

func TestReturnsSuccessWhenValidationSucceededAndCheckSkipped(t *testing.T) {
	t.Errorf("Not specified yet")
}

func TestRetriesCheckAndReturnsSuccessOnCorrectAttempt(t *testing.T) {
	t.Errorf("Not specified yet")
}

func TestRetriesCheckAndFailsAfterMaxAttempts(t *testing.T) {
	t.Errorf("Not specified yet")
}

func defaultConfig() AppConfig {
	return AppConfig{
		SkipSchemaCompatibilityCheck: false,
	}
}
