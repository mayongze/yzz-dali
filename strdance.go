package homeassistant

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"huiqun/homeassistant/hadevice"
	"log"
)

type StrDanceFunction byte

const (
	ReadStatus StrDanceFunction = 0x10 // 读取状态
)

type StrDancePackage struct {
	// 55 send 22 recv
	Identifier byte
	DeviceAddr byte
	Function   byte
	Data       [4]byte
	Checksum   byte
}

func (s *StrDancePackage) calculateChecksum() *StrDancePackage {
	var checksum uint8
	checksum = s.Identifier + s.DeviceAddr + s.Function
	for _, b := range s.Data {
		checksum += b
	}
	s.Checksum = checksum
	return s
}

func (s *StrDancePackage) MarshalBinary() ([]byte, error) {
	s.calculateChecksum()
	buf := bytes.NewBufferString("")
	_ = binary.Write(buf, binary.BigEndian, s)
	return buf.Bytes(), nil
}

func (s *StrDancePackage) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	_ = binary.Read(buf, binary.BigEndian, s)
	return nil
}

func (s *StrDancePackage) HexDump() string {
	buf, _ := s.MarshalBinary()
	return fmt.Sprintf("%X", buf)
}

func NewStrDanceForwardPackage(addr byte, function byte, data [4]byte) *StrDancePackage {
	pkg := &StrDancePackage{
		Identifier: 0x55,
		DeviceAddr: addr,
		Function:   function,
		Data:       data,
	}
	pkg.calculateChecksum()
	return pkg
}

type StrDanceDevice struct {
	sBus       *SBus
	BusAddress int
	Name       string
	UID        string
	State      *SwitchState
	// OutputPkg  chan *StrDancePackage
}

func NewStrDanceDevice(sBus *SBus, busAddress int, uid, name string) *StrDanceDevice {
	strDanceDevice := &StrDanceDevice{
		sBus:       sBus,
		BusAddress: busAddress,
		Name:       name,
		UID:        uid,
		State:      GlobalState.GetSwitchState(uid),
		// OutputPkg:  make(chan *StrDancePackage, 100),
	}
	sBus.StrDeviceMap[busAddress] = strDanceDevice
	return strDanceDevice
}

func (s *StrDanceDevice) GetDeviceType() DeviceType {
	return StrDanceDeviceType
}

func (s *StrDanceDevice) Recv(data []byte) {
	pkg := new(StrDancePackage)
	err := pkg.UnmarshalBinary(data)
	if err != nil {
		log.Printf("解析数据包出错: %v", err)
		return
	}
	log.Printf("收到数据包: %+v\n", pkg)

	if pkg.Function == 0x10 {
		//pkg.Data 4字节 32位 1代表ON 0代表OFF
		var switchState int32
		err = binary.Read(bytes.NewReader(pkg.Data[:]), binary.BigEndian, &switchState)
		if err != nil {
			log.Printf("解析数据包Data出错: %v", err)
			return
		}
		if switchState&1 == 1 {
			s.State.L1 = "ON"
		} else {
			s.State.L1 = "OFF"
		}
		if switchState>>1&1 == 1 {
			s.State.L2 = "ON"
		} else {
			s.State.L2 = "OFF"
		}
		if switchState>>2&1 == 1 {
			s.State.L3 = "ON"
		} else {
			s.State.L3 = "OFF"
		}
		if switchState>>3&1 == 1 {
			s.State.L4 = "ON"
		} else {
			s.State.L4 = "OFF"
		}
		log.Printf("状态数据: %+v\n", s.State)
		token := GetMqttCli().Publish(fmt.Sprintf("yzz-dali/%s", s.UID), 0, false, s.State.GetBytes())
		if token.Error() != nil {
			log.Printf("发送状态出错: %v", token.Error())
		}
	}

	if pkg.Function == 0x20 || pkg.Function == 0x21 || pkg.Function == 0x22 {
		var switchState int32
		err = binary.Read(bytes.NewReader(pkg.Data[:]), binary.BigEndian, &switchState)
		if err != nil {
			log.Printf("解析数据包Data出错: %v", err)
			return
		}
		str := ""
		for i := 0; i <= 3; i++ {
			if switchState>>i&1 == 1 {
				str = fmt.Sprintf("l%d", i+1)
				break
			}
		}
		switch pkg.Function {
		case 0x20:
			str = fmt.Sprintf("%s_%s", str, "single")
		case 0x21:
			str = fmt.Sprintf("%s_%s", str, "double")
		case 0x22:
			str = fmt.Sprintf("%s_%s", str, "hold")
		}
		s.State.Action = str
		token := GetMqttCli().Publish(fmt.Sprintf("yzz-dali/%s", s.UID), 0, false, s.State.GetBytes())
		if token.Error() != nil {
			log.Printf("发送状态出错: %v", token.Error())
		}
	}

	if pkg.Function == 0x23 || pkg.Function == 0x24 {
		var switchState int32
		err = binary.Read(bytes.NewReader(pkg.Data[:]), binary.BigEndian, &switchState)
		if err != nil {
			log.Printf("解析数据包Data出错: %v", err)
			return
		}
		action := "press"
		if pkg.Function == 0x24 {
			action = "release"
		}
		if switchState&1 == 1 {
			s.State.StateL1 = action
		}
		if switchState>>1&1 == 1 {
			s.State.StateL2 = action
		}
		if switchState>>2&1 == 1 {
			s.State.StateL3 = action
		}
		if switchState>>3&1 == 1 {
			s.State.StateL4 = action
		}
		token := GetMqttCli().Publish(fmt.Sprintf("yzz-dali/%s", s.UID), 0, false, s.State.GetBytes())
		if token.Error() != nil {
			log.Printf("发送状态出错: %v", token.Error())
		}
	}

}

