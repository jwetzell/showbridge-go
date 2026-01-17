package module

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
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

type SIPDTMFCall struct {
	inDialog *diago.DialogServerSession
	lock     sync.Mutex
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "sip.dtmf.server",
		New: func(config config.ModuleConfig) (Module, error) {
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
			return &SIPDTMFServer{config: config, IP: ipString, Port: int(portNum), Transport: transportString, Separator: separatorString, logger: CreateLogger(config)}, nil
		},
	})
}

func (sds *SIPDTMFServer) Id() string {
	return sds.config.Id
}

func (sds *SIPDTMFServer) Type() string {
	return sds.config.Type
}

func (sds *SIPDTMFServer) Run(ctx context.Context) error {
	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return errors.New("sip.dtmf.server unable to get router from context")
	}
	sds.router = router
	sds.ctx = ctx

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
				dialogContext := context.WithValue(sds.ctx, sipCallContextKey("call"), &SIPDTMFCall{
					inDialog: inDialog,
				})
				sds.router.HandleInput(dialogContext, sds.Id(), SIPDTMFMessage{
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
	call, ok := ctx.Value(sipCallContextKey("call")).(*SIPDTMFCall)

	if !ok {
		return errors.New("sip.dtmf.server output must originate from sip.dtmf.server input")
	}

	gotLock := call.lock.TryLock()

	if !gotLock {
		return errors.New("sip.dtmf.server call is already locked")
	}

	if call.inDialog.LoadState() == sip.DialogStateEnded {
		return errors.New("sip.dtmf.server inDialog already ended")
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
