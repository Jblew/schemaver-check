package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "verify-data",
				Usage: "verify-data [json-schema-file] [json file]",
				Action: func(c *cli.Context) error {
					fmt.Println("Not implemented yet")
					return nil
				},
			},
			{
				Name:  "assert-compatibility",
				Usage: "assert-compatibility [json-schema-file] [compatibility_url]",
				Action: func(c *cli.Context) error {
					fmt.Println("Not implemented yet")
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
