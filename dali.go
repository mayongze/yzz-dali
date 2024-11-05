package homeassistant

import (
	"errors"
	"fmt"
	"github.com/antlabs/gstl/set"
	"math"
	"strings"
)

type Address uint16

func BroadcastAddr() Address {
	return Address(0xFE00)
}

func GroupAddr(addr int) Address {
	if addr < 0 || addr > 15 {
		panic("invalid group address")
	}
	return 0x8000 | Address(addr<<9)
}

func ShortAddr(addr int) Address {
	if addr < 0 || addr > 63 {
		panic("invalid short address")
	}
	return Address(addr << 9)
}

func DAPC(deviceAddr Address, level int) *DAPCCommand {
	return NewDAPCCommand(deviceAddr, level)
}

// Off 0xff
func Off(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x00)
	cmd.isReply = false
	return &cmd
}

func Up(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x01)
	cmd.isReply = false
	return &cmd
}

func Down(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x02)
	cmd.isReply = false
	return &cmd
}

func StepUp(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x03)
	cmd.isReply = false
	return &cmd
}

func StepDown(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x04)
	cmd.isReply = false
	return &cmd
}

func RecallMaxLevel(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x05)
	cmd.isReply = false
	return &cmd
}

func RecallMinLevel(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x06)
	cmd.isReply = false
	return &cmd
}

func StepDownAndOff(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x07)
	cmd.isReply = false
	return &cmd
}

func OnAndStepUp(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x08)
	cmd.isReply = false
	return &cmd
}

func EnableDAPCSequence(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x09)
	cmd.isReply = false
	return &cmd
}

func GoToLastActiveLevel(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x0A)
	cmd.isReply = false
	return &cmd
}

func ContinuousUp(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x0B)
	cmd.isReply = false
	return &cmd
}

func ContinuousDown(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x0C)
	cmd.isReply = false
	return &cmd
}

func GoToScene(deviceAddr Address, scene int) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x10, byte(scene))
	cmd.isReply = false
	return &cmd
}

func Reset(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x20)
	cmd.isReply = false
	return &cmd
}

func QueryMaxLevel(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0xA1)
	cmd.isReply = true
	return &cmd
}

func QueryMinLevel(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0xA2)
	cmd.isReply = true
	return &cmd
}

func QueryPowerOnLevel(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0xA3)
	cmd.isReply = true
	return &cmd
}

func QuerySystemFailureLevel(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0xA4)
	cmd.isReply = true
	return &cmd
}

func QueryFastFadeTime(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0xFD)
	cmd.isReply = true
	cmd.deviceType = 6
	return &cmd
}

func QueryMinFastFadeTime(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0xFE)
	cmd.isReply = true
	cmd.deviceType = 6
	return &cmd
}

func QueryActualLevel(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0xA0)
	cmd.isReply = true
	return &cmd
}

func SetDTR0(value byte) *SpecialCommand {
	cmd := NewSpecialCommand(0xA3, value)
	return cmd
}

func SetDTR1(value byte) *SpecialCommand {
	cmd := NewSpecialCommand(0xC3, value)
	return cmd
}

func SetDTR2(value byte) *SpecialCommand {
	cmd := NewSpecialCommand(0xC5, value)
	return cmd
}

func QueryColourValue(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 250)
	cmd.deviceType = 8
	return &cmd
	// return uint16(deviceAddr) | 0x100 | 250
}

func SetTemporaryColourTemperature(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 231)
	cmd.deviceType = 8
	cmd.isReply = false
	return &cmd
}

func SetPowerOnLevel(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x2d)
	//cmd.deviceType = 8
	cmd.isReply = false
	return &cmd
}

func Activate(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 226)
	cmd.deviceType = 8
	cmd.isReply = false
	return &cmd
}

func QueryContentDTR0(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 152)
	return &cmd
	// return uint16(deviceAddr) | 0x100 | 152
}

type QueryDeviceTypeCmd StandardCommand

func NewQueryDeviceTypeCmd(deviceAddr Address) *QueryDeviceTypeCmd {
	cmd := QueryDeviceTypeCmd(NewStandardCommand(deviceAddr, 0x99))
	cmd.isReply = true
	return &cmd
}

type QueryNextDeviceType StandardCommand

func NewQueryNextDeviceTypeCmd(deviceAddr Address) *QueryNextDeviceType {
	cmd := QueryNextDeviceType(NewStandardCommand(deviceAddr, 0xA7))
	cmd.isReply = true
	return &cmd
}

