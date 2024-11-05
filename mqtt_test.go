package homeassistant

import (
	"huiqun/homeassistant/adapter"
	"testing"
	"time"
)

func TestMMMM(t *testing.T) {
	InitSysConfig("yzz-dali.yml")
	InitGlobalState("state.js")
	host := "192.168.123.129"
	port := "8899"
	tcpClient := adapter.NewTcpClientAdapter(host, port)
	sBus := NewSBus(tcpClient)

	InitGlobalDeviceList(sBus)
	InitDaliDriver(NewSaveKeyDALIDriver(NewSaveKeyDevice(sBus, daliSwitchSeq, 0x01)))
	RegisterOnConnect(SubscriptionCommandSet)
	InitMQTTClient(Cfg.Mqtt)

	for _, device := range GlobalDeviceList {
		err := device.InitDiscover()
		if err != nil {
			t.Logf("InitDiscover error: %v", err)
		}
	}
	// 每隔10s写出一次状态
	go func() {
		for range time.NewTicker(10 * time.Second).C {
			SaveGlobalState("state.js")
		}
	}()
	time.Sleep(50 * time.Minute)
}
