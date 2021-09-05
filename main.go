package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "schemaver-check",
		Usage: "schemaver-check --data-file [path] --definition-name [path]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "data-file",
				Usage:    "Local mock JSON data file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "definition-name",
				Usage:    "JSON Schema definition name",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			conf := AppConfig{
				SkipSchemaCompatibilityCheck:        os.Getenv("SCHEMAVERCHECK_SKIP_REMOTE_CHECK") == "1",
				DefinitionName:                      c.String("definition-name"),
				SchemaFilePath:                      os.Getenv("SCHEMAVERCHECK_SCHEMA_FILE"),
				DataFilePath:                        c.String("data-file"),
				CompatibilityCheckEndpointURLFormat: os.Getenv("SCHEMAVERCHECK_ENDPOINT_URL_FORMAT"),
				CompatibilityCheckRetryCount:        10,
				CompatibilityCheckRetryInterval:     500 * time.Millisecond,
				FnValidator:                         ValidateAgainstSpecificDefinition,
				FnChecker:                           CheckSchemaVerCompatibility,
			}
			success, err := runApp(conf)
			if err != nil {
				return err
			}
			if !success {
				os.Exit(1)
			}
			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
