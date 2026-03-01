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
	UserAgent string
	Separator string
	logger    *slog.Logger
	cancel    context.CancelFunc
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
		New: func(moduleConfig config.ModuleConfig) (Module, error) {
			params := moduleConfig.Params

			portNum, err := params.GetInt("port")
			if err != nil {

				if errors.Is(err, config.ErrParamNotFound) {
					portNum = 5060
				} else {
					return nil, fmt.Errorf("sip.dtmf.server port error: %w", err)
				}
			}

			ipString, err := params.GetString("ip")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					ipString = "0.0.0.0"
				} else {
					return nil, fmt.Errorf("sip.dtmf.server ip error: %w", err)
				}
			}

			transportString, err := params.GetString("transport")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					transportString = "udp"
				} else {
					return nil, fmt.Errorf("sip.dtmf.server transport error: %w", err)
				}
			}

			userAgentString, err := params.GetString("userAgent")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					userAgentString = "showbridge"
				} else {
					return nil, fmt.Errorf("sip.dtmf.server userAgent error: %w", err)
				}
			}

			separatorString, err := params.GetString("separator")
			if err != nil {
				return nil, fmt.Errorf("sip.dtmf.server separator error: %w", err)
			}

			if len(separatorString) != 1 {
				return nil, errors.New("sip.dtmf.server separator must be a single character")
			}

			if !strings.ContainsRune("0123456789*#ABCD", rune(separatorString[0])) {
				return nil, errors.New("sip.dtmf.server separator must be a valid DTMF character")
			}
			return &SIPDTMFServer{config: moduleConfig, IP: ipString, Port: int(portNum), Transport: transportString, UserAgent: userAgentString, Separator: separatorString, logger: CreateLogger(moduleConfig)}, nil
		},
	})
}

func (sds *SIPDTMFServer) Id() string {
	return sds.config.Id
}

func (sds *SIPDTMFServer) Type() string {
	return sds.config.Type
}

func (sds *SIPDTMFServer) Start(ctx context.Context) error {
	sds.logger.Debug("running")
	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return errors.New("sip.dtmf.server unable to get router from context")
	}
	sds.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	sds.ctx = moduleContext
	sds.cancel = cancel

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
		time.Sleep(time.Millisecond * time.Duration(payloadDTMFResponse.PostWait))
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

func (sds *SIPDTMFServer) Stop() {
	sds.cancel()
}
