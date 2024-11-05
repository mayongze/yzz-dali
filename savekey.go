package homeassistant

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type SystemInfo struct {
	MaxLightNum       uint8
	CurrLightNum      uint8
	CommunicationType uint8
	Transparent       uint8
	GroupAndScene     uint8
}

type DaliCommandType byte

const (
	DaliCommandTypeTx   DaliCommandType = 0x11
	DaliCommandTypeRx   DaliCommandType = 0x12
	DaliCommandTypeInfo DaliCommandType = 0x14
)

type DALIBusStatus string

const (
	DALIBusStatusOK          DALIBusStatus = "总线正常"
	DALIBusStatusNoPower     DALIBusStatus = "未接上DALI电源或出现故障"
	DALIBusStatusDataErr     DALIBusStatus = "接收数据开始出错"
	DALIBusStatusDataLong    DALIBusStatus = "接收数据太长出错"
	DALIBusStatusBusLevelErr DALIBusStatus = "接收数据总线电平出错"
	DALIBusStatusUnknown     DALIBusStatus = "未知错误"
)

type DALICMDOpStatus string

const (
	DALICMDOpStatusFailed  DALICMDOpStatus = "操作失败"
	DALICMDOpStatusOK      DALICMDOpStatus = "操作成功"
	DALICMDOpStatusExit    DALICMDOpStatus = "退出设置"
	DALICMDOpStatusSeqDiff DALICMDOpStatus = "设备序号不一致"
	DALICMDOpStatusNotSup  DALICMDOpStatus = "功能不支持"
	DALICMDOpStatusUnauth  DALICMDOpStatus = "此功能未授权"
	DALICMDOpStatusUnknow  DALICMDOpStatus = "未知错误"
)

type SaveKeyDALIPackage struct {
	CommandType    DaliCommandType
	Control        byte
	SequenceNumber byte
	DeviceID       [3]byte
	Status         byte
	GroupControl   byte
	InstructionSet [6]byte
	Data14         byte
}

func (d *SaveKeyDALIPackage) GetBusStatus() DALIBusStatus {
	// 高4位是总线状态 低4位是操作状态
	switch d.Status >> 4 {
	case 0:
		return DALIBusStatusOK
	case 1:
		return DALIBusStatusNoPower
	case 2:
		return DALIBusStatusDataErr
	case 3:
		return DALIBusStatusDataLong
	case 4:
		return DALIBusStatusBusLevelErr
	}
	return DALIBusStatusUnknown
}

func (d *SaveKeyDALIPackage) GetInstructionSet() []uint16 {
	var insSet []uint16
	for idx := 0; idx < int(d.GroupControl&0x0F); idx++ {
		insSet = append(insSet, uint16(d.InstructionSet[idx*2])<<8|uint16(d.InstructionSet[idx*2+1]))
	}
	return insSet
}

func (d *SaveKeyDALIPackage) GetOpStatus() DALICMDOpStatus {
	switch d.Status & 0x0F {
	case 0:
		return DALICMDOpStatusFailed
	case 1:
		return DALICMDOpStatusOK
	case 2:
		return DALICMDOpStatusExit
	case 4:
		return DALICMDOpStatusSeqDiff
	case 5:
		return DALICMDOpStatusNotSup
	case 6:
		return DALICMDOpStatusUnauth
	}
	return DALICMDOpStatusUnknow
}

func NewSaveKeyDALIForwardPkg(deviceID [3]byte, instructionSet []uint16) (*SaveKeyDALIPackage, error) {
	payload := &SaveKeyDALIPackage{
		CommandType:    DaliCommandTypeTx,
		Control:        0x00,
		SequenceNumber: 0x01,
		DeviceID:       deviceID,
		Status:         0x00,
		GroupControl:   0x00,
	}
	if len(instructionSet) > 3 || len(instructionSet) == 0 {
		return nil, fmt.Errorf("invalid instruction set length")
	}
	payload.GroupControl = uint8(2<<4) | uint8(len(instructionSet))
	for idx, ins := range instructionSet {
		payload.InstructionSet[idx*2] = byte(ins >> 8)
		payload.InstructionSet[idx*2+1] = byte(ins)
	}
	return payload, nil
}

func NewDALISystemInfoForwardPkg() *SaveKeyDALIPackage {
	payload := &SaveKeyDALIPackage{
		CommandType:    DaliCommandTypeInfo,
		Control:        0x00,
		SequenceNumber: 0x01,
		DeviceID:       [3]byte{0xFF, 0xFF, 0xFF},
		Status:         0x00,
		GroupControl:   0x00,
	}
	return payload
}