type QueryFadeTimeFadeRateCmd StandardCommand

type FadeTimeFadeRate struct {
	FadeTime float64
	FadeRate int
}

func (cmd *QueryFadeTimeFadeRateCmd) GetFadeTimeFadeRate() FadeTimeFadeRate {
	ft := cmd.val >> 4
	rst := FadeTimeFadeRate{}
	if ft > 0 {
		rst.FadeTime = 0.5 * math.Sqrt(math.Pow(2, float64(ft)))
	}
	// rst.FadeRate = float64(506) / math.Pow(float64(2), float64(cmd.val&0x0F)/float64(2))
	rst.FadeRate = int(cmd.val) & 0x0F
	return rst
}

func NewQueryFadeTimeFadeRateCmd(deviceAddr Address) *QueryFadeTimeFadeRateCmd {
	cmd := QueryFadeTimeFadeRateCmd(NewStandardCommand(deviceAddr, 0xA5))
	cmd.isReply = true
	return &cmd
}

type QueryExtendedFadeTimeCmd StandardCommand

type ExtendedFadeTimeMultiplier byte

const (
	ExtendedFadeTimeMultiplierZero           = 0
	ExtendedFadeTimeMultiplier100Millisecond = 1
	ExtendedFadeTimeMultiplier1Second        = 2
	ExtendedFadeTimeMultiplier10Seconds      = 3
	ExtendedFadeTimeMultiplier1Minute        = 4
)

// https://infosys.beckhoff.com/english.php?content=../content/1033/tcplclib_tc3_dali/5627748491.html&id=

// 0 100 200 300 400 500 600 700 800 900 1000 1100 1200 1300 1400 1500 1600 2000 3000 4000 5000 6000 7000 8000 9000 10000
// 11000 12000 13000 14000 15000 16000 20000 30000 40000 50000 60000 70000 80000 90000 100000 110000 120000 130000 140000 150000 160000
// 180000 240000

// QueryMinFastFadeTime --> 25ms一个跳变 --> 25*27

type ExtendedFadeTime struct {
	ExtendedFadeTimeMultiplier ExtendedFadeTimeMultiplier
	ExtendedFadeTimeBase       int
}

func (cmd *ExtendedFadeTime) Millisecond() int {
	switch cmd.ExtendedFadeTimeMultiplier {
	case ExtendedFadeTimeMultiplierZero:
		return 0
	case ExtendedFadeTimeMultiplier100Millisecond:
		return int(cmd.ExtendedFadeTimeBase) * 100
	case ExtendedFadeTimeMultiplier1Second:
		return int(cmd.ExtendedFadeTimeBase) * 1000
	case ExtendedFadeTimeMultiplier10Seconds:
		return int(cmd.ExtendedFadeTimeBase) * 10000
	case ExtendedFadeTimeMultiplier1Minute:
		return int(cmd.ExtendedFadeTimeBase) * 60000
	}
	return 0
}

func (cmd *QueryExtendedFadeTimeCmd) GetExtendedFadeTime() ExtendedFadeTime {
	return ExtendedFadeTime{
		ExtendedFadeTimeMultiplier: ExtendedFadeTimeMultiplier(cmd.val >> 4),
		ExtendedFadeTimeBase:       int(cmd.val&0x0F) + 1,
	}
}

func NewQueryExtendedFadeTimeCmd(deviceAddr Address) *QueryExtendedFadeTimeCmd {
	cmd := QueryExtendedFadeTimeCmd(NewStandardCommand(deviceAddr, 0xA8))
	cmd.isReply = true
	return &cmd
}

func SetFastFadeTime(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0xe4)
	cmd.isReply = false
	cmd.deviceType = 6
	cmd.sendTwice = true
	return &cmd
}

func SetFadeTimeCmd(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x2e)
	cmd.sendTwice = true
	cmd.isReply = false
	return &cmd
}

func SetFadeRateCmd(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x2f)
	cmd.sendTwice = true
	cmd.isReply = false
	return &cmd
}

func SetExtendedFadeTimeCmd(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x30)
	cmd.isReply = false
	return &cmd
}

func QueryGroupsZeroToSeven(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0xc0)
	cmd.isReply = true
	return &cmd
}

func QueryGroupsEightToFifteen(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0xc1)
	cmd.isReply = true
	return &cmd
}

