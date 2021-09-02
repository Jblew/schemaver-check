package main

import "testing"

func TestValidData(t *testing.T) {
	result, err := ValidateAgainstSpecificDefinition(ValidationParams{
		SchemaPath:     "mock/schema.json",
		DataPath:       "mock/data_valid.json",
		DefinitionName: "ChartSpec",
	})

	if err != nil {
		t.Errorf("Validation failed %+v", err)
	}

	if !result.IsOk {
		t.Errorf("Valid data marked as invalid %+v", result)
	}
}

func TestInvalidData(t *testing.T) {
	result, err := ValidateAgainstSpecificDefinition(ValidationParams{
		SchemaPath:     "mock/schema.json",
		DataPath:       "mock/data_invalid.json",
		DefinitionName: "ChartSpec",
	})

	if err != nil {
		t.Errorf("Validation failed %+v", err)
	}

	if result.IsOk {
		t.Errorf("Invalid data marked as valid %+v", result)
	}
}

func TestWrongDefinition(t *testing.T) {
	result, err := ValidateAgainstSpecificDefinition(ValidationParams{
		SchemaPath:     "mock/schema.json",
		DataPath:       "mock/data_valid.json",
		DefinitionName: "CompanyStructureSpec",
	})

	if err != nil {
		t.Errorf("Validation failed %+v", err)
	}

	if result.IsOk {
		t.Errorf("Valid data validated against wrong deginition marked as valid %+v", result)
	}
}

func TestMissingFile(t *testing.T) {
	_, err := ValidateAgainstSpecificDefinition(ValidationParams{
		SchemaPath:     "mock/no_such_file.json",
		DataPath:       "mock/data_valid.json",
		DefinitionName: "ChartSpec",
	})

	if err == nil {
		t.Errorf("Should fail on missing file")
	}
}
