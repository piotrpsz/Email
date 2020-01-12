package pop3

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"Email/pop3/command"
	"Email/pop3/response"
	"Email/shared/tr"
	"Email/socket"
	"Email/socket/plain"
	"Email/socket/stls"
)

type Security uint8

const (
	Plain Security = iota
	TLS
)

const (
	CR    = byte(0x0d)
	LF    = byte(0x0a)
	HTAB  = byte(0x09)
	SPACE = byte(0x20)
)

type POP3 struct {
	socket socket.Socket
}

type ListItem struct {
	ID   int
	Size int
}

func (i ListItem) String() string {
	return fmt.Sprintf("ListItem (id: %d, size: %d)", i.ID, i.Size)
}

type HeaderItem struct {
	Key   string
	Value string
}

func (h HeaderItem) String() string {
	return fmt.Sprintf("(%s : %s)", h.Key, h.Value)
}

func New(server string, port int, security Security) *POP3 {
	if security == TLS {
		if tlsSocket := stls.ConnectTo(server, port); tlsSocket.Valid() {
			return &POP3{socket: tlsSocket}
		}
		return nil
	}
	if security == Plain {
		if plainSocket := plain.ConnectTo(server, port); plainSocket.Valid() {
			return &POP3{socket: plainSocket}
		}
	}
	return nil
}

func (p *POP3) Close() {
	if p.socket != nil {
		p.socket.Close()
		p.socket = nil
	}
}

func (p *POP3) Quit() *response.Response {
	return command.Send(p.socket, "QUIT", nil)
}

func (p *POP3) Auth(user, password string) (*response.Response, bool) {
	var resp *response.Response

	if resp = command.Send(p.socket, "USER", []string{user}); resp != nil && resp.IsOK() {
		if resp = command.Send(p.socket, "PASS", []string{password}); resp != nil && resp.IsOK() {
			return resp, true
		}
	}

	return resp, false
}

func (p *POP3) Stat() (*response.Response, bool) {
	var resp *response.Response

	if resp = command.Send(p.socket, "STAT", nil); resp != nil && resp.IsOK() {
		return resp, true
	}
	return resp, false
}

/// List - read information about E-mails.
/// Returns the array of ListInfo (can be empty).
func (p *POP3) List() []ListItem {
	var resp *response.Response

	if resp = command.Send(p.socket, "LIST", nil); resp != nil {
		if resp.IsOK() {
			return p.ListRead()
		}
	}
	return nil
}

func (p *POP3) ListRead() []ListItem {
	var result []ListItem

	for {
		if data, ok := p.socket.Read(512); ok {
			text := string(data)
			if text == "." {
				break
			}
			if items := splitToItems(text, " "); len(items) == 2 {
				if id, err := strconv.Atoi(items[0]); tr.IsOK(err) {
					if size, err := strconv.Atoi(items[1]); tr.IsOK(err) {
						item := ListItem{ID: id, Size: size}
						result = append(result, item)
					}
				}
			}
			continue
		}
		break
	}
	return result
}

func (p *POP3) Read() *response.Response {
	return response.Read(p.socket)
}

func splitToItems(text, sep string) []string {
	var data []string

	if items := strings.Split(text, sep); len(items) > 0 {
		for _, item := range items {
			data = append(data, strings.TrimSpace(item))
		}
	}

	return data
}

/********************************************************************
*                                                                   *
*                         E - M A I L                               *
*                                                                   *
********************************************************************/

func (p *POP3) ReadEmail(item ListItem) bool {
	fmt.Println(item)
	number := fmt.Sprintf("%d", item.ID)
	if resp := command.Send(p.socket, "RETR", []string{number}); resp != nil && resp.IsOK() {
		if data := p.socket.ReadHeader(); data != nil {
			fmt.Println("Size:", len(data))
			if hdr := p.parseHeader(data); hdr != nil {
				for _, i := range hdr {
					fmt.Printf("%s:%s\n", i.Key, i.Value)
				}
				fmt.Println()
				if data := p.socket.ReadBody(); data != nil {
					fmt.Println(string(data))
					fmt.Println("Size:", len(data))
				} else {
					fmt.Println()
					fmt.Println("DUPA")
				}
			}

			return true
		}
	}
	fmt.Println("ERROR")
	return false
}

func (p *POP3) parseHeader(data []byte) []HeaderItem {
	var header []HeaderItem
	var acc []byte

	buffer := bytes.Split(data, []byte{CR, LF})
	acc = append(acc, buffer[0]...)
	//printAsText(buffer[0])
	i := 1

	for i < len(buffer) {
		//printAsText(buffer[i])
		if len(buffer[i]) > 0 {
			if buffer[i][0] == HTAB || buffer[i][0] == SPACE {
				if buffer[i][0] == HTAB {
					//acc = append(acc, []byte{' '}...)
				}
				acc = append(acc, bytes.TrimSpace(buffer[i])...)
				i++
				continue
			}

			if len(acc) > 0 {
				if key, value, ok := keyAndValue(acc); ok {
					header = append(header, HeaderItem{Key: key, Value: value})
				}
			}
			acc = []byte{}
			acc = append(acc, buffer[i]...)
		}
		i++
	}

	if len(acc) > 0 {
		if key, value, ok := keyAndValue(acc); ok {
			header = append(header, HeaderItem{Key: key, Value: value})
		}
	}

	return header
}

func keyAndValue(data []byte) (string, string, bool) {
	if idx := indexOf(data, ':'); idx != -1 {
		keyBytes := bytes.TrimSpace(data[:idx])
		valueBytes := bytes.TrimSpace(data[idx+1:])
		return string(keyBytes), string(valueBytes), true

	}
	return "", "", false
}

func printAsText(data []byte) {
	for _, c := range data {
		if c >= 33 && c <= 126 {
			fmt.Printf("%c", c)
		} else {
			fmt.Printf("(0x%02x)", c)
		}
	}
	fmt.Println()
}

func isAlpha(c byte) bool {
	return (c >= 65 && c <= 90) || (c >= 97 && c <= 122)
}

func isDigit(c byte) bool {
	return c >= 48 && c <= 57
}

func isOtherValid(c byte) bool {
	if c == '=' {
		return true
	}
	if c == '+' || c == '-' {
		return true
	}
	if c == ':' {
		return true
	}
	if c == '(' || c == ')' {
		return true
	}
	if c == '[' || c == ']' {
		return true
	}
	if c == '"' {
		return true
	}
	if c == '@' {
		return true
	}

	return false
}

func indexOf(data []byte, c byte) int {
	for i, item := range data {
		if item == c {
			return i
		}
	}
	return -1
}