func (d *SaveKeyDALIPackage) AddInstructionSet(instructionSet uint16) error {
	currentLen := int(d.GroupControl & 0x0F)
	if currentLen > 1 {
		return fmt.Errorf("invalid instruction set length")
	}
	d.GroupControl = uint8(2<<4) | uint8(currentLen+1)
	d.InstructionSet[currentLen*2] = byte(instructionSet >> 8)
	d.InstructionSet[currentLen*2+1] = byte(instructionSet)
	return nil
}

func (d *SaveKeyDALIPackage) SetInstructionSet(instructionSet []uint16) error {
	if len(instructionSet) > 2 || len(instructionSet) == 0 {
		return fmt.Errorf("invalid instruction set length")
	}
	d.GroupControl = uint8(2<<4) | uint8(len(instructionSet))
	for idx, ins := range instructionSet {
		d.InstructionSet[idx*2] = byte(ins >> 8)
		d.InstructionSet[idx*2+1] = byte(ins)
	}
	return nil
}

func (d *SaveKeyDALIPackage) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, binary.BigEndian, *d)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d *SaveKeyDALIPackage) BuildSaveKeyPackage(deviceAddr int) (*SaveKeyPackage, error) {
	buf, err := d.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return NewSaveKeyForwardPackage(ServiceTypeDali).SetPayload(buf).SetDeviceAddr(uint8(deviceAddr)), nil
}

type CommDirection uint8

const (
	DirectionTx CommDirection = 0x8E
	DirectionRx CommDirection = 0x0E
)

type ServiceType uint8

const (
	ServiceTypeDali ServiceType = 0x19
)

type SaveKeyPackage struct {
	Identifier  byte
	Direction   CommDirection
	DeviceAddr  uint8
	ServiceType ServiceType

	DataLength uint8
	Payload    []byte
	Checksum   byte
}

func (s *SaveKeyPackage) calculateChecksum() byte {
	var checksum uint8
	checksum = s.Identifier + uint8(s.Direction) + s.DeviceAddr + uint8(s.ServiceType) + s.DataLength
	for _, b := range s.Payload {
		checksum += b
	}
	s.Checksum = byte(checksum)
	return byte(checksum)
}

func (s *SaveKeyPackage) SetPayload(payload []byte) *SaveKeyPackage {
	s.Payload = payload
	s.DataLength = uint8(len(payload))
	return s
}

func (s *SaveKeyPackage) SetDeviceAddr(deviceAddr uint8) *SaveKeyPackage {
	s.DeviceAddr = deviceAddr
	return s
}

func (s *SaveKeyPackage) MarshalBinary() ([]byte, error) {
	if len(s.Payload) == 0 {
		return nil, errors.New("invalid payload")
	}
	s.calculateChecksum()
	buf := bytes.NewBufferString("")
	_ = binary.Write(buf, binary.BigEndian, s.Identifier)
	_ = binary.Write(buf, binary.BigEndian, s.Direction)
	_ = binary.Write(buf, binary.BigEndian, s.DeviceAddr)
	_ = binary.Write(buf, binary.BigEndian, s.ServiceType)
	_ = binary.Write(buf, binary.BigEndian, s.DataLength)
	_ = binary.Write(buf, binary.BigEndian, s.Payload)
	_ = binary.Write(buf, binary.BigEndian, s.Checksum)
	return buf.Bytes(), nil
}

func (s *SaveKeyPackage) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	_ = binary.Read(buf, binary.BigEndian, &s.Identifier)
	_ = binary.Read(buf, binary.BigEndian, &s.Direction)
	_ = binary.Read(buf, binary.BigEndian, &s.DeviceAddr)
	_ = binary.Read(buf, binary.BigEndian, &s.ServiceType)
	_ = binary.Read(buf, binary.BigEndian, &s.DataLength)
	s.Payload = make([]byte, s.DataLength)
	_ = binary.Read(buf, binary.BigEndian, &s.Payload)
	_ = binary.Read(buf, binary.BigEndian, &s.Checksum)
	return nil
}

func (s *SaveKeyPackage) HexDump() string {
	if len(s.Payload) > 0 {
		buf, _ := s.MarshalBinary()
		return fmt.Sprintf("%X", buf)
	}
	return ""
}

func NewSaveKeyForwardPackage(serviceType ServiceType) *SaveKeyPackage {
	return &SaveKeyPackage{
		Identifier: 0x11,
		Direction:  DirectionTx,
		// DeviceAddr:  deviceAddr,
		DeviceAddr:  0,
		ServiceType: serviceType,
	}
}

