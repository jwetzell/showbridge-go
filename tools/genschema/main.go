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
	switch schemaType {
	case "processors":
		outputSchema(os.Stdout, schema.GetProcessorsSchema())
	case "modules":
		outputSchema(os.Stdout, schema.GetModulesSchema())
	case "config":
		outputSchema(os.Stdout, &schema.ConfigSchema)
	case "routes":
		outputSchema(os.Stdout, &schema.RoutesConfigSchema)
	case "api":
		outputSchema(os.Stdout, &schema.ApiConfigSchema)
	default:
	}

}

func outputSchema(writer io.Writer, jsonSchema *jsonschema.Schema) {
	schemaJSON, err := json.Marshal(jsonSchema)
	if err != nil {
		panic(err)
	}
	writer.Write(schemaJSON)
}
