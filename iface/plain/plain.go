package plain

import (
	"bufio"
	"bytes"
	"fmt"
	"net"

	"Email/iface"
	"Email/shared/tr"
)

const (
	defaultPlainPort = 110
)

type PlainSocket struct {
	writer *net.TCPConn
	reader *bufio.Reader
}

func ConnectTo(addr string, port int) (PlainSocket, bool) {

	s := PlainSocket{}
	if port == -1 {
		port = defaultPlainPort
	}
	server := fmt.Sprintf("%s:%d", addr, port)
	tr.Info("Connect to server: %s", server)

	if addr, err := net.ResolveIPAddr("ip", addr); tr.IsOK(err) {
		tcpAddr := net.TCPAddr{IP: addr.IP, Port: port, Zone: addr.Zone}
		if conn, err := net.DialTCP("tcp", nil, &tcpAddr); tr.IsOK(err) {
			if err := conn.SetKeepAlive(true); tr.IsOK(err) {
				if reader := bufio.NewReader(conn); reader != nil {
					s.reader = reader
					s.writer = conn
					return s, true
				}
			}
		}
	}

	return s, false
}

func (s PlainSocket) Close() {
	if s.writer != nil {
		s.writer.Close()
		s.writer = nil
		s.reader = nil
	}
}

func (s PlainSocket) Write(data []byte) bool {
	if _, err := s.writer.Write(data); tr.IsOK(err) {
		return true
	}
	return false
}

func (s PlainSocket) ReadByte() (byte, bool) {
	buffer := []byte{0}
	if n, err := s.reader.Read(buffer); tr.IsOK(err) {
		if n == 1 {
			return buffer[0], true
		}
	}
	return 0, false
}

func (s PlainSocket) ReadBytes(n int) []byte {
	buffer := make([]byte, n)
	if _, err := s.reader.Read(buffer); tr.IsOK(err) {
		return buffer
	}
	return nil
}

func (s PlainSocket) Read(count int) ([]byte, bool) {
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

func (s PlainSocket) ReadHeader() []byte {
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

func (s PlainSocket) ReadBody() []byte {
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
