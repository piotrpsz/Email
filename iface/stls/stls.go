package stls

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"

	"Email/iface"
	"Email/shared/tr"
)

const (
	defaultTLSPort = 995
)

type TLSSocket struct {
	writer *tls.Conn
	reader *bufio.Reader
}

func ConnectTo(addr string, port int) (TLSSocket, bool) {
	tr.In()
	defer tr.Out()

	s := TLSSocket{}
	if port == -1 {
		port = defaultTLSPort
	}
	server := fmt.Sprintf("%s:%d", addr, port)
	tr.Info("Connect to server: %s", server)

	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	if conn, err := tls.Dial("tcp", server, tlsConfig); tr.IsOK(err) {
		if reader := bufio.NewReader(conn); reader != nil {
			s.reader = reader
			s.writer = conn
			return s, true
		}
	}
	return s, false
}

func (s TLSSocket) Close() {
	if s.writer != nil {
		s.writer.Close()
		s.writer = nil
		s.reader = nil
	}
}

func (s TLSSocket) Write(data []byte) bool {
	if _, err := s.writer.Write(data); tr.IsOK(err) {
		return true
	}
	return false
}

func (s TLSSocket) ReadByte() (byte, bool) {
	buffer := []byte{0}
	if n, err := s.reader.Read(buffer); tr.IsOK(err) {
		if n == 1 {
			return buffer[0], true
		}
	}
	return 0, false
}

func (s TLSSocket) ReadBytes(n int) []byte {
	buffer := make([]byte, n)
	if _, err := s.reader.Read(buffer); tr.IsOK(err) {
		return buffer
	}
	return nil
}

func (s TLSSocket) Read(count int) ([]byte, bool) {
	var buffer []byte
	n := 0

	for n < count {
		if c, ok := s.ReadByte(); ok {
			buffer = append(buffer, c)
			n += 1

			if c == iface.LF && n > 1 {
				if buffer[n-2] == iface.CR {
					return buffer[:n-2], true
				}
			}
			continue
		}
		break
	}
	return buffer, false
}

func (s TLSSocket) ReadHeader() []byte {
	var buffer []byte
	n := 0

	for {
		if c, ok := s.ReadByte(); ok {
			buffer = append(buffer, c)
			n += 1

			if bytes.HasSuffix(buffer, []byte("\r\n\r\n")) {
				return buffer
			}
			continue
		}
		break
	}

	return buffer
}

func (s TLSSocket) ReadBody() []byte {
	var buffer []byte
	n := 0

	for {
		if c, ok := s.ReadByte(); ok {
			buffer = append(buffer, c)
			n += 1

			if bytes.HasSuffix(buffer, []byte("\r\n.\r\n")) {
				return buffer
			}
			continue
		}
		break
	}

	return nil
}
