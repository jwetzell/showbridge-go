package module

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/emiago/diago"
	"github.com/emiago/diago/media"
	"github.com/emiago/sipgo"
	"github.com/emiago/sipgo/sip"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type SIPCallServer struct {
	config    config.ModuleConfig
	ctx       context.Context
	router    route.RouteIO
	IP        string
	Port      int
	Transport string
}

type SIPCallMessage struct {
	To string
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "sip.call.server",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {
			params := config.Params
			portNum := 5060

			port, ok := params["port"]
			if ok {
				specificPortNum, ok := port.(float64)

				if !ok {
					return nil, fmt.Errorf("sip.call.server port must be a number")
				}
				portNum = int(specificPortNum)
			}

			ipString := "0.0.0.0"

			ip, ok := params["ip"]
			if ok {

				specificIpString, ok := ip.(string)

				if !ok {
					return nil, fmt.Errorf("sip.call.server ip must be a string")
				}
				ipString = specificIpString
			}

			transportString := "udp"

			transport, ok := params["transport"]
			if ok {

				specificTransportString, ok := transport.(string)

				if !ok {
					return nil, fmt.Errorf("sip.call.server transport must be a string")
				}
				transportString = specificTransportString
			}
			return &SIPCallServer{config: config, ctx: ctx, router: router, IP: ipString, Port: int(portNum), Transport: transportString}, nil
		},
	})
}

func (sds *SIPCallServer) Id() string {
	return sds.config.Id
}

func (sds *SIPCallServer) Type() string {
	return sds.config.Type
}

func (sds *SIPCallServer) Run() error {
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
	slog.Debug("router context done in module", "id", sds.Id())
	return nil
}

func (sds *SIPCallServer) HandleCall(inDialog *diago.DialogServerSession) {
	inDialog.Trying()
	inDialog.Ringing()
	inDialog.Answer()
	sds.router.HandleInput(sds.Id(), SIPCallMessage{
		To: inDialog.ToUser(),
	})
	<-inDialog.Context().Done()
}

func (sds *SIPCallServer) Output(payload any) error {
	return fmt.Errorf("sip.call.server output is not implemented")
}
