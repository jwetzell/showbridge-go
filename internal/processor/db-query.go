package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"text/template"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type DbQuery struct {
	config   config.ProcessorConfig
	ModuleId string
	Query    *template.Template
	logger   *slog.Logger
}

func (dq *DbQuery) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
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

	db := dbModule.Database()
	if db == nil {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("db.query module with id %s returned nil database", dq.ModuleId)
	}

	var queryBuffer bytes.Buffer
	err := dq.Query.Execute(&queryBuffer, wrappedPayload)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	// support proper parameterized queries
	rows, err := db.QueryContext(ctx, queryBuffer.String())
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

	results := make([]map[string]any, 0)

	for rows.Next() {
		columnValues := make([]interface{}, len(columns))

		for i := range columnValues {
			columnValues[i] = new(interface{})
		}

		if err := rows.Scan(columnValues...); err != nil {
			wrappedPayload.End = true
			return wrappedPayload, fmt.Errorf("db.query error scanning row: %w", err)
		}

		rowMap := make(map[string]any)
		for i, colName := range columns {
			value := *columnValues[i].(*interface{})
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

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "db.query",
		Title: "Query Database",
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
