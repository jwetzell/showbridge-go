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
			id, ok := params["id"]

			if !ok {
				return nil, fmt.Errorf("freed.create requires an id parameter")
			}

			idString, ok := id.(string)

			if !ok {
				return nil, fmt.Errorf("freed.create id must be a string")
			}

			idTemplate, err := template.New("id").Parse(idString)

			if err != nil {
				return nil, err
			}

			pan, ok := params["pan"]

			if !ok {
				return nil, fmt.Errorf("freed.create requires an pan parameter")
			}

			panString, ok := pan.(string)

			if !ok {
				return nil, fmt.Errorf("freed.create pan must be a string")
			}

			panTemplate, err := template.New("pan").Parse(panString)

			tilt, ok := params["tilt"]

			if !ok {
				return nil, fmt.Errorf("freed.create requires an tilt parameter")
			}

			tiltString, ok := tilt.(string)

			if !ok {
				return nil, fmt.Errorf("freed.create tilt must be a string")
			}

			tiltTemplate, err := template.New("tilt").Parse(tiltString)

			roll, ok := params["roll"]

			if !ok {
				return nil, fmt.Errorf("freed.create requires an roll parameter")
			}

			rollString, ok := roll.(string)

			if !ok {
				return nil, fmt.Errorf("freed.create roll must be a string")
			}

			rollTemplate, err := template.New("roll").Parse(rollString)

			if err != nil {
				return nil, err
			}

			posX, ok := params["posX"]

			if !ok {
				return nil, fmt.Errorf("freed.create requires a posX parameter")
			}

			posXString, ok := posX.(string)

			if !ok {
				return nil, fmt.Errorf("freed.create posX must be a string")
			}

			posXTemplate, err := template.New("posX").Parse(posXString)

			if err != nil {
				return nil, err
			}

			posY, ok := params["posY"]

			if !ok {
				return nil, fmt.Errorf("freed.create requires a posY parameter")
			}

			posYString, ok := posY.(string)

			if !ok {
				return nil, fmt.Errorf("freed.create posY must be a string")
			}

			posYTemplate, err := template.New("posY").Parse(posYString)

			if err != nil {
				return nil, err
			}

			posZ, ok := params["posZ"]

			if !ok {
				return nil, fmt.Errorf("freed.create requires a posZ parameter")
			}

			posZString, ok := posZ.(string)

			if !ok {
				return nil, fmt.Errorf("freed.create posZ must be a string")
			}

			posZTemplate, err := template.New("posZ").Parse(posZString)

			if err != nil {
				return nil, err
			}

			zoom, ok := params["zoom"]

			if !ok {
				return nil, fmt.Errorf("freed.create requires an zoom parameter")
			}

			zoomString, ok := zoom.(string)

			if !ok {
				return nil, fmt.Errorf("freed.create zoom must be a string")
			}

			zoomTemplate, err := template.New("zoom").Parse(zoomString)

			focus, ok := params["focus"]

			if !ok {
				return nil, fmt.Errorf("freed.create requires an focus parameter")
			}

			focusString, ok := focus.(string)

			if !ok {
				return nil, fmt.Errorf("freed.create focus must be a string")
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