func SetScene(deviceAddr Address, scene int) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x40, byte(scene))
	cmd.isReply = false
	return &cmd
}

func RemoveFromScene(deviceAddr Address, scene int) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x50, byte(scene))
	cmd.isReply = false
	return &cmd
}

func AddToGroup(deviceAddr Address, group int) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x60, byte(group))
	cmd.isReply = false
	return &cmd
}

func RemoveFromGroup(deviceAddr Address, group int) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x70, byte(group))
	cmd.isReply = false
	return &cmd
}

func QuerySceneLevel(deviceAddr Address, scene int) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0xb0, byte(scene))
	cmd.isReply = true
	return &cmd
}

func IdentifyDevice(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0x25)
	cmd.sendTwice = true
	cmd.isReply = false
	return &cmd
}

func ReadMemoryLocation(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0xc5)
	cmd.isReply = true
	return &cmd
}

type QueryColourStatusCmd StandardCommand

type ColourStatus struct {
	XYColourPointOutOfRange             bool `bitmap:"0"`
	ColourTemperatureTcOutOfRange       bool `bitmap:"1"`
	AutoCalibrationRunning              bool `bitmap:"2"`
	AutoCalibrationSuccessful           bool `bitmap:"3"`
	ColourTypeXYActive                  bool `bitmap:"4"`
	ColourTypeColourTemperatureTcActive bool `bitmap:"5"`
	ColourTypePrimaryNActive            bool `bitmap:"6"`
	ColourTypeRGBWAFActive              bool `bitmap:"7"`
}

func (cmd *QueryColourStatusCmd) GetColourStatus() ColourStatus {
	return FromBitmap(cmd.val, ColourStatus{})
}

func NewQueryColourStatusCommand(deviceAddr Address) *QueryColourStatusCmd {
	cmdVal := 248
	cmd := QueryColourStatusCmd(NewStandardCommand(deviceAddr, byte(cmdVal)))
	cmd.isReply = true
	cmd.deviceType = 8
	return &cmd
}

func EnableDeviceType(t byte) *SpecialCommand {
	return NewSpecialCommand(0xC1, t)
}

func QueryDimmingCurve(deviceAddr Address) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0xee)
	cmd.deviceType = 6
	cmd.isReply = true
	return &cmd
}

type DimmingCurve byte

const (
	DimmingCurveStandard DimmingCurve = 0
	DimmingCurveLinear   DimmingCurve = 1
)

func ParseDimmingCurve(str string) DimmingCurve {
	switch str {
	case "standard":
		return DimmingCurveStandard
	case "linear":
		return DimmingCurveLinear
	}
	return 255
}

func (dc DimmingCurve) Byte() byte {
	return byte(dc)
}

func (dc DimmingCurve) String() string {
	switch dc {
	case DimmingCurveStandard:
		return "standard"
	case DimmingCurveLinear:
		return "linear"
	}
	return ""
}

func SelectDimmingCurve(deviceAddr Address, curve DimmingCurve) *StandardCommand {
	cmd := NewStandardCommand(deviceAddr, 0xe3, byte(curve))
	cmd.deviceType = 6
	cmd.sendTwice = true
	cmd.isReply = false
	return &cmd
}

type DaliDriverI interface {
	SendCommand(commands []DaliCommand) error
}

type DaliDriver struct {
	DaliDriverI
}

var daliDriver *DaliDriver

func GetDalDriver() *DaliDriver {
	return daliDriver
}

func InitDaliDriver(drv DaliDriverI) {
	daliDriver = NewDaliDriver(drv)
}

func NewDaliDriver(drv DaliDriverI) *DaliDriver {
	return &DaliDriver{drv}
}

var ErrNoResponse = errors.New("no response to initial query")

func (drv *DaliDriver) QueryDeviceType(deviceAddr Address) ([]byte, error) {
	result := make([]byte, 0)
	cmd := NewQueryDeviceTypeCmd(deviceAddr)
	err := drv.SendCommand([]DaliCommand{cmd})
	if err != nil {
		if strings.EqualFold(err.Error(), "resp timeout") {
			return nil, ErrNoResponse
		}
		return nil, err
	}
	var val byte = 0
	if cmd.val == 0xff {
		for val != 0xfe {
			cmd := NewQueryNextDeviceTypeCmd(deviceAddr)
			err := drv.SendCommand([]DaliCommand{cmd})
			if err != nil {
				return nil, err
			}
			val = cmd.val
			if val != 0xfe {
				result = append(result, cmd.val)
			}
		}
	} else {
		result = append(result, cmd.val)
	}
	return result, nil
}

