package homeassistant

import (
	"encoding/hex"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"huiqun/homeassistant/adapter"
	"io"
	"log"
	"testing"
	"time"
)

const sss string = "1"

var daliSwitchSeq = [3]byte{0xA0, 0xAD, 0xFB}

func TestNewSaveKeyProtocol(t *testing.T) {
	//host := ""
	//port := 888
	//tcpClient := adapter.NewTcpClientAdapter(host, port)

	ctrl := gomock.NewController(t)
	tcpClientMock := adapter.NewMockBusAdapter(ctrl)
	defer ctrl.Finish()

	reader, writer := io.Pipe()
	expect := "118e01190f110001a0adfb0021ff06000000000048"
	expectWriteArgs, _ := hex.DecodeString(expect)
	tcpClientMock.EXPECT().Write(gomock.Eq(expectWriteArgs)).Return(len(expectWriteArgs), nil).Do(func(interface{}) {
		resp, _ := hex.DecodeString("110e01190f110001a0adfb0121ff050000000000c8")
		_, _ = writer.Write(resp)
	}).AnyTimes()
	tcpClientMock.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
		return reader.Read(p)
	}).AnyTimes()
	sBus := NewSBus(tcpClientMock)
	skDlDriver := NewSaveKeyDALIDriver(NewSaveKeyDevice(sBus, daliSwitchSeq, 0x01))

	daliCommand := RecallMinLevel(BroadcastAddr())
	err := skDlDriver.SendInstructionSet([]uint16{daliCommand.instruction})
	assert.Equal(t, nil, err)

	time.Sleep(10 * time.Second)
}

func TestYzz(t *testing.T) {
	host := "192.168.123.129"
	port := "8899"
	tcpClient := adapter.NewTcpClientAdapter(host, port)
	skDlDriver := NewDaliDriver(NewSaveKeyDALIDriver(NewSaveKeyDevice(NewSBus(tcpClient), daliSwitchSeq, 0x01)))

	//serialPortName := "/dev/cu.usbserial-1140"
	// baud := 19200
	//serialPort := adapter.NewSerialAdapter(serialPortName, baud)
	//sBus := NewSBus(serialPort)
	//skDlDriver := NewDaliDriver(NewSaveKeyDALIDriver(NewSaveKeyDevice(sBus, daliSwitchSeq, 0x01)))

	//types, err := skDlDriver.QueryDeviceType(ShortAddr(2))
	//assert.Equal(t, nil, err, types)
	var err error
	time.Sleep(time.Second)
	addr := ShortAddr(3)

	err = skDlDriver.SetFastFadeTime(addr, 23)
	assert.Equal(t, nil, err)

	t0 := time.Now()

	//err = skDlDriver.SetDT8ColourValueTc(ShortAddr(2), TcKelvinMirek(2700))
	//assert.Equal(t, nil, err)
	//
	//err = skDlDriver.SetPowerOnLevel(ShortAddr(2), 100)
	//assert.Equal(t, nil, err)
	//
	//err = skDlDriver.SetLevel(ShortAddr(2), 100)
	//assert.Equal(t, nil, err)

	tcMirek, err := skDlDriver.QueryDT8ColourValue(addr, ColourTemperatureTC)
	assert.Equal(t, nil, err, tcMirek)
	tcCoolest, err := skDlDriver.QueryDT8ColourValue(addr, ColourTemperatureTcCoolest)
	assert.Equal(t, nil, err, tcCoolest)
	tcPhysicalCoolest, err := skDlDriver.QueryDT8ColourValue(addr, ColourTemperatureTcPhysicalCoolest)
	assert.Equal(t, nil, err, tcPhysicalCoolest)
	tcWarmest, err := skDlDriver.QueryDT8ColourValue(addr, ColourTemperatureTcWarmest)
	assert.Equal(t, nil, err, tcWarmest)
	tcPhysicalWarmest, err := skDlDriver.QueryDT8ColourValue(addr, ColourTemperatureTcPhysicalWarmest)
	assert.Equal(t, nil, err, tcPhysicalWarmest)
	log.Printf("tcMirek: %dK, tcCoolest: %dK, tcPhysicalCoolest: %dK, tcWarmest: %dK, tcPhysicalWarmest: %dK",
		TcKelvinMirek(tcMirek), TcKelvinMirek(tcCoolest), TcKelvinMirek(tcPhysicalCoolest),
		TcKelvinMirek(tcWarmest), TcKelvinMirek(tcPhysicalWarmest))

	ftfr, err := skDlDriver.QueryFadeTimeFadeRate(addr)
	assert.Equal(t, nil, err, ftfr)
	fastFadeTime, err := skDlDriver.QueryFastFadeTime(addr)
	assert.Equal(t, nil, err, fastFadeTime)
	minFastFadeTime, err := skDlDriver.QueryMinFastFadeTime(addr)
	assert.Equal(t, nil, err, minFastFadeTime)
	fff, err := skDlDriver.QueryExtendedFadeTime(addr)
	assert.Equal(t, nil, err, fff)

	minLevel, err := skDlDriver.QueryMinLevel(addr)
	assert.Equal(t, nil, err, minLevel)
	maxLevel, err := skDlDriver.QueryMaxLevel(addr)
	assert.Equal(t, nil, err, maxLevel)

	log.Printf("time: %dms", time.Since(t0).Milliseconds())
	// for i := 0; i < 5; i++ {
	//daliCommand := RecallMaxLevel(ShortAddr(2))
	//err := skDlDriver.SendInstructionSet([]uint16{daliCommand})
	//assert.Equal(t, nil, err)
	//
	//time.Sleep(2 * time.Second)
	//daliCommand = RecallMinLevel(ShortAddr(2))
	//err = skDlDriver.SendInstructionSet([]uint16{daliCommand})
	//assert.Equal(t, nil, err)

	time.Sleep(2 * time.Second)
	// }
}

