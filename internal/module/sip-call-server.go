package module

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/emiago/diago"
	"github.com/emiago/diago/media"
	"github.com/emiago/sipgo"
	"github.com/emiago/sipgo/sip"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
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
	cancel    context.CancelFunc
}

type SIPCallMessage struct {
	To string
}

type SIPCall struct {
	inDialog *diago.DialogServerSession
	lock     sync.Mutex
}

type sipCallContextKey string

func init() {
	RegisterModule(ModuleRegistration{
		Type: "sip.call.server",
		New: func(config config.ModuleConfig) (Module, error) {
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

			return &SIPCallServer{config: config, IP: ipString, Port: int(portNum), Transport: transportString, UserAgent: userAgentString, logger: CreateLogger(config)}, nil
		},
	})
}

func (scs *SIPCallServer) Id() string {
	return scs.config.Id
}

func (scs *SIPCallServer) Type() string {
	return scs.config.Type
}

func (scs *SIPCallServer) Run(ctx context.Context) error {
	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return errors.New("sip.call.server unable to get router from context")
	}
	scs.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	scs.ctx = moduleContext
	scs.cancel = cancel

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

	dialogContext := context.WithValue(scs.ctx, sipCallContextKey("call"), &SIPCall{
		inDialog: inDialog,
	})
	scs.router.HandleInput(dialogContext, scs.Id(), SIPCallMessage{
		To: inDialog.ToUser(),
	})
}

func (scs *SIPCallServer) Output(ctx context.Context, payload any) error {

	call, ok := ctx.Value(sipCallContextKey("call")).(*SIPCall)

	if !ok {
		return errors.New("sip.call.server output must originate from sip.call.server input")
	}

	gotLock := call.lock.TryLock()

	if !gotLock {
		return errors.New("sip.call.server call is already locked")
	}

	if call.inDialog.LoadState() == sip.DialogStateEnded {
		return errors.New("sip.call.server inDialog already ended")
	}

	payloadDTMFResponse, ok := payload.(processor.SipDTMFResponse)

	if ok {
		dtmfWriter := call.inDialog.AudioWriterDTMF()

		time.Sleep(time.Millisecond * time.Duration(payloadDTMFResponse.PreWait))
		for i, dtmfRune := range payloadDTMFResponse.Digits {
			err := dtmfWriter.WriteDTMF(dtmfRune)

			if err != nil {
				return fmt.Errorf("sip.dtmf.server error output dtmf digit at index %d", i)
			}
		}
		time.Sleep(time.Millisecond * time.Duration(payloadDTMFResponse.PreWait))
		return nil
	}

	payloadAudioFileResponse, ok := payload.(processor.SipAudioFileResponse)

	if ok {
		audioFile, err := os.Open(payloadAudioFileResponse.AudioFile)
		if err != nil {
			return err
		}
		defer audioFile.Close()

		playback, err := call.inDialog.PlaybackCreate()

		if err != nil {
			return err
		}

		time.Sleep(time.Millisecond * time.Duration(payloadAudioFileResponse.PreWait))

		_, err = playback.Play(audioFile, "audio/wav")

		time.Sleep(time.Millisecond * time.Duration(payloadAudioFileResponse.PostWait))

		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("sip.dtmf.server can only output SipDTMFResponse or SipAudioFileResponse")
}

func (scs *SIPCallServer) Stop() {
	scs.cancel()
}
