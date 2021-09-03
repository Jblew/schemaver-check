package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "jsonschema-compatibility-checker",
		Usage: "jsonschema-compatibility-checker --data-file [path] --definition-name [path]",
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
				SkipSchemaCompatibilityCheck:        os.Getenv("JSCC_SKIP_SCHEMAVER_COMPATIBILITY_CHECK") == "1",
				DefinitionName:                      c.String("definition-name"),
				SchemaFilePath:                      os.Getenv("JSCC_SCHEMA_FILE_PATH"),
				DataFilePath:                        c.String("data-file"),
				CompatibilityCheckEndpointURLFormat: os.Getenv("JSCC_COMPATIBILITY_CHECK_ENDPOINT_URL_FORMAT"),
				CompatibilityCheckRetryCount:        10,
				CompatibilityCheckRetryInterval:     1 * time.Second,
			}
			success, message, err := runApp(conf)
			if err != nil {
				return err
			}
			log.Println(message)
			if !success {
				os.Exit(1)
			}
			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