func TestQueryInfo(t *testing.T) {
	host := "192.168.123.129"
	port := "8899"
	tcpClient := adapter.NewTcpClientAdapter(host, port)
	sBus := NewSBus(tcpClient)
	skDlDriver := NewDaliDriver(NewSaveKeyDALIDriver(NewSaveKeyDevice(sBus, daliSwitchSeq, 0x01)))
	addr := ShortAddr(2)

	var err error

	err = skDlDriver.SetLevel(addr, 0)

	//err = skDlDriver.SetFadeTime(addr, 0)
	//assert.Equal(t, nil, err)
	ftfr, err := skDlDriver.QueryFadeTimeFadeRate(addr)
	assert.Equal(t, nil, err, ftfr)

	err = skDlDriver.SetExtendedFadeTime(addr, ExtendedFadeTime{
		ExtendedFadeTimeBase:       12,
		ExtendedFadeTimeMultiplier: ExtendedFadeTimeMultiplier10Seconds})
	assert.Equal(t, nil, err)
	fff, err := skDlDriver.QueryExtendedFadeTime(addr)
	assert.Equal(t, nil, err, fff)
}

func TestCmd(t *testing.T) {
	host := "192.168.123.129"
	port := "8899"
	tcpClient := adapter.NewTcpClientAdapter(host, port)
	sBus := NewSBus(tcpClient)
	skDlDriver := NewDaliDriver(NewSaveKeyDALIDriver(NewSaveKeyDevice(sBus, daliSwitchSeq, 0x01)))
	addr := ShortAddr(3)

	var err error

	err = skDlDriver.GoToLastActiveLevel(addr)
	assert.Equal(t, nil, err)
	err = skDlDriver.SetLevel(addr, 100)
	assert.Equal(t, nil, err)
	time.Sleep(time.Second)
	err = skDlDriver.GoToLastActiveLevel(addr)
	assert.Equal(t, nil, err)

	err = skDlDriver.SetLevel(addr, 0)
	assert.Equal(t, nil, err)

	err = skDlDriver.RecallMaxLevel(addr)
	assert.Equal(t, nil, err)
	time.Sleep(time.Second)
	err = skDlDriver.RecallMinLevel(addr)
	assert.Equal(t, nil, err)
	time.Sleep(time.Second)
	err = skDlDriver.SetLevel(addr, 100)
	assert.Equal(t, nil, err)
	time.Sleep(time.Second)
	err = skDlDriver.SetLevel(addr, 254)
	assert.Equal(t, nil, err)
	time.Sleep(time.Second)
	err = skDlDriver.Off(addr)
	assert.Equal(t, nil, err)
	time.Sleep(time.Second)
	err = skDlDriver.SetLevel(addr, 255)
	assert.Equal(t, nil, err)
	//err = skDlDriver.GoToLastActiveLevel(addr)
	//assert.Equal(t, nil, err)
	time.Sleep(5 * time.Second)
	err = skDlDriver.OnAndStepUp(addr)
	assert.Equal(t, nil, err)
	time.Sleep(time.Second)
	err = skDlDriver.OnAndStepUp(addr)
	assert.Equal(t, nil, err)
	time.Sleep(time.Second)
	err = skDlDriver.Up(addr)
	assert.Equal(t, nil, err)
	time.Sleep(time.Second)
	err = skDlDriver.SetLevel(addr, 0)
	assert.Equal(t, nil, err)
}

