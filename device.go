package homeassistant

import (
	"log"
)

type Device interface {
	GetDeviceType() DeviceType
	InitDiscover() error
	HandleCommand(stateKV StateKV) error
}

type BroadcastEntity struct {
}

type DeviceType string

const (
	GearShortDeviceType DeviceType = "GearShort"
	GroupDeviceType     DeviceType = "DaliGroup"
	StrDanceDeviceType  DeviceType = "StrDance"
)

type DeviceList map[string]Device

func (d DeviceList) GetGearDevice(uid string) *GearShortEntity {
	device, ok := GlobalDeviceList[uid].(*GearShortEntity)
	if !ok {
		return nil
	}
	return device
}

var GlobalDeviceList = make(DeviceList)
var GlobalStrDanceDeviceList = make(map[string]*StrDanceDevice)

func InitGlobalDeviceList(sBus *SBus) {
	for uid, device := range Cfg.Devices {
		if device.Type == GearShortDeviceType {
			GlobalDeviceList[uid] = &GearShortEntity{
				Addr:  device.Addr,
				UID:   uid,
				Name:  device.Name,
				State: GlobalState.GetLightState(uid),
			}
		}
		if device.Type == StrDanceDeviceType {
			GlobalStrDanceDeviceList[uid] = NewStrDanceDevice(sBus, device.Addr, uid, device.Name)
		}
	}
	for uid, group := range Cfg.Groups {
		groupDevice := &GroupEntity{
			Addr:  group.Addr,
			UID:   uid,
			Name:  group.Name,
			State: GlobalState.GetLightState(uid),
		}
		for _, deviceUID := range group.Devices {
			gearDevice := GlobalDeviceList.GetGearDevice(deviceUID)
			if gearDevice == nil {
				log.Fatal("device not found or type error: ", deviceUID)
			}
			groupDevice.Devices = append(groupDevice.Devices, *gearDevice)
		}
		GlobalDeviceList[uid] = groupDevice
	}
}
