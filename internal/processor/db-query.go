package processor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type DbQuery struct {
	config   config.ProcessorConfig
	ModuleId string
	Query    string
	logger   *slog.Logger
}

func (dq *DbQuery) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	ctxModules := ctx.Value(common.ModulesContextKey)
	if ctxModules == nil {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("db.query unable to get modules from context")
	}

	moduleMap, ok := ctxModules.(map[string]common.Module)
	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("db.query modules from context has wrong type")
	}

	module, ok := moduleMap[dq.ModuleId]
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

	rows, err := db.QueryContext(ctx, dq.Query)
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
			rowMap[colName] = columnValues[i]
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
		Type: "db.query",
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
			return &DbQuery{config: config, ModuleId: moduleIdString, Query: queryString, logger: slog.Default().With("component", "processor", "type", config.Type)}, nil
		},
	})
}