func TestScanShortAddr(t *testing.T) {
	host := "192.168.123.129"
	port := "8899"
	tcpClient := adapter.NewTcpClientAdapter(host, port)
	sBus := NewSBus(tcpClient)
	skDlDriver := NewDaliDriver(NewSaveKeyDALIDriver(NewSaveKeyDevice(sBus, daliSwitchSeq, 0x01)))

	for i := 0; i < 64; i++ {
		addr := ShortAddr(i)
		types, err := skDlDriver.QueryDeviceType(addr)
		if errors.Is(err, ErrNoResponse) {
			continue
		}
		assert.Equal(t, nil, err)
		t.Logf("addr: %d, type: %v", i, types)
	}
}

func TestGroup(t *testing.T) {
	host := "192.168.123.129"
	port := "8899"
	tcpClient := adapter.NewTcpClientAdapter(host, port)
	sBus := NewSBus(tcpClient)
	skDlDriver := NewDaliDriver(NewSaveKeyDALIDriver(NewSaveKeyDevice(sBus, daliSwitchSeq, 0x01)))

	var err error
	groupAddr := GroupAddr(1)
	device0 := ShortAddr(0)
	device1 := ShortAddr(1)
	device2 := ShortAddr(2)
	device3 := ShortAddr(3)
	_ = groupAddr

	// skDlDriver.SetLevel(groupAddr, 0)
	//err = skDlDriver.SetSceneLevel(BroadcastAddr(), 1, 0)
	//assert.Equal(t, nil, err)
	//err = skDlDriver.GoToScene(BroadcastAddr(), 2)
	//l, err := skDlDriver.QuerySceneLevel(device2, 0)
	//assert.Equal(t, nil, err, l)
	//err = skDlDriver.SetExtendedFadeTime(groupAddr, ExtendedFadeTime{
	//	ExtendedFadeTimeBase:       1,
	//	ExtendedFadeTimeMultiplier: ExtendedFadeTimeMultiplier1Second})
	//assert.Equal(t, nil, err)

	rsp, err := skDlDriver.QueryGroups(device1)
	list := rsp.ToSlice()
	assert.Equal(t, nil, err, list)

	err = skDlDriver.AddToGroup([]Address{device0, device1, device2, device3}, 1)
	assert.Equal(t, nil, err)

	// 循环16个地址 16个场景
}
