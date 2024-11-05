package adapter

import (
	"go.bug.st/serial"
	"log"
	"time"
)

type SerialAdapter struct {
	portName string
	baud     int
	port     serial.Port
}

func NewSerialAdapter(portName string, baud int) *SerialAdapter {
	mode := &serial.Mode{
		BaudRate: baud,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatalf("无法连接到串口: %v", err)
	}
	t := &SerialAdapter{
		portName: portName,
		baud:     baud,
		port:     port,
	}
	return t
}

func (t *SerialAdapter) Close() error {
	return t.port.Close()
}

func (t *SerialAdapter) Write(data []byte) (int, error) {
	n, err := t.port.Write(data)
	if err != nil {
		log.Fatalf("无法写入数据: %v", err)
	}
	return n, nil
}

func (t *SerialAdapter) Read(data []byte) (int, error) {
	t.port.SetReadTimeout(20 * time.Millisecond)
	var first time.Time
	l := 0
	for {
		buf := make([]byte, 128)
		n, err := t.port.Read(buf)
		if n == 0 && err == nil {
			if !first.IsZero() && time.Since(first) > 20*time.Millisecond {
				return l, nil
			}
			continue
		}
		if err != nil {
			log.Fatalf("无法读取数据: %v", err)
		}
		if n > 0 {
			if first.IsZero() {
				first = time.Now()
			}
			copy(data[l:], buf[:n])
			l += n
		}
	}
}
