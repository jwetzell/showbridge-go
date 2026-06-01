package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"text/template"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "db.query",
		Title: "Query Database",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"module": {
					Title:       "Module ID",
					Description: "ID of the database module to query",
					Type:        "string",
				},
				"query": {
					Title:       "Query",
					Description: "SQL query to execute",
					Type:        "string",
				},
			},
			Required:             []string{"module", "query"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (Processor, error) {

			params := config.Params

			moduleIdString, err := params.GetString("module")
			if err != nil {
				return nil, fmt.Errorf("db.query module error: %w", err)
			}

			queryString, err := params.GetString("query")
			if err != nil {
				return nil, fmt.Errorf("db.query query error: %w", err)
			}

			queryTemplate, err := template.New("query").Parse(queryString)

			if err != nil {
				return nil, err
			}
			return &DbQuery{config: config, ModuleId: moduleIdString, Query: queryTemplate, logger: slog.Default().With("component", "processor", "type", config.Type)}, nil
		},
	})
}

type DbQuery struct {
	config   config.ProcessorConfig
	ModuleId string
	Query    *template.Template
	logger   *slog.Logger
	module   common.DatabaseModule
}

func (dq *DbQuery) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	if dq.module == nil {
		if wrappedPayload.Modules == nil {
			wrappedPayload.End = true
			return wrappedPayload, errors.New("db.query wrapped payload has no modules")
		}

		module, ok := wrappedPayload.Modules[dq.ModuleId]
		if !ok {
			wrappedPayload.End = true
			return wrappedPayload, fmt.Errorf("db.query unable to find module with id: %s", dq.ModuleId)
		}

		dbModule, ok := module.(common.DatabaseModule)
		if !ok {
			wrappedPayload.End = true
			return wrappedPayload, fmt.Errorf("db.query module with id %s is not a DatabaseModule", dq.ModuleId)
		}
		dq.module = dbModule
	}

	var queryBuffer bytes.Buffer
	err := dq.Query.Execute(&queryBuffer, wrappedPayload)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	// support proper parameterized queries
	rows, err := dq.module.QueryContext(ctx, queryBuffer.String())
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("db.query error executing query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("db.query error getting columns: %w", err)
	}

	// TODO(jwetzell): optimize this
	results := make([]map[string]any, 0)

	for rows.Next() {
		columnValues := make([]any, len(columns))

		for i := range columnValues {
			columnValues[i] = new(any)
		}

		if err := rows.Scan(columnValues...); err != nil {
			wrappedPayload.End = true
			return wrappedPayload, fmt.Errorf("db.query error scanning row: %w", err)
		}

		rowMap := make(map[string]any)
		for i, colName := range columns {
			value := *columnValues[i].(*any)
			rowMap[colName] = value
		}
		results = append(results, rowMap)
	}

	if len(results) == 0 {
		wrappedPayload.Payload = nil
		return wrappedPayload, nil
	} else if len(results) == 1 {
		wrappedPayload.Payload = results[0]
		return wrappedPayload, nil
	}
	wrappedPayload.Payload = results
	return wrappedPayload, nil
}

func (dq *DbQuery) Type() string {
	return dq.config.Type
}
