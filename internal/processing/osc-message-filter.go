package processing

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jwetzell/osc-go"
)

type OSCMessageFilter struct {
	config  ProcessorConfig
	Address *regexp.Regexp
}

func (o *OSCMessageFilter) Process(ctx context.Context, payload any) (any, error) {

	payloadMessage, ok := payload.(osc.OSCMessage)

	if !ok {
		return nil, fmt.Errorf("osc.message.filter can only operate on OSCMessage payloads")
	}

	if !o.Address.MatchString(payloadMessage.Address) {
		return nil, nil
	}

	return payloadMessage, nil
}

func (o *OSCMessageFilter) Type() string {
	return o.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "osc.message.filter",
		New: func(config ProcessorConfig) (Processor, error) {
			params := config.Params
			address, ok := params["address"]

			if !ok {
				return nil, fmt.Errorf("osc.message.filter requires an address parameter")
			}

			addressString, ok := address.(string)

			if !ok {
				return nil, fmt.Errorf("osc.message.filter address must be a string")
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
