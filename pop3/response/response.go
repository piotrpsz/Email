package response

import (
	"bytes"

	"Email/socket"
)

const (
	maxResponseSize = 512
	okPrefix        = "+OK"
	errPrefix       = "-ERR"
)

var (
	okBytes, errBytes []byte
)

type Response struct {
	Status string
	Data   string
}

func init() {
	okBytes = []byte(okPrefix)
	errBytes = []byte(errPrefix)
}

func Read(sck socket.Socket) *Response {
	if data, ok := sck.Read(maxResponseSize); ok {
		if bytes.HasPrefix(data, okBytes) {
			return &Response{
				Status: okPrefix,
				Data:   string(bytes.TrimSpace(data[len(okBytes):])),
			}
		}

		if bytes.HasPrefix(data, errBytes) {
			return &Response{
				Status: errPrefix,
				Data:   string(bytes.TrimSpace(data[len(okBytes):])),
			}
		}

	}
	return nil
}

func (r *Response) IsOK() bool {
	return r.Status == okPrefix
}