type SaveKeyDevice struct {
	sBus         *SBus
	SerialNumber [3]byte
	DeviceType   string
	BusAddress   int
	Name         string
	UID          string

	DALIOutputPkg chan *SaveKeyDALIPackage
}

func NewSaveKeyDevice(sBus *SBus, serialNumber [3]byte, busAddress int) *SaveKeyDevice {
	device := &SaveKeyDevice{
		sBus:          sBus,
		SerialNumber:  serialNumber,
		BusAddress:    busAddress,
		DALIOutputPkg: make(chan *SaveKeyDALIPackage, 100),
	}
	sBus.SaveKeyDeviceMap[busAddress] = device
	return device
}

func (s *SaveKeyDevice) SendDALIForward(daliPackage *SaveKeyDALIPackage) error {
	saveKeyPkg, err := daliPackage.BuildSaveKeyPackage(s.BusAddress)
	if err != nil {
		return err
	}
	return s.sBus.RawSend(saveKeyPkg)
}

func (s *SaveKeyDevice) ReadDALIPackage(timeout time.Duration) *SaveKeyDALIPackage {
	select {
	case <-time.After(timeout):
		return nil
	case pkg := <-s.DALIOutputPkg:
		return pkg
	}
}

func (s *SaveKeyDevice) Recv(data []byte) {
	saveKeyPkg := new(SaveKeyPackage)
	err := saveKeyPkg.UnmarshalBinary(data)
	if err != nil {
		log.Printf("解析数据包出错: %v", err)
		return
	}
	if saveKeyPkg.ServiceType == ServiceTypeDali {
		daliPkg := new(SaveKeyDALIPackage)
		err = binary.Read(bytes.NewBuffer(saveKeyPkg.Payload), binary.BigEndian, daliPkg)
		if err != nil {
			log.Printf("解析DALI数据包出错: %v", err)
			return
		}
		if daliPkg.SequenceNumber == 1 {
			s.DALIOutputPkg <- daliPkg
		} else {
			log.Printf("检测到其他操作序列号: %d", daliPkg.SequenceNumber)
			log.Printf("[debug] daliPkg: %+vnstruction, busStatus:%s opStatus:%s instructionSet:%x", daliPkg,
				daliPkg.GetBusStatus(), daliPkg.GetOpStatus(), daliPkg.GetInstructionSet())
		}
	} else {
		log.Printf("未知服务类型: %d", saveKeyPkg.ServiceType)
	}
}

type SaveKeyDALIDriver struct {
	lock          sync.Mutex
	SaveKeyDevice *SaveKeyDevice
}

func NewSaveKeyDALIDriver(device *SaveKeyDevice) *SaveKeyDALIDriver {
	return &SaveKeyDALIDriver{
		SaveKeyDevice: device,
	}
}

func (sd *SaveKeyDALIDriver) rawSend(instructionSet []uint16) error {
	p, err := NewSaveKeyDALIForwardPkg(sd.SaveKeyDevice.SerialNumber, instructionSet)
	if err != nil {
		return err
	}
	var respPkg *SaveKeyDALIPackage
	for i := 0; i < 2; i++ {
		if err = sd.SaveKeyDevice.SendDALIForward(p); err != nil {
			return err
		}
		// 先等发送成功消息
		respPkg = sd.SaveKeyDevice.ReadDALIPackage(500 * time.Millisecond)
		if respPkg != nil {
			break
		}
	}
	if respPkg == nil {
		return fmt.Errorf("resp timeout")
	}
	log.Printf("[debug] daliPkg: %+vnstruction, busStatus:%s opStatus:%s instructionSet:%x", respPkg,
		respPkg.GetBusStatus(), respPkg.GetOpStatus(), respPkg.GetInstructionSet())

	if respPkg.GetBusStatus() != DALIBusStatusOK {
		return fmt.Errorf("%s", respPkg.GetBusStatus())
	}
	if respPkg.GetOpStatus() != DALICMDOpStatusOK {
		return fmt.Errorf("%s", respPkg.GetOpStatus())
	}
	return nil
}

