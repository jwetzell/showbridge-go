<div align="center">

# showbridge (go edition)

[![Coverage](https://github.com/jwetzell/showbridge-go/wiki/coverage.svg)](https://raw.githack.com/wiki/jwetzell/showbridge-go/coverage.html)
Simple protocol router _/s_

</div>

<p align="center">
	<a href="https://github.com/jwetzell/showbridge-go/releases">Releases</a> ·
	<a href="https://docs.showbridge.io">Documentation</a>
</p>

### Supported Protocols

- HTTP
- UDP
- TCP
- [MQTT](https://mqtt.org/)
- [NATS](https://nats.io/)
- [PosiStageNet](https://posistage.net/)
- MIDI (not included in pre-built binaries yet)
- Serial (not included in pre-built binaries yet)
- [OSC](https://opensoundcontrol.stanford.edu/spec-1_0.html)
- [FreeD](https://ptzoptics.com/freed/)
- [SIP](https://en.wikipedia.org/wiki/Session_Initiation_Protocol)
- [Redis](https://redis.io/)

### CLI Usage

```
NAME:
   showbridge - Simple protocol router /s

USAGE:
   showbridge [global options]

GLOBAL OPTIONS:
   --config string      path to config file (default: "./config.yaml") [$SHOWBRIDGE_CONFIG]
   --log-level string   set log level (default: "info") [$SHOWBRIDGE_LOG_LEVEL]
   --log-format string  log format to use (default: "text") [$SHOWBRIDGE_LOG_FORMAT]
   --trace              enable OpenTelemetry tracing [$SHOWBRIDGE_TRACE]
   --help, -h           show help
   --version, -v        print the version
```