func (drv *DaliDriver) QueryColourStatus(deviceAddr Address) (ColourStatus, error) {
	cmd := NewQueryColourStatusCommand(deviceAddr)
	err := drv.SendCommand([]DaliCommand{cmd})
	if err != nil {
		return ColourStatus{}, err
	}
	return cmd.GetColourStatus(), nil
}

func (drv *DaliDriver) QueryActualLevel(deviceAddr Address) (int, error) {
	cmd := QueryActualLevel(deviceAddr)
	err := drv.SendCommand([]DaliCommand{cmd})
	if err != nil {
		return 0, err
	}
	return int(cmd.val), nil
}

func (drv *DaliDriver) QueryMinLevel(deviceAddr Address) (int, error) {
	cmd := QueryMinLevel(deviceAddr)
	err := drv.SendCommand([]DaliCommand{cmd})
	if err != nil {
		return 0, err
	}
	return int(cmd.val), nil
}

func (drv *DaliDriver) QueryMaxLevel(deviceAddr Address) (int, error) {
	cmd := QueryMaxLevel(deviceAddr)
	err := drv.SendCommand([]DaliCommand{cmd})
	if err != nil {
		return 0, err
	}
	return int(cmd.val), nil
}

func (drv *DaliDriver) QueryPowerOnLevel(deviceAddr Address) (int, error) {
	cmd := QueryPowerOnLevel(deviceAddr)
	err := drv.SendCommand([]DaliCommand{cmd})
	if err != nil {
		return 0, err
	}
	return int(cmd.val), nil
}

func (drv *DaliDriver) QuerySystemFailureLevel(deviceAddr Address) (int, error) {
	cmd := QuerySystemFailureLevel(deviceAddr)
	err := drv.SendCommand([]DaliCommand{cmd})
	if err != nil {
		return 0, err
	}
	return int(cmd.val), nil
}

func (drv *DaliDriver) QueryFastFadeTime(deviceAddr Address) (int, error) {
	cmd := QueryFastFadeTime(deviceAddr)
	err := drv.SendCommand([]DaliCommand{cmd})
	if err != nil {
		return 0, err
	}
	return int(cmd.val), nil
}

func (drv *DaliDriver) QueryMinFastFadeTime(deviceAddr Address) (int, error) {
	cmd := QueryMinFastFadeTime(deviceAddr)
	err := drv.SendCommand([]DaliCommand{cmd})
	if err != nil {
		return 0, err
	}
	return int(cmd.val), nil
}

