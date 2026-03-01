package processor

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"text/template"

	freeD "github.com/jwetzell/free-d-go"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type FreeDCreate struct {
	config config.ProcessorConfig
	Id     *template.Template
	Pan    *template.Template
	Tilt   *template.Template
	Roll   *template.Template
	PosX   *template.Template
	PosY   *template.Template
	PosZ   *template.Template
	Zoom   *template.Template
	Focus  *template.Template
}

func (fc *FreeDCreate) Process(ctx context.Context, payload any) (any, error) {

	var idBuffer bytes.Buffer
	err := fc.Id.Execute(&idBuffer, payload)

	if err != nil {
		return nil, err
	}

	idString := idBuffer.String()

	idNum, err := strconv.ParseUint(idString, 10, 8)

	if err != nil {
		return nil, err
	}

	var panBuffer bytes.Buffer
	err = fc.Pan.Execute(&panBuffer, payload)

	if err != nil {
		return nil, err
	}

	panString := panBuffer.String()

	panNum, err := strconv.ParseFloat(panString, 32)

	if err != nil {
		return nil, err
	}

	var tiltBuffer bytes.Buffer
	err = fc.Tilt.Execute(&tiltBuffer, payload)

	if err != nil {
		return nil, err
	}

	tiltString := tiltBuffer.String()

	tiltNum, err := strconv.ParseFloat(tiltString, 32)

	if err != nil {
		return nil, err
	}

	var rollBuffer bytes.Buffer
	err = fc.Tilt.Execute(&rollBuffer, payload)

	if err != nil {
		return nil, err
	}

	rollString := rollBuffer.String()

	rollNum, err := strconv.ParseFloat(rollString, 32)

	if err != nil {
		return nil, err
	}

	var posXBuffer bytes.Buffer
	err = fc.PosX.Execute(&posXBuffer, payload)

	if err != nil {
		return nil, err
	}

	posXString := posXBuffer.String()

	posXNum, err := strconv.ParseFloat(posXString, 32)

	if err != nil {
		return nil, err
	}

	var posYBuffer bytes.Buffer
	err = fc.PosY.Execute(&posYBuffer, payload)

	if err != nil {
		return nil, err
	}

	posYString := posYBuffer.String()

	posYNum, err := strconv.ParseFloat(posYString, 32)

	if err != nil {
		return nil, err
	}

	var posZBuffer bytes.Buffer
	err = fc.PosZ.Execute(&posZBuffer, payload)

	if err != nil {
		return nil, err
	}

	posZString := posZBuffer.String()

	posZNum, err := strconv.ParseFloat(posZString, 32)

	if err != nil {
		return nil, err
	}

	var zoomBuffer bytes.Buffer
	err = fc.Zoom.Execute(&zoomBuffer, payload)

	if err != nil {
		return nil, err
	}

	zoomString := zoomBuffer.String()

	zoomNum, err := strconv.ParseInt(zoomString, 10, 32)

	if err != nil {
		return nil, err
	}

	var focusBuffer bytes.Buffer
	err = fc.Zoom.Execute(&focusBuffer, payload)

	if err != nil {
		return nil, err
	}

	focusString := focusBuffer.String()

	focusNum, err := strconv.ParseInt(focusString, 10, 32)

	if err != nil {
		return nil, err
	}

	payloadMessage := freeD.FreeDPosition{
		ID:    uint8(idNum),
		Pan:   float32(panNum),
		Tilt:  float32(tiltNum),
		Roll:  float32(rollNum),
		PosX:  float32(posXNum),
		PosY:  float32(posYNum),
		PosZ:  float32(posZNum),
		Zoom:  int32(zoomNum),
		Focus: int32(focusNum),
	}

	return payloadMessage, nil
}

func (fc *FreeDCreate) Type() string {
	return fc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "freed.create",
		New: func(config config.ProcessorConfig) (Processor, error) {

			// TODO(jwetzell): make some params optional
			params := config.Params
			idString, err := params.GetString("id")
			if err != nil {
				return nil, fmt.Errorf("freed.create id error: %w", err)
			}

			idTemplate, err := template.New("id").Parse(idString)

			if err != nil {
				return nil, err
			}

			panString, err := params.GetString("pan")
			if err != nil {
				return nil, fmt.Errorf("freed.create pan error: %w", err)
			}

			panTemplate, err := template.New("pan").Parse(panString)

			tiltString, err := params.GetString("tilt")
			if err != nil {
				return nil, fmt.Errorf("freed.create tilt error: %w", err)
			}

			tiltTemplate, err := template.New("tilt").Parse(tiltString)

			rollString, err := params.GetString("roll")
			if err != nil {
				return nil, fmt.Errorf("freed.create roll error: %w", err)
			}

			rollTemplate, err := template.New("roll").Parse(rollString)

			if err != nil {
				return nil, err
			}

			posXString, err := params.GetString("posX")
			if err != nil {
				return nil, fmt.Errorf("freed.create posX error: %w", err)
			}

			posXTemplate, err := template.New("posX").Parse(posXString)

			if err != nil {
				return nil, err
			}

			posYString, err := params.GetString("posY")
			if err != nil {
				return nil, fmt.Errorf("freed.create posY error: %w", err)
			}

			posYTemplate, err := template.New("posY").Parse(posYString)

			if err != nil {
				return nil, err
			}

			posZString, err := params.GetString("posZ")
			if err != nil {
				return nil, fmt.Errorf("freed.create posZ error: %w", err)
			}

			posZTemplate, err := template.New("posZ").Parse(posZString)

			if err != nil {
				return nil, err
			}

			zoomString, err := params.GetString("zoom")
			if err != nil {
				return nil, fmt.Errorf("freed.create zoom error: %w", err)
			}

			zoomTemplate, err := template.New("zoom").Parse(zoomString)

			focusString, err := params.GetString("focus")
			if err != nil {
				return nil, fmt.Errorf("freed.create focus error: %w", err)
			}

			focusTemplate, err := template.New("focus").Parse(focusString)

			return &FreeDCreate{
				config: config,
				Id:     idTemplate,
				Pan:    panTemplate,
				Tilt:   tiltTemplate,
				Roll:   rollTemplate,
				PosX:   posXTemplate,
				PosY:   posYTemplate,
				PosZ:   posZTemplate,
				Zoom:   zoomTemplate,
				Focus:  focusTemplate,
			}, nil
		},
	})
}
