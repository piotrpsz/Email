package plain

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"Email/shared/tr"
	"Email/socket"
)

const (
	defaultPlainPort = 110
)

type TCPConnPlain struct {
	conn   *net.TCPConn
	reader io.Reader
	writer io.Writer
}

func ConnectTo(addr string, port int) TCPConnPlain {
	if port == -1 {
		port = defaultPlainPort
	}

	server := fmt.Sprintf("%s:%d", addr, port)

	if addr, err := net.ResolveIPAddr("ip", addr); tr.IsOK(err) {
		tcpAddr := net.TCPAddr{IP: addr.IP, Port: port, Zone: addr.Zone}
		if conn, err := net.DialTCP("tcp", nil, &tcpAddr); tr.IsOK(err) {
			if err := conn.SetKeepAlive(true); tr.IsOK(err) {
				if reader := bufio.NewReader(conn); reader != nil {
					tr.Info("Connect to server (PLAIN): %s", server)
					return TCPConnPlain{conn: conn, reader: reader, writer: conn}
				}
			}
		}
	}

	return TCPConnPlain{}
}

func (s TCPConnPlain) Valid() bool {
	return s.conn != nil
}

func (s TCPConnPlain) Close() {
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
		s.writer = nil
		s.reader = nil
	}
}

func (s TCPConnPlain) Write(data []byte) bool {
	return socket.Write(s.writer, data)
}
func (s TCPConnPlain) ReadByte() (byte, bool) {
	return socket.ReadByte(s.reader)
}
func (s TCPConnPlain) ReadBytes(nbytes int) []byte {
	return socket.ReadBytes(s.reader, nbytes)
}
func (s TCPConnPlain) Read(maxBytesCount int) ([]byte, bool) {
	return socket.Read(s.reader, maxBytesCount)
}
func (s TCPConnPlain) ReadHeader() []byte {
	return socket.ReadHeader(s.reader)
}
func (s TCPConnPlain) ReadBody() []byte {
	return socket.ReadBody(s.reader)
}