func (drv *DaliDriver) SetFastFadeTime(deviceAddr Address, time int) error {
	err := drv.SendCommand([]DaliCommand{SetDTR0(byte(time)), SetFastFadeTime(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) QueryDT8ColourValue(deviceAddr Address, query QueryColourValueDTR) (int, error) {
	qalCmd := QueryActualLevel(deviceAddr)
	colourCmd := QueryColourValue(deviceAddr)
	dtr0Cmd := QueryContentDTR0(deviceAddr)
	err := drv.SendCommand([]DaliCommand{qalCmd, SetDTR0(query.Byte()), colourCmd, dtr0Cmd})
	if err != nil {
		return 0, err
	}
	lsb := dtr0Cmd.val
	msb := colourCmd.val
	// lsb 低8bit msb 高8bit 组合成 16bit int
	// 42的误差
	return int(msb)<<8 | int(lsb), nil
}

func (drv *DaliDriver) SetDT8ColourValueTc(deviceAddr Address, tc int) error {
	lsb := byte(tc & 0xff)
	msb := byte(tc >> 8)
	err := drv.SendCommand([]DaliCommand{SetDTR0(lsb), SetDTR1(msb), SetTemporaryColourTemperature(deviceAddr), Activate(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) SetLevel(deviceAddr Address, level int) error {
	err := drv.SendCommand([]DaliCommand{DAPC(deviceAddr, level)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) SetPowerOnLevel(deviceAddr Address, level byte) error {
	err := drv.SendCommand([]DaliCommand{SetDTR0(level), SetPowerOnLevel(deviceAddr), SetPowerOnLevel(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) Off(deviceAddr Address) error {
	err := drv.SendCommand([]DaliCommand{Off(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) Up(deviceAddr Address) error {
	err := drv.SendCommand([]DaliCommand{Up(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) Down(deviceAddr Address) error {
	err := drv.SendCommand([]DaliCommand{Down(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) StepUp(deviceAddr Address) error {
	err := drv.SendCommand([]DaliCommand{StepUp(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) StepDown(deviceAddr Address) error {
	err := drv.SendCommand([]DaliCommand{StepDown(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

// RecallMaxLevel 不经过渐变调到最大亮度
func (drv *DaliDriver) RecallMaxLevel(deviceAddr Address) error {
	err := drv.SendCommand([]DaliCommand{RecallMaxLevel(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

// RecallMinLevel 不经过渐变调到最小亮度
func (drv *DaliDriver) RecallMinLevel(deviceAddr Address) error {
	err := drv.SendCommand([]DaliCommand{RecallMinLevel(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) StepDownAndOff(deviceAddr Address) error {
	err := drv.SendCommand([]DaliCommand{StepDownAndOff(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) OnAndStepUp(deviceAddr Address) error {
	err := drv.SendCommand([]DaliCommand{OnAndStepUp(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) EnableDAPCSequence(deviceAddr Address) error {
	err := drv.SendCommand([]DaliCommand{EnableDAPCSequence(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

// GoToLastActiveLevel off后调到上次亮度
func (drv *DaliDriver) GoToLastActiveLevel(deviceAddr Address) error {
	err := drv.SendCommand([]DaliCommand{GoToLastActiveLevel(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) ContinuousUp(deviceAddr Address) error {
	err := drv.SendCommand([]DaliCommand{ContinuousUp(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) ContinuousDown(deviceAddr Address) error {
	err := drv.SendCommand([]DaliCommand{ContinuousDown(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) GoToScene(deviceAddr Address, scene int) error {
	err := drv.SendCommand([]DaliCommand{GoToScene(deviceAddr, scene)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) Reset(deviceAddr Address) error {
	err := drv.SendCommand([]DaliCommand{Reset(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) QueryFadeTimeFadeRate(deviceAddr Address) (FadeTimeFadeRate, error) {
	cmd := NewQueryFadeTimeFadeRateCmd(deviceAddr)
	err := drv.SendCommand([]DaliCommand{cmd})
	if err != nil {
		return FadeTimeFadeRate{}, err
	}
	return cmd.GetFadeTimeFadeRate(), nil
}

func (drv *DaliDriver) QueryExtendedFadeTime(deviceAddr Address) (ExtendedFadeTime, error) {
	cmd := NewQueryExtendedFadeTimeCmd(deviceAddr)
	err := drv.SendCommand([]DaliCommand{cmd})
	if err != nil {
		return ExtendedFadeTime{}, err
	}
	return cmd.GetExtendedFadeTime(), nil
}

func (drv *DaliDriver) SetFadeTime(deviceAddr Address, fadeTime float64) error {
	// T=0.5(sqrt(pow(2,DTR))) seconds, 已知T求DTR
	dtr0 := byte(uint8(2 * math.Log2(2*fadeTime)))
	// dtr0 := byte(fadeTime)
	err := drv.SendCommand([]DaliCommand{SetDTR0(dtr0), SetFadeTimeCmd(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) SetFadeRate(deviceAddr Address, fadeRate byte) error {
	err := drv.SendCommand([]DaliCommand{SetDTR0(fadeRate), SetFadeRateCmd(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) SetExtendedFadeTime(deviceAddr Address, eft ExtendedFadeTime) error {
	dtr0 := byte(eft.ExtendedFadeTimeMultiplier<<4) | byte(eft.ExtendedFadeTimeBase-1)
	err := drv.SendCommand([]DaliCommand{SetDTR0(dtr0), SetExtendedFadeTimeCmd(deviceAddr), SetExtendedFadeTimeCmd(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) QueryGroups(deviceAddr Address) (*set.Set[int], error) {
	rst := set.New[int]()
	ztsCmd := QueryGroupsZeroToSeven(deviceAddr)
	etfCmd := QueryGroupsEightToFifteen(deviceAddr)
	err := drv.SendCommand([]DaliCommand{ztsCmd, etfCmd})
	if err != nil {
		return nil, err
	}
	for i := 0; i < 8; i++ {
		if ztsCmd.val&(1<<i) != 0 {
			rst.Set(i)
		}
	}
	for i := 0; i < 8; i++ {
		if etfCmd.val&(1<<i) != 0 {
			rst.Set(i + 8)
		}
	}
	return rst, nil
}

func (drv *DaliDriver) AddToGroup(deviceAddr []Address, group int) error {
	for _, addr := range deviceAddr {
		err := drv.SendCommand([]DaliCommand{AddToGroup(addr, group), AddToGroup(addr, group)})
		if err != nil {
			return err
		}
	}
	return nil
}

func (drv *DaliDriver) RemoveFromGroup(deviceAddr []Address, group int) error {
	for _, addr := range deviceAddr {
		err := drv.SendCommand([]DaliCommand{RemoveFromGroup(addr, group), RemoveFromGroup(addr, group)})
		if err != nil {
			return err
		}
	}
	return nil
}

func (drv *DaliDriver) QuerySceneLevel(deviceAddr Address, scene int) (int, error) {
	cmd := QuerySceneLevel(deviceAddr, scene)
	err := drv.SendCommand([]DaliCommand{cmd})
	if err != nil {
		return 0, err
	}
	return int(cmd.val), nil
}

func (drv *DaliDriver) SetSceneLevel(deviceAddr Address, scene int, level byte) error {
	err := drv.SendCommand([]DaliCommand{SetDTR0(level), SetScene(deviceAddr, scene), SetScene(deviceAddr, scene)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) RemoveFromScene(deviceAddr Address, scene int) error {
	err := drv.SendCommand([]DaliCommand{RemoveFromScene(deviceAddr, scene), RemoveFromScene(deviceAddr, scene)})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) QueryDimmingCurve(deviceAddr Address) (DimmingCurve, error) {
	cmd := QueryDimmingCurve(deviceAddr)
	err := drv.SendCommand([]DaliCommand{cmd})
	if err != nil {
		return 0, err
	}
	return DimmingCurve(cmd.val), nil
}

func (drv *DaliDriver) SelectDimmingCurve(deviceAddr Address, curve DimmingCurve) error {
	cmd := SelectDimmingCurve(deviceAddr, curve)
	err := drv.SendCommand([]DaliCommand{SetDTR0(byte(curve)), cmd})
	if err != nil {
		return err
	}
	return nil
}

func (drv *DaliDriver) IdentifyDevice(deviceAddr Address) error {
	err := drv.SendCommand([]DaliCommand{IdentifyDevice(deviceAddr)})
	if err != nil {
		return err
	}
	return nil
}

type MemoryBank0 struct {
	LastMemoryBank  int
	GTIN            uint64
	FirmwareVersion string
	SerialNumber    uint64
	HardwareVersion string

	VersionNumber101 byte
	VersionNumber102 byte
	VersionNumber103 byte

	ControlDeviceUnitNumber byte
	ControlGearUnitNumber   byte
	ControlGearUnitIndex    byte
}

func (drv *DaliDriver) ReadMemoryLocation(deviceAddr Address, bank byte) (*MemoryBank0, error) {
	cmd := ReadMemoryLocation(deviceAddr)
	err := drv.SendCommand([]DaliCommand{SetDTR1(bank), SetDTR0(0), cmd})
	if err != nil {
		return nil, err
	}
	drv.SendCommand([]DaliCommand{ReadMemoryLocation(deviceAddr)})
	length := cmd.val
	bytes := make([]byte, length+1)
	bytes[0] = length
	command := make([]DaliCommand, 0)
	for i := 2; i <= int(length); i++ {
		command = append(command, ReadMemoryLocation(deviceAddr))
	}
	err = drv.SendCommand(command)
	if err != nil {
		return nil, err
	}
	for i := 2; i <= int(length); i++ {
		bytes[i] = command[i-2].Val()
	}
	rst := &MemoryBank0{
		LastMemoryBank: int(bytes[2]),
	}

	rst.GTIN, _ = ReadBigEndianUInt64(bytes[3 : 3+6])
	rst.FirmwareVersion = fmt.Sprintf("V%d.%d", bytes[9], bytes[10])
	rst.SerialNumber, _ = ReadBigEndianUInt64(bytes[11 : 11+8])
	rst.HardwareVersion = fmt.Sprintf("V%d.%d", bytes[19], bytes[20])
	rst.VersionNumber101 = bytes[21]
	rst.VersionNumber102 = bytes[22]
	rst.VersionNumber103 = bytes[23]
	rst.ControlDeviceUnitNumber = bytes[24]
	rst.ControlGearUnitNumber = bytes[25]
	rst.ControlGearUnitIndex = bytes[26]
	return rst, nil
}
