package processor

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type OSCMessageFilter struct {
	config  config.ProcessorConfig
	Address *regexp.Regexp
}

func (omf *OSCMessageFilter) Process(ctx context.Context, payload any) (any, error) {

	payloadMessage, ok := payload.(osc.OSCMessage)

	if !ok {
		return nil, errors.New("osc.message.filter can only operate on OSCMessage payloads")
	}

	if !omf.Address.MatchString(payloadMessage.Address) {
		return nil, nil
	}

	return payloadMessage, nil
}

func (omf *OSCMessageFilter) Type() string {
	return omf.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "osc.message.filter",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params
			addressString, err := params.GetString("address")
			if err != nil {
				return nil, fmt.Errorf("osc.message.filter address error: %w", err)
			}

			addressPattern := strings.ReplaceAll(addressString, "?", ".")
			addressPattern = strings.ReplaceAll(addressPattern, "*", "[^/]*")
			addressPattern = strings.ReplaceAll(addressPattern, "[!", "[^")
			addressPattern = strings.ReplaceAll(addressPattern, "{", "(")
			addressPattern = strings.ReplaceAll(addressPattern, "}", ")")
			addressPattern = strings.ReplaceAll(addressPattern, ",", "|")

			addressPatternRegexp, err := regexp.Compile(fmt.Sprintf("^%s$", addressPattern))

			if err != nil {
				return nil, err
			}

			return &OSCMessageFilter{config: config, Address: addressPatternRegexp}, nil
		},
	})
}
