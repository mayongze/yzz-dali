package adapter

import (
	"io"
	"log"
	"net"
)

// mockgen -destination=mocks_test.go -package=mocksgithub.com/your-username/your-repo MyService
type BusAdapter interface {
	io.ReadWriteCloser
}

type TcpClientAdapter struct {
	addr string
	conn net.Conn
	Data chan []byte
}

func NewTcpClientAdapter(host string, port string) *TcpClientAdapter {
	t := &TcpClientAdapter{
		addr: net.JoinHostPort(host, port),
		Data: make(chan []byte, 100),
	}
	conn, err := net.Dial("tcp", t.addr)
	if err != nil {
		log.Fatalf("无法连接到服务器: %v", err)
	}
	t.conn = conn
	return t
}

//func (t *TcpClientAdapter) ReadLoop() error {
//	buf := make([]byte, 128)
//	for {
//		n, err := t.conn.Read(buf)
//		if err != nil {
//			log.Printf("读取数据时出错: %v", err)
//			return err
//		}
//		data := buf[:n]
//		log.Printf("收到数据包: %x\n", data)
//		if t.Data != nil {
//			t.Data <- data
//		}
//	}
//}

func (t *TcpClientAdapter) Close() error {
	return t.conn.Close()
}

func (t *TcpClientAdapter) Write(data []byte) (int, error) {
	n, err := t.conn.Write(data)
	if err != nil {
		log.Fatalf("无法写入数据: %v", err)
	}
	return n, nil
}

func (t *TcpClientAdapter) Read(data []byte) (int, error) {
	n, err := t.conn.Read(data)
	if err != nil {
		log.Fatalf("无法读取数据: %v", err)
	}
	return n, nil
}
