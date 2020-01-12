package socket

import (
	"bytes"
	"fmt"
	"io"

	"Email/shared/tr"
)

const (
	CR = byte(0x0d)
	LF = byte(0x0a)
)

type Socket interface {
	Valid() bool
	Write([]byte) bool
	ReadByte() (byte, bool)
	ReadBytes(n int) []byte
	Read(count int) ([]byte, bool)
	ReadHeader() []byte
	ReadBody() []byte
	Close()
}

func Write(writer io.Writer, data []byte) bool {
	if _, err := writer.Write(data); tr.IsOK(err) {
		return true
	}
	return false
}

func ReadByte(reader io.Reader) (byte, bool) {
	buffer := []byte{0}
	if n, err := reader.Read(buffer); tr.IsOK(err) {
		if n == 1 {
			return buffer[0], true
		}
	}
	return 0, false
}

func ReadBytes(reader io.Reader, nbytes int) []byte {
	buffer := make([]byte, nbytes)
	if _, err := reader.Read(buffer); tr.IsOK(err) {
		return buffer
	}
	return nil
}

func Read(reader io.Reader, maxBytesCount int) ([]byte, bool) {
	var buffer []byte
	n := 0

	for n < maxBytesCount {
		if c, ok := ReadByte(reader); ok {
			buffer = append(buffer, c)
			n += 1

			if c == LF && n > 1 {
				if buffer[n-2] == CR {
					return buffer[:n-2], true
				}
			}
			continue
		}
		break
	}
	return buffer, false
}

func ReadHeader(reader io.Reader) []byte {
	var buffer []byte
	n := 0

	for {
		if c, ok := ReadByte(reader); ok {
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

func ReadBody(reader io.Reader) []byte {
	var buffer []byte
	n := 0

	for {
		if c, ok := ReadByte(reader); ok {
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

func PrintByte(c byte) {
	if c >= 33 && c <= 126 {
		fmt.Printf("%c", c)
	} else {
		fmt.Printf("(0x%02x)", c)
	}
}
