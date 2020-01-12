package iface

const (
	CR = byte(0x0d)
	LF = byte(0x0a)
)

type Socket interface {
	Write([]byte) bool
	ReadByte() (byte, bool)
	ReadBytes(n int) []byte
	Read(count int) ([]byte, bool)
	ReadHeader() []byte
	ReadBody() []byte
	Close()
}

type TCPInterface struct {
	socket Socket
}

func New(s Socket) *TCPInterface {
	return &TCPInterface{socket: s}
}

func (i *TCPInterface) Close() {
	i.socket.Close()
}

func (i *TCPInterface) Write(data []byte) bool {
	return i.socket.Write(data)
}

func (i *TCPInterface) readByte() (byte, bool) {
	return i.socket.ReadByte()
}

func (i *TCPInterface) ReadBytes(n int) []byte {
	return i.socket.ReadBytes(n)
}

func (i *TCPInterface) Read(count int) ([]byte, bool) {
	return i.socket.Read(count)
}

func (i *TCPInterface) ReadHeader() []byte {
	return i.socket.ReadHeader()
}

func (i *TCPInterface) ReadBody() []byte {
	return i.socket.ReadBody()
}

/*
func printByte(c byte) {
	if c >= 33 && c <= 126 {
		fmt.Printf("%c", c)
	} else {
		fmt.Printf("(0x%02x)", c)
	}
}

*/
