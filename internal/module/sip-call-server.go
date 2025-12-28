package module

import (
	"context"
	"errors"
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
	dg        *diago.Diago
	logger    *slog.Logger
}

type SIPCallMessage struct {
	To string
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "sip.call.server",
		New: func(ctx context.Context, config config.ModuleConfig) (Module, error) {
			params := config.Params
			portNum := 5060

			port, ok := params["port"]
			if ok {
				specificPortNum, ok := port.(float64)

				if !ok {
					return nil, errors.New("sip.call.server port must be a number")
				}
				portNum = int(specificPortNum)
			}

			ipString := "0.0.0.0"

			ip, ok := params["ip"]
			if ok {

				specificIpString, ok := ip.(string)

				if !ok {
					return nil, errors.New("sip.call.server ip must be a string")
				}
				ipString = specificIpString
			}

			transportString := "udp"

			transport, ok := params["transport"]
			if ok {

				specificTransportString, ok := transport.(string)

				if !ok {
					return nil, errors.New("sip.call.server transport must be a string")
				}
				transportString = specificTransportString
			}

			userAgentString := "showbridge"

			userAgent, ok := params["userAgent"]
			if ok {

				specificTransportString, ok := userAgent.(string)

				if !ok {
					return nil, errors.New("sip.call.server userAgent must be a string")
				}
				userAgentString = specificTransportString
			}

			router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

			if !ok {
				return nil, errors.New("sip.call.server unable to get router from context")
			}
			return &SIPCallServer{config: config, ctx: ctx, router: router, IP: ipString, Port: int(portNum), Transport: transportString, UserAgent: userAgentString, logger: CreateLogger(config)}, nil
		},
	})
}

func (scs *SIPCallServer) Id() string {
	return scs.config.Id
}

func (scs *SIPCallServer) Type() string {
	return scs.config.Type
}

func (scs *SIPCallServer) Run() error {
	diagoLogger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	ua, _ := sipgo.NewUA(
		sipgo.WithUserAgent(scs.UserAgent),
		sipgo.WithUserAgentTransportLayerOptions(sip.WithTransportLayerLogger(diagoLogger)),
		sipgo.WithUserAgentTransactionLayerOptions(sip.WithTransactionLayerLogger(diagoLogger)),
	)
	defer ua.Close()

	sip.SetDefaultLogger(diagoLogger)
	media.SetDefaultLogger(diagoLogger)
	dg := diago.NewDiago(ua, diago.WithLogger(diagoLogger), diago.WithTransport(
		diago.Transport{
			Transport: scs.Transport,
			BindHost:  scs.IP,
			BindPort:  scs.Port,
		},
	))

	go func() {
		dg.Serve(scs.ctx, func(inDialog *diago.DialogServerSession) {
			scs.HandleCall(inDialog)
		})
	}()

	scs.dg = dg

	<-scs.ctx.Done()
	scs.logger.Debug("done")
	return nil
}

func (scs *SIPCallServer) HandleCall(inDialog *diago.DialogServerSession) {
	inDialog.Trying()
	inDialog.Ringing()
	inDialog.Answer()
	scs.router.HandleInput(scs.ctx, scs.Id(), SIPCallMessage{
		To: inDialog.ToUser(),
	})
	<-inDialog.Context().Done()
}

func (scs *SIPCallServer) Output(ctx context.Context, payload any) error {

	payloadMsg, ok := payload.(string)
	if !ok {
		return errors.New("sip.call.server output payload must be of type string")
	}

	if scs.dg == nil {
		return errors.New("sip.call.server diago is not initialized")
	}

	var uri sip.Uri
	err := sip.ParseUri(payloadMsg, &uri)
	if err != nil {
		return fmt.Errorf("sip.call.server output payload is not a valid SIP URI: %s", err)
	}
	outDialog, err := scs.dg.NewDialog(uri, diago.NewDialogOptions{
		Transport: scs.Transport,
	})

	if err != nil {
		return fmt.Errorf("sip.call.server failed to create new dialog: %s", err)
	}

	err = outDialog.Invite(scs.ctx, diago.InviteClientOptions{})
	if err != nil {
		return fmt.Errorf("sip.call.server failed to send invite: %s", err)
	}

	err = outDialog.Ack(scs.ctx)
	if err != nil {
		return fmt.Errorf("sip.call.server failed to send ack: %s", err)
	}
	// TODO(jwetzell): make this configurable
	// NOTE(jwetzell): wait 5 seconds before hanging up the call
	time.Sleep(5 * time.Second)
	err = outDialog.Hangup(scs.ctx)
	if err != nil {
		return fmt.Errorf("sip.call.server failed to hangup call: %s", err)
	}
	return nil
}
