package homeassistant

import (
	"fmt"
	"huiqun/homeassistant/hadevice"
	"log"
)

type GroupEntity struct {
	Addr    int
	UID     string
	Name    string
	State   *LightState
	Devices []GearShortEntity
}

func (g *GroupEntity) GetAddr() Address {
	return GroupAddr(g.Addr)
}

func (g *GroupEntity) GetDeviceType() DeviceType {
	return GroupDeviceType
}

func (s *GroupEntity) getHaDevice() hadevice.Device {
	return hadevice.Device{
		Identifiers:  []string{fmt.Sprintf("yzz-dali_%s", s.UID)},
		Manufacturer: `DALI-YZz`,
		Model:        "dali group",
		Name:         s.Name,
		UID:          s.UID,
		// ViaDevice:    "dali_bridge_yzz",
	}
}

func (g *GroupEntity) InitDiscover() error {
	topic := fmt.Sprintf("homeassistant/light/%s/light/config", g.UID)
	lightCfg := hadevice.NewLight(g.getHaDevice(), g.Name)
	token := GetMqttCli().Publish(topic, 0, true, lightCfg.GetBytes())
	if token.Error() != nil {
		return token.Error()
	}
	return g.QueryState()
}

func (g *GroupEntity) QueryState() error {
	// 发布mqtt
	token := GetMqttCli().Publish(fmt.Sprintf("yzz-dali/%s", g.UID), 0, false, g.State.GetBytes())
	if token.Error() != nil {
		return token.Error()
	}
	// 任意状态是关的都设置当前为关闭状态
	for _, device := range g.Devices {
		if !device.Init {
			if err := device.QueryState(); err != nil {
				return err
			}
		}
		if device.State.State == "OFF" {
			g.State.State = "OFF"
		}
	}
	return nil
}

func (g *GroupEntity) UpdateState(kv StateKV) error {
	g.State.UpdateValue(kv)
	token := GetMqttCli().Publish(fmt.Sprintf("yzz-dali/%s", g.UID), 0, false, g.State.GetBytes())
	return token.Error()
}

func (g *GroupEntity) HandleCommand(stateKV StateKV) error {
	drv := GetDalDriver()
	state, ok := GetStateV[string](stateKV, "state")
	if ok && state == "Toggle" {
		if g.State.State == "ON" {
			state = "OFF"
		} else {
			state = "ON"
		}
		stateKV["state"] = state
	}

	if ok && state == "OFF" {
		err := drv.SetLevel(g.GetAddr(), 0)
		if err != nil {
			log.Print("SetPower error: ", err)
		}
	}
	if ok && state == "ON" {
		brightness, ok := GetStateV[float64](stateKV, "brightness")
		if ok {
			err := drv.SetLevel(g.GetAddr(), int(brightness))
			if err != nil {
				log.Print("SetLevel error: ", err)
			}
		} else if g.State.State == "OFF" {
			err := drv.SetLevel(g.GetAddr(), g.State.LastNonZeroBrightness)
			if err != nil {
				log.Print("SetLevel error: ", err)
			}
		}
		colorTemp, ok := GetStateV[float64](stateKV, "color_temp")
		if ok {
			err := drv.SetDT8ColourValueTc(g.GetAddr(), int(colorTemp))
			if err != nil {
				log.Print("SetColorTemp error: ", err)
			}
		}
	}
	if err := g.UpdateState(stateKV); err != nil {
		return err
	}
	for _, device := range g.Devices {
		if err := device.UpdateState(stateKV); err != nil {
			return err
		}
	}
	return nil
}
