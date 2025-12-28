package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type HTTPServer struct {
	config config.ModuleConfig
	Port   uint16
	ctx    context.Context
	router route.RouteIO
	logger *slog.Logger
}

type ResponseIOError struct {
	Index        int      `json:"index"`
	OutputErrors []string `json:"outputErrors"`
	ProcessError *string  `json:"processError"`
	InputError   *string  `json:"inputError"`
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
		Type: "http.server",
		New: func(ctx context.Context, config config.ModuleConfig) (Module, error) {
			params := config.Params
			port, ok := params["port"]
			if !ok {
				return nil, errors.New("http.server requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, errors.New("http.server port must be uint16")
			}

			router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

			if !ok {
				return nil, errors.New("http.server unable to get router from context")
			}

			return &HTTPServer{Port: uint16(portNum), config: config, ctx: ctx, router: router, logger: CreateLogger(config)}, nil
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
	if hs.router != nil {
		inputContext := context.WithValue(hs.ctx, httpServerContextKey("responseWriter"), &responseWriter)
		aRouteFound, routingErrors := hs.router.HandleInput(inputContext, hs.Id(), r)
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

						if responseIOError.InputError != nil {
							errorMsg := responseIOError.InputError.Error()
							errorToAdd.InputError = &errorMsg
						}

						if responseIOError.ProcessError != nil {
							errorMsg := responseIOError.ProcessError.Error()
							errorToAdd.ProcessError = &errorMsg
						}

						if responseIOError.OutputErrors != nil {
							outputErrorMsgs := []string{}

							for _, outputError := range responseIOError.OutputErrors {
								outputErrorMsgs = append(outputErrorMsgs, outputError.Error())
							}

							errorToAdd.OutputErrors = outputErrorMsgs
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

func (hs *HTTPServer) Run() error {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", hs.Port),
		Handler: hs,
	}

	go func() {
		<-hs.ctx.Done()
		httpServer.Close()
	}()

	err := httpServer.ListenAndServe()
	// TODO(jwetzell): handle server closed error differently
	if err != nil {
		if err.Error() != "http: Server closed" {
			return err
		}
	}

	<-hs.ctx.Done()
	hs.logger.Debug("done")
	return nil
}

func (hs *HTTPServer) Output(ctx context.Context, payload any) error {
	responseWriter, ok := ctx.Value(httpServerContextKey("responseWriter")).(*HTTPServerResponseWriter)

	if !ok {
		return errors.New("http.server output must originate from an http.server input")
	}

	payloadResponse, ok := payload.(processor.HTTPResponse)

	if !ok {
		return errors.New("http.server is only able to output HTTPResponse")
	}

	responseWriter.WriteHeader(payloadResponse.Status)

	responseWriter.Write(payloadResponse.Body)

	return nil
}
