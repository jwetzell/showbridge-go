package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

type HTTPServer struct {
	config       config.ModuleConfig
	Port         uint16
	ctx          context.Context
	inputHandler common.InputHandler
	logger       *slog.Logger
	cancel       context.CancelFunc
	server       *http.Server
	serverMu     sync.Mutex
}

type ResponseIOError struct {
	Index        int     `json:"index"`
	ProcessError *string `json:"processError"`
}

type IOResponseData struct {
	IOErrors []ResponseIOError `json:"ioErrors"`
	Message  string            `json:"message"`
	Status   string            `json:"status"`
}

type httpServerContextKey string

type HTTPServerResponseWriter struct {
	http.ResponseWriter
	done bool
}

func (hsrw *HTTPServerResponseWriter) WriteHeader(status int) {
	hsrw.done = true
	hsrw.ResponseWriter.WriteHeader(status)
}

func (hsrw *HTTPServerResponseWriter) Write(data []byte) (int, error) {
	hsrw.done = true
	return hsrw.ResponseWriter.Write(data)
}

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "http.server",
		Title: "HTTP Server",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"port": {
					Title:   "Port",
					Type:    "integer",
					Minimum: jsonschema.Ptr[float64](1024),
					Maximum: jsonschema.Ptr[float64](65535),
				},
			},
			Required:             []string{"port"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params
			portNum, err := params.GetInt("port")
			if err != nil {
				return nil, fmt.Errorf("http.server port error: %w", err)
			}
			return &HTTPServer{Port: uint16(portNum), config: config, logger: CreateLogger(config)}, nil
		},
	})
}

func (hs *HTTPServer) Id() string {
	return hs.config.Id
}

func (hs *HTTPServer) Type() string {
	return hs.config.Type
}

func (hs *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responseWriter := HTTPServerResponseWriter{ResponseWriter: w}

	response := IOResponseData{
		Message: "routing successful",
		Status:  "ok",
	}
	if hs.inputHandler != nil {
		inputContext := context.WithValue(hs.ctx, httpServerContextKey("responseWriter"), &responseWriter)
		aRouteFound, routingErrors := hs.inputHandler(inputContext, hs.Id(), r)
		if !responseWriter.done {
			if aRouteFound {
				if routingErrors != nil {
					w.WriteHeader(http.StatusInternalServerError)
					response.Status = "error"
					response.Message = "routing failed"

					response.IOErrors = []ResponseIOError{}
					for _, responseIOError := range routingErrors {
						errorToAdd := ResponseIOError{
							Index: responseIOError.Index,
						}

						if responseIOError.ProcessError != nil {
							errorMsg := responseIOError.ProcessError.Error()
							errorToAdd.ProcessError = &errorMsg
						}

						response.IOErrors = append(response.IOErrors, errorToAdd)

					}
					json.NewEncoder(w).Encode(response)
					return
				} else {
					w.WriteHeader(http.StatusOK)
					response.Message = "routing successful"
					json.NewEncoder(w).Encode(response)
					return
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
				response.Status = "error"
				response.Message = "no matching routes found"
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = "no router registered"
		response.Status = "error"
		json.NewEncoder(w).Encode(response)
		return
	}
}

func (hs *HTTPServer) Start(ctx context.Context, inputHandler common.InputHandler) error {
	hs.logger.Debug("running")
	hs.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	hs.ctx = moduleContext
	hs.cancel = cancel

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", hs.Port),
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           hs,
	}

	hs.serverMu.Lock()
	hs.server = httpServer
	hs.serverMu.Unlock()

	err := httpServer.ListenAndServe()
	// TODO(jwetzell): handle server closed error differently
	if err != nil {
		if err.Error() != "http: Server closed" {
			return err
		}
	}

	<-hs.ctx.Done()
	return nil
}

func (hs *HTTPServer) Output(ctx context.Context, payload any) error {
	responseWriter, ok := ctx.Value(httpServerContextKey("responseWriter")).(*HTTPServerResponseWriter)

	if !ok {
		return errors.New("http.server output must originate from an http.server input")
	}

	payloadResponse, ok := common.GetAnyAs[processor.HTTPResponse](payload)

	if !ok {
		return errors.New("http.server is only able to output HTTPResponse")
	}

	if responseWriter.done {
		return errors.New("http.server response writer has already been written to")
	}

	responseWriter.WriteHeader(payloadResponse.Status)
	responseWriter.Write(payloadResponse.Body)
	return nil
}

func (hs *HTTPServer) Stop() {
	if hs.cancel != nil {
		hs.cancel()
	}
	hs.serverMu.Lock()
	defer hs.serverMu.Unlock()
	if hs.server != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		hs.server.Shutdown(shutdownCtx)
		shutdownCancel()
		<-shutdownCtx.Done()
		hs.server = nil
	}
	hs.logger.Debug("done")
}
