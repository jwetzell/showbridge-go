package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/schema"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("schema type argument is required")
		return
	}
	schemaType := os.Args[1]
	var err error
	switch schemaType {
	case "processors":
		err = outputSchema(os.Stdout, schema.GetProcessorsSchema())
	case "modules":
		err = outputSchema(os.Stdout, schema.GetModulesSchema())
	case "config":
		err = outputSchema(os.Stdout, &schema.ConfigSchema)
	case "routes":
		err = outputSchema(os.Stdout, &schema.RoutesConfigSchema)
	case "api":
		err = outputSchema(os.Stdout, &schema.ApiConfigSchema)
	default:
	}

	if err != nil {
		panic(err)
	}
}

func outputSchema(writer io.Writer, jsonSchema *jsonschema.Schema) error {
	schemaJSON, err := json.Marshal(jsonSchema)
	if err != nil {
		panic(err)
	}
	_, err = writer.Write(schemaJSON)
	return err
}
