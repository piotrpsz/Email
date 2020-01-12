package stls

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"

	"Email/shared/tr"
	"Email/socket"
)

const (
	defaultTLSPort = 995
)

type TCPConnTLS struct {
	conn   *tls.Conn
	reader io.Reader
	writer io.Writer
}

func ConnectTo(addr string, port int) TCPConnTLS {
	if port == -1 {
		port = defaultTLSPort
	}

	server := fmt.Sprintf("%s:%d", addr, port)
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	if conn, err := tls.Dial("tcp", server, tlsConfig); tr.IsOK(err) {
		if reader := bufio.NewReader(conn); reader != nil {
			log.Printf("Connected to server (TLS): %s", server)
			return TCPConnTLS{conn: conn, reader: reader, writer: conn}
		}
	}
	return TCPConnTLS{}
}

func (s TCPConnTLS) Valid() bool {
	return s.conn != nil
}

func (s TCPConnTLS) Close() {
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
		s.writer = nil
		s.reader = nil
	}
}

func (s TCPConnTLS) Write(data []byte) bool {
	return socket.Write(s.writer, data)
}
func (s TCPConnTLS) ReadByte() (byte, bool) {
	return socket.ReadByte(s.reader)
}
func (s TCPConnTLS) ReadBytes(nbytes int) []byte {
	return socket.ReadBytes(s.reader, nbytes)
}
func (s TCPConnTLS) Read(maxBytesCount int) ([]byte, bool) {
	return socket.Read(s.reader, maxBytesCount)
}
func (s TCPConnTLS) ReadHeader() []byte {
	return socket.ReadHeader(s.reader)
}
func (s TCPConnTLS) ReadBody() []byte {
	return socket.ReadBody(s.reader)
}