func (sd *SaveKeyDALIDriver) SendCommand(commands []DaliCommand) error {
	sd.lock.Lock()
	defer sd.lock.Unlock()
	// command中要添加指令序列号来对应关系
	// 同一批命令使用相同的指令序列号
	commandsFix := make([]DaliCommand, 0)
	for _, command := range commands {
		if command.DeviceType() != 0 {
			commandsFix = append(commandsFix, EnableDeviceType(command.DeviceType()))
		}
		commandsFix = append(commandsFix, command)
		if command.SendTwice() {
			commandsFix = append(commandsFix, command)
		}
	}
	instructionSet := make([]uint16, 0)
	for _, command := range commandsFix {
		instructionSet = append(instructionSet, command.Instruction())
	}

	// instructionSet拆分成2个一组进行发送
	for i := 0; i < len(instructionSet); i += 3 {
		if i+2 < len(instructionSet) {
			t0 := time.Now()
			if err := sd.rawSend([]uint16{instructionSet[i], instructionSet[i+1], instructionSet[i+2]}); err != nil {
				return err
			}
			for j := i; j < i+3; j++ {
				command := commandsFix[j]
				if !command.RequiresReply() {
					continue
				}
				respPkg := sd.SaveKeyDevice.ReadDALIPackage(300 * time.Millisecond)
				if respPkg == nil {
					return fmt.Errorf("resp timeout")
				}
				err := command.ReadReply(respPkg.InstructionSet[:])
				if err != nil {
					return err
				}
			}
			log.Printf("send 3条指令 %d-%d cost: %dms", i, i+2, time.Since(t0).Milliseconds())
		} else if i+1 < len(instructionSet) {
			t0 := time.Now()
			if err := sd.rawSend([]uint16{instructionSet[i], instructionSet[i+1]}); err != nil {
				return err
			}
			for j := i; j < i+2; j++ {
				command := commandsFix[j]
				if !command.RequiresReply() {
					continue
				}
				respPkg := sd.SaveKeyDevice.ReadDALIPackage(300 * time.Millisecond)
				if respPkg == nil {
					return fmt.Errorf("resp timeout")
				}
				err := command.ReadReply(respPkg.InstructionSet[:])
				if err != nil {
					return err
				}
			}
			log.Printf("send 2条指令 %d-%d cost: %dms", i, i+1, time.Since(t0).Milliseconds())
		} else {
			t0 := time.Now()
			if err := sd.rawSend([]uint16{instructionSet[i]}); err != nil {
				return err
			}
			for j := i; j < i+1; j++ {
				command := commandsFix[j]
				if !command.RequiresReply() {
					continue
				}
				respPkg := sd.SaveKeyDevice.ReadDALIPackage(300 * time.Millisecond)
				if respPkg == nil {
					return fmt.Errorf("resp timeout")
				}
				err := command.ReadReply(respPkg.InstructionSet[:])
				if err != nil {
					return err
				}
			}
			log.Printf("send 1条指令 %d-%d cost: %dms", i, i, time.Since(t0).Milliseconds())
		}
	}
	//for _, command := range commandsFix {
	//	if !command.RequiresReply() {
	//		continue
	//	}
	//	respPkg := sd.SaveKeyDevice.ReadDALIPackage(300 * time.Millisecond)
	//	if respPkg == nil {
	//		return fmt.Errorf("resp timeout")
	//	}
	//	err := command.ReadReply(respPkg.InstructionSet[:])
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}

func (sd *SaveKeyDALIDriver) SendInstructionSet(instructionSet []uint16) error {
	p, err := NewSaveKeyDALIForwardPkg(sd.SaveKeyDevice.SerialNumber, instructionSet)
	if err != nil {
		return err
	}
	if err = sd.SaveKeyDevice.SendDALIForward(p); err != nil {
		return err
	}
	// 先等发送成功消息
	respPkg := sd.SaveKeyDevice.ReadDALIPackage(time.Second)
	if respPkg == nil {
		return fmt.Errorf("resp timeout")
	}
	log.Printf("[debug] daliPkg: %+v nstruction, busStatus:%s opStatus:%s instructionSet:%x", respPkg,
		respPkg.GetBusStatus(), respPkg.GetOpStatus(), respPkg.GetInstructionSet())

	if respPkg.GetBusStatus() != DALIBusStatusOK {
		return fmt.Errorf("%s", respPkg.GetBusStatus())
	}
	if respPkg.GetOpStatus() != DALICMDOpStatusOK {
		return fmt.Errorf("%s", respPkg.GetOpStatus())
	}

	return nil
}

func (sd *SaveKeyDALIDriver) GetSystemInfo() error {
	p := NewDALISystemInfoForwardPkg()
	if err := sd.SaveKeyDevice.SendDALIForward(p); err != nil {
		return err
	}
	return nil
}
