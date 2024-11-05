package homeassistant

import (
	"encoding"
	"fmt"
	"io"
	"log"
	"sync"
)

type SBus struct {
	mux              sync.Mutex
	BusConn          io.ReadWriteCloser
	SaveKeyDeviceMap map[int]*SaveKeyDevice
	StrDeviceMap     map[int]*StrDanceDevice
}

func NewSBus(bus io.ReadWriteCloser) *SBus {
	sBus := &SBus{
		BusConn:          bus,
		SaveKeyDeviceMap: make(map[int]*SaveKeyDevice),
		StrDeviceMap:     make(map[int]*StrDanceDevice),
	}
	go func() {
		_ = sBus.ReadLoop()
	}()
	return sBus
}

func (s *SBus) RawSend(pkg encoding.BinaryMarshaler) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	bs, err := pkg.MarshalBinary()
	if err != nil {
		return err
	}
	log.Printf("发送数据包: %x\n", bs)
	n, err := s.BusConn.Write(bs)
	if err != nil {
		return err
	}
	if n != len(bs) {
		return fmt.Errorf("write bs not match")
	}
	return nil
}

func (s *SBus) ReadLoop() error {
	buf := make([]byte, 128)
	for {
		n, err := s.BusConn.Read(buf)
		if err != nil {
			log.Printf("读取数据时出错: %v", err)
			return err
		}
		data := buf[:n]
		if len(data) > 4 && data[0] == 0x11 && data[1] == 0x0E {
			pkgSize := data[4]
			if len(data) < int(pkgSize)+1+5 {
				n1, err := s.BusConn.Read(buf[n:])
				if err != nil {
					log.Printf("读取数据时出错: %v", err)
					return err
				}
				data = append(data, buf[n:n+n1]...)
			}
			log.Printf("收到数据包: %x\n", data)
			deviceAddr := data[2]
			device, ok := s.SaveKeyDeviceMap[int(deviceAddr)]
			if !ok {
				log.Printf("device not found: %d", deviceAddr)
				continue
			}
			device.Recv(data)
		} else if len(data) == 8 && (data[0] == 0x55 || data[0] == 0x22) {
			// 固定长度
			deviceAddr := data[1]
			device, ok := s.StrDeviceMap[int(deviceAddr)]
			if !ok {
				log.Printf("device not found: %d", deviceAddr)
				continue
			}
			device.Recv(data)
		}
	}
}
