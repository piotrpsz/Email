package command

import (
	"Email/iface"
	"Email/pop3/response"
)

const (
	crlf      = "\r\n"
	separator = " "
)

func Send(ifc *iface.TCPInterface, cmd string, args []string) *response.Response {
	buffer := cmd
	for _, item := range args {
		buffer += separator
		buffer += item
	}
	buffer += crlf

	if ifc.Write([]byte(buffer)) {
		return response.Read(ifc)
	}
	return nil
}
