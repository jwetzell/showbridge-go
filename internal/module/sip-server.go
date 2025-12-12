package module

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/emiago/diago"
	"github.com/emiago/sipgo"
	"github.com/emiago/sipgo/sip"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type SIPServer struct {
	config config.ModuleConfig
	ctx    context.Context
	router route.RouteIO
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.sip.server",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {

			return &SIPServer{config: config, ctx: ctx, router: router}, nil
		},
	})
}

func (ss *SIPServer) Id() string {
	return ss.config.Id
}

func (ss *SIPServer) Type() string {
	return ss.config.Type
}

func (ss *SIPServer) Run() error {

	ua, _ := sipgo.NewUA()

	diagoLogger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	sip.SetDefaultLogger(diagoLogger)
	tu := diago.NewDiago(ua, diago.WithLogger(diagoLogger), diago.WithTransport(
		diago.Transport{
			Transport: "udp",
			BindHost:  "0.0.0.0",
			BindPort:  5060,
		},
	))

	err := tu.Serve(ss.ctx, func(inDialog *diago.DialogServerSession) {
		ss.HandleCall(inDialog)
	})

	if err != nil {
		return err
	}

	<-ss.ctx.Done()
	slog.Debug("router context done in module", "id", ss.Id())
	return nil
}

func (ss *SIPServer) HandleCall(inDialog *diago.DialogServerSession) error {
	inDialog.Trying()
	inDialog.Ringing()
	inDialog.Answer()

	reader := inDialog.AudioReaderDTMF()
	userString := ""
	return reader.Listen(func(dtmf rune) error {
		if dtmf == '#' {
			if ss.router != nil {
				ss.router.HandleInput(ss.Id(), userString)
			}
			userString = ""
		} else {
			userString += string(dtmf)
		}
		return nil
	}, 5*time.Second)
}

func (ss *SIPServer) Output(payload any) error {
	return fmt.Errorf("net.sip.server output is not implemented")
}
