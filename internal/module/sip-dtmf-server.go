package module

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/emiago/diago"
	"github.com/emiago/diago/media"
	"github.com/emiago/sipgo"
	"github.com/emiago/sipgo/sip"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type SIPDTMFServer struct {
	config    config.ModuleConfig
	ctx       context.Context
	router    route.RouteIO
	IP        string
	Port      int
	Transport string
	Separator string
	logger    *slog.Logger
}

type SIPDTMFMessage struct {
	To     string
	Digits string
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "sip.dtmf.server",
		New: func(ctx context.Context, config config.ModuleConfig) (Module, error) {
			params := config.Params
			portNum := 5060

			port, ok := params["port"]
			if ok {
				specificPortNum, ok := port.(float64)

				if !ok {
					return nil, errors.New("sip.dtmf.server port must be a number")
				}
				portNum = int(specificPortNum)
			}

			ipString := "0.0.0.0"

			ip, ok := params["ip"]
			if ok {

				specificIpString, ok := ip.(string)

				if !ok {
					return nil, errors.New("sip.dtmf.server ip must be a string")
				}
				ipString = specificIpString
			}

			transportString := "udp"

			transport, ok := params["transport"]
			if ok {

				specificTransportString, ok := transport.(string)

				if !ok {
					return nil, errors.New("sip.dtmf.server transport must be a string")
				}
				transportString = specificTransportString
			}

			separator, ok := params["separator"]
			if !ok {
				return nil, errors.New("sip.dtmf.server requires a separator parameter")
			}
			separatorString, ok := separator.(string)
			if !ok {
				return nil, errors.New("sip.dtmf.server separator must be a string")
			}

			if len(separatorString) != 1 {
				return nil, errors.New("sip.dtmf.server separator must be a single character")
			}

			if !strings.ContainsRune("0123456789*#ABCD", rune(separatorString[0])) {
				return nil, errors.New("sip.dtmf.server separator must be a valid DTMF character")
			}
			router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

			if !ok {
				return nil, errors.New("sip.dtmf.server unable to get router from context")
			}
			return &SIPDTMFServer{config: config, ctx: ctx, router: router, IP: ipString, Port: int(portNum), Transport: transportString, Separator: separatorString, logger: CreateLogger(config)}, nil
		},
	})
}

func (sds *SIPDTMFServer) Id() string {
	return sds.config.Id
}

func (sds *SIPDTMFServer) Type() string {
	return sds.config.Type
}

func (sds *SIPDTMFServer) Run() error {
	diagoLogger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	ua, _ := sipgo.NewUA(
		sipgo.WithUserAgentTransportLayerOptions(sip.WithTransportLayerLogger(diagoLogger)),
		sipgo.WithUserAgentTransactionLayerOptions(sip.WithTransactionLayerLogger(diagoLogger)),
	)
	defer ua.Close()

	sip.SetDefaultLogger(diagoLogger)
	media.SetDefaultLogger(diagoLogger)
	dg := diago.NewDiago(ua, diago.WithLogger(diagoLogger), diago.WithTransport(
		diago.Transport{
			Transport: sds.Transport,
			BindHost:  sds.IP,
			BindPort:  sds.Port,
		},
	))

	err := dg.Serve(sds.ctx, func(inDialog *diago.DialogServerSession) {
		sds.HandleCall(inDialog)
	})

	if err != nil {
		return err
	}

	<-sds.ctx.Done()
	sds.logger.Debug("done")
	return nil
}

func (sds *SIPDTMFServer) HandleCall(inDialog *diago.DialogServerSession) error {
	inDialog.Trying()
	inDialog.Ringing()
	inDialog.Answer()

	reader := inDialog.AudioReaderDTMF()
	userString := ""
	return reader.Listen(func(dtmf rune) error {
		if dtmf == rune(sds.Separator[0]) {
			if sds.router != nil {
				sds.router.HandleInput(sds.ctx, sds.Id(), SIPDTMFMessage{
					To:     inDialog.ToUser(),
					Digits: userString,
				})
			}
			userString = ""
		} else {
			userString += string(dtmf)
		}
		return nil
	}, 5*time.Second)
}

func (sds *SIPDTMFServer) Output(ctx context.Context, payload any) error {
	return errors.New("sip.dtmf.server output is not implemented")
}
