package showbridge

import (
	"fmt"
	"log/slog"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

type MIDIClient struct {
	config     ModuleConfig
	router     *Router
	InputPort  string
	OutputPort string
	SendFunc   func(midi.Message) error
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "misc.midi.client",
		New: func(config ModuleConfig) (Module, error) {
			params := config.Params
			input, ok := params["input"]

			if !ok {
				return nil, fmt.Errorf("net.mqtt.client requires a input parameter")
			}

			inputString, ok := input.(string)

			if !ok {
				return nil, fmt.Errorf("misc.midi.client input must be a string")
			}

			output, ok := params["output"]

			if !ok {
				return nil, fmt.Errorf("net.mqtt.client requires a output parameter")
			}

			outputString, ok := output.(string)

			if !ok {
				return nil, fmt.Errorf("misc.midi.client output must be a string")
			}

			return &MIDIClient{config: config, InputPort: inputString, OutputPort: outputString}, nil
		},
	})
}

func (mc *MIDIClient) Id() string {
	return mc.config.Id
}

func (mc *MIDIClient) Type() string {
	return mc.config.Type
}

func (mc *MIDIClient) RegisterRouter(router *Router) {
	mc.router = router
}

func (mc *MIDIClient) Run() error {
	defer midi.CloseDriver()

	in, err := midi.FindInPort(mc.InputPort)
	if err != nil {
		return fmt.Errorf("misc.midi.client can't find input port: %s", mc.InputPort)
	}

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		if mc.router != nil {
			mc.router.HandleInput(mc.Id(), msg)
		}
	}, midi.UseSysEx())
	defer stop()

	out, err := midi.FindOutPort(mc.OutputPort)

	if err != nil {
		return fmt.Errorf("misc.midi.client can't find output port: %s", mc.OutputPort)
	}

	send, err := midi.SendTo(out)
	if err != nil {
		return err
	}

	mc.SendFunc = send

	if err != nil {
		return err
	}

	<-mc.router.Context.Done()
	slog.Debug("router context done in module", "id", mc.config.Id)
	return nil
}

func (mc *MIDIClient) Output(payload any) error {
	if mc.SendFunc == nil {
		return fmt.Errorf("misc.midi.client output is not setup")
	}

	payloadMessage, ok := payload.(midi.Message)

	if !ok {
		return fmt.Errorf("misc.midi.client can only ouptut midi.Message")
	}

	return mc.SendFunc(payloadMessage)
}
