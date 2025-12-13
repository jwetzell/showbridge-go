package module

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

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
	UserAgent string
	diag      *diago.Diago
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

			userAgentString := "showbridge"

			userAgent, ok := params["userAgent"]
			if ok {

				specificTransportString, ok := userAgent.(string)

				if !ok {
					return nil, fmt.Errorf("sip.call.server userAgent must be a string")
				}
				userAgentString = specificTransportString
			}
			return &SIPCallServer{config: config, ctx: ctx, router: router, IP: ipString, Port: int(portNum), Transport: transportString, UserAgent: userAgentString}, nil
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
		sipgo.WithUserAgent(sds.UserAgent),
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

	go func() {
		dg.Serve(sds.ctx, func(inDialog *diago.DialogServerSession) {
			sds.HandleCall(inDialog)
		})
	}()

	sds.diag = dg

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

	payloadMsg, ok := payload.(string)
	if !ok {
		return fmt.Errorf("sip.call.server output payload must be of type string")
	}

	if sds.diag == nil {
		return fmt.Errorf("sip.call.server diago not initialized")
	}

	var uri sip.Uri
	err := sip.ParseUri(payloadMsg, &uri)
	if err != nil {
		return fmt.Errorf("sip.call.server output payload is not a valid SIP URI: %v", err)
	}
	outDialog, err := sds.diag.NewDialog(uri, diago.NewDialogOptions{
		Transport: sds.Transport,
	})

	if err != nil {
		return fmt.Errorf("sip.call.server failed to create new dialog: %v", err)
	}
	outDialog.Invite(sds.ctx, diago.InviteClientOptions{})
	outDialog.Ack(sds.ctx)

	// TODO(jwetzell): make this configurable
	// NOTE(jwetzell): wait 5 seconds before hanging up the call
	time.Sleep(5 * time.Second)
	outDialog.Hangup(sds.ctx)
	return nil
}
