package showbridge

import (
	"context"
	"fmt"
	"net"
	"time"
)

type TCPClient struct {
	Host string
	Port uint16
}

func init() {
	RegisterProtocol(ProtocolRegistration{
		Type: "tcp.client",
		New: func(params map[string]any) (Protocol, error) {

			host, ok := params["host"]

			if !ok {
				return nil, fmt.Errorf("tcp client requires a host parameter")
			}

			hostString, ok := host.(string)

			if !ok {
				return nil, fmt.Errorf("tcp client host must be uint16")
			}

			port, ok := params["port"]
			if !ok {
				return nil, fmt.Errorf("tcp client requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, fmt.Errorf("tcp client port must be uint16")
			}

			return TCPClient{Host: hostString, Port: uint16(portNum)}, nil
		},
	})
}

func (ts TCPClient) Run(ctx context.Context) error {
	for {
		client, err := net.Dial("tcp", fmt.Sprintf(":%d", ts.Port))
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 2)
			continue
		}

		buffer := make([]byte, 1024)
		select {
		case <-ctx.Done():
			return nil
		default:
		READ:
			for {
				select {
				case <-ctx.Done():
					return nil
				default:
					byteCount, err := client.Read(buffer)

					if err != nil {
						fmt.Println("connection closed")
						break READ
					}

					if byteCount > 0 {
						fmt.Println(buffer[0:byteCount])
					}
				}

			}
		}

	}

}