func (s *StrDanceDevice) getHaDevice() hadevice.Device {
	return hadevice.Device{
		Identifiers:  []string{fmt.Sprintf("yzz-dali_strdance_%s", s.UID)},
		Manufacturer: `StrDance`,
		Model:        "switch-4-gang",
		Name:         s.Name,
		UID:          s.UID,
	}
}

func (s *StrDanceDevice) InitDiscover() error {
	lCount := 4
	for i := 1; i <= lCount; i++ {
		topic := fmt.Sprintf("homeassistant/switch/%s/l%d/config", s.UID, i)
		switchCfg := hadevice.NewSwitch(s.getHaDevice(), fmt.Sprintf("开关%d", i), fmt.Sprintf("l%d", i), "ON", "OFF")
		token := GetMqttCli().Publish(topic, 0, true, switchCfg.GetBytes())
		if token.Error() != nil {
			return token.Error()
		}
		topic = fmt.Sprintf("homeassistant/sensor/%s/state_l%d/config", s.UID, i)
		sensorCfg := hadevice.NewSensor(s.getHaDevice(), fmt.Sprintf("按键%d状态", i), fmt.Sprintf("state_l%d", i))
		token = GetMqttCli().Publish(topic, 0, true, sensorCfg.GetBytes())
		if token.Error() != nil {
			return token.Error()
		}
	}

	topic := fmt.Sprintf("homeassistant/sensor/%s/action/config", s.UID)
	sensorCfg := hadevice.NewSensor(s.getHaDevice(), "Action", "action")
	token := GetMqttCli().Publish(topic, 0, true, sensorCfg.GetBytes())
	if token.Error() != nil {
		return token.Error()
	}
	// 发送查询命令
	return s.sBus.RawSend(NewStrDanceForwardPackage(byte(s.BusAddress), byte(ReadStatus), [4]byte{}))
}

func (s *StrDanceDevice) HandleCommand(stateKV StateKV) error {
	data := [4]byte{}
	var state int64
	if s.State.L1 == "ON" {
		state |= 1
	}
	if s.State.L2 == "ON" {
		state |= 1 << 1
	}
	if s.State.L3 == "ON" {
		state |= 1 << 2
	}
	if s.State.L4 == "ON" {
		state |= 1 << 3
	}
	l1, ok := GetStateV[string](stateKV, "l1")
	if ok {
		if l1 == "ON" {
			state |= 1
		} else {
			state &^= 1
		}
	}
	l2, ok := GetStateV[string](stateKV, "l2")
	if ok {
		if l2 == "ON" {
			state |= 1 << 1
		} else {
			state &^= 1 << 1
		}
	}
	l3, ok := GetStateV[string](stateKV, "l3")
	if ok {
		if l3 == "ON" {
			state |= 1 << 2
		} else {
			state &^= 1 << 2
		}
	}
	l4, ok := GetStateV[string](stateKV, "l4")
	if ok {
		if l4 == "ON" {
			state |= 1 << 3
		} else {
			state &^= 1 << 3
		}
	}
	binary.BigEndian.PutUint32(data[:], uint32(state))
	return s.sBus.RawSend(NewStrDanceForwardPackage(byte(s.BusAddress), 0x33, data))
}
