package homeassistant

import (
	"fmt"
	"huiqun/homeassistant/hadevice"
	"log"
	"strconv"
	"time"
)

type GearShortEntity struct {
	Addr  int
	UID   string
	Name  string
	Init  bool
	State *LightState
}

func (g *GearShortEntity) GetAddr() Address {
	return ShortAddr(g.Addr)
}

func (g *GearShortEntity) GetDeviceType() DeviceType {
	return GearShortDeviceType
}

func (s *GearShortEntity) getHaDevice() hadevice.Device {
	return hadevice.Device{
		Identifiers:  []string{fmt.Sprintf("yzz-dali_%s", s.UID)},
		Manufacturer: `LTECH`,
		Model:        "SE-12-100-500-W2D",
		Name:         s.Name,
		UID:          s.UID,
		// ViaDevice:    "dali_bridge_yzz",
	}
}

func (g *GearShortEntity) InitDiscover() error {
	topic := fmt.Sprintf("homeassistant/light/%s/light/config", g.UID)
	lightCfg := hadevice.NewLight(g.getHaDevice(), g.Name)
	token := GetMqttCli().Publish(topic, 0, true, lightCfg.GetBytes())
	if token.Error() != nil {
		return token.Error()
	}
	selectCfg := hadevice.NewSelect(g.getHaDevice(), "调光曲线", "dimming_curve",
		[]string{"standard", "linear"})
	token = GetMqttCli().Publish(fmt.Sprintf("homeassistant/select/%s/%s/config", g.UID, "dimming_curve"),
		0, true, selectCfg.GetBytes())
	if token.Error() != nil {
		return token.Error()
	}
	transitionCfg := hadevice.NewNumber(g.getHaDevice(), "过渡时间",
		"transition", 0.1, 16, 0.1, "Second")
	token = GetMqttCli().Publish(fmt.Sprintf("homeassistant/number/%s/%s/config", g.UID, "transition"),
		0, true, transitionCfg.GetBytes())
	if token.Error() != nil {
		return token.Error()
	}
	fadeRateCfg := hadevice.NewNumber(g.getHaDevice(), "调光速率",
		"fade_rate", 1, 15, 1, "")
	token = GetMqttCli().Publish(fmt.Sprintf("homeassistant/number/%s/%s/config", g.UID, "fade_rate"),
		0, true, fadeRateCfg.GetBytes())
	if token.Error() != nil {
		return token.Error()
	}

	UpBtnCfg := hadevice.NewButton(g.getHaDevice(), "Dim Up", "command", "up")
	token = GetMqttCli().Publish(fmt.Sprintf("homeassistant/button/%s/%s/config", g.UID, "up"), 0, true, UpBtnCfg.GetBytes())
	if token.Error() != nil {
		return token.Error()
	}

	DownBtnCfg := hadevice.NewButton(g.getHaDevice(), "Dim Down", "command", "down")
	token = GetMqttCli().Publish(fmt.Sprintf("homeassistant/button/%s/%s/config", g.UID, "down"), 0, true, DownBtnCfg.GetBytes())
	if token.Error() != nil {
		return token.Error()
	}

	stepUpBtnCfg := hadevice.NewButton(g.getHaDevice(), "Step Up", "command", "step_up")
	token = GetMqttCli().Publish(fmt.Sprintf("homeassistant/button/%s/%s/config", g.UID, "step_up"), 0, true, stepUpBtnCfg.GetBytes())
	if token.Error() != nil {
		return token.Error()
	}

	stepDownBtnCfg := hadevice.NewButton(g.getHaDevice(), "Step Down", "command", "step_down")
	token = GetMqttCli().Publish(fmt.Sprintf("homeassistant/button/%s/%s/config", g.UID, "step_down"), 0, true, stepDownBtnCfg.GetBytes())
	if token.Error() != nil {
		return token.Error()
	}

	// todo: 查询
	SceneCfg := hadevice.NewScene(g.getHaDevice(), "Scene 1", "scene", "1")
	token = GetMqttCli().Publish(fmt.Sprintf("homeassistant/button/%s/%s/config", g.UID, "scene1"), 0, true, SceneCfg.GetBytes())
	if token.Error() != nil {
		return token.Error()
	}

	IdentifyDeviceCfg := hadevice.NewButton(g.getHaDevice(), "Identify Device", "command", "identify")
	IdentifyDeviceCfg.DeviceClass = "identify"
	token = GetMqttCli().Publish(fmt.Sprintf("homeassistant/button/%s/%s/config", g.UID, "identify"), 0, true, IdentifyDeviceCfg.GetBytes())
	if token.Error() != nil {
		return token.Error()
	}

	return g.QueryState()
}

func (g *GearShortEntity) UpdateState(kv StateKV) error {
	g.State.UpdateValue(kv)
	token := GetMqttCli().Publish(fmt.Sprintf("yzz-dali/%s", g.UID), 0, false, g.State.GetBytes())
	return token.Error()
}

func (g *GearShortEntity) HandleCommand(stateKV StateKV) error {
	// OFF -> ON
	// OFF -> ON brightness
	// OFF -> ON color_temp
	//     -> OFF
	var err error
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
		if ok && g.State.ColorTemp != int(colorTemp) {
			err := drv.SetDT8ColourValueTc(g.GetAddr(), int(colorTemp))
			if err != nil {
				log.Print("SetColorTemp error: ", err)
			}
		}
	}

	command, ok := GetStateV[string](stateKV, "command")
	if ok {
		switch command {
		case "up":
			if err = drv.Up(g.GetAddr()); err != nil {
				log.Print("SetLevel error: ", err)
			}
		case "step_up":
			if err = drv.OnAndStepUp(g.GetAddr()); err != nil {
				log.Print("SetLevel error: ", err)
			} else {
				if g.State.State == "OFF" {
					stateKV["state"] = "ON"
					stateKV["brightness"] = float64(1)
				} else if g.State.Brightness < 254 {
					stateKV["brightness"] = float64(g.State.Brightness + 1)
				}
			}
		case "step_down":
			if err = drv.StepDownAndOff(g.GetAddr()); err != nil {
				log.Print("SetLevel error: ", err)
			} else {
				if g.State.Brightness == 1 {
					stateKV["state"] = "OFF"
				}
				if g.State.Brightness > 1 {
					stateKV["brightness"] = float64(g.State.Brightness - 1)
				}
			}
		case "down":
			if err = drv.Down(g.GetAddr()); err != nil {
				log.Print("SetLevel error: ", err)
			}
		case "identify":
			if err = drv.IdentifyDevice(g.GetAddr()); err != nil {
				log.Print("IdentifyDevice error: ", err)
			}
		}

		if command == "up" || command == "down" {
			// 修正亮度值
			time.AfterFunc(200*time.Millisecond, func() {
				kv := make(StateKV)
				level, err := drv.QueryActualLevel(g.GetAddr())
				if err != nil {
					log.Print("QueryActualLevel error: ", err)
				} else {
					kv["brightness"] = float64(level)
					_ = g.UpdateState(kv)
				}
			})
		}
	}

	scene, ok := GetStateV[string](stateKV, "scene")
	if ok {
		sceneId, _ := strconv.Atoi(scene)
		if err = drv.GoToScene(g.GetAddr(), sceneId); err != nil {
			log.Print("GoToScene error: ", err)
		}
		level, err := drv.QuerySceneLevel(g.GetAddr(), sceneId)
		if err != nil {
			log.Print("QuerySceneLevel error: ", err)
		}
		stateKV["brightness"] = float64(level)
	}

	dimmingCurve, ok := GetStateV[string](stateKV, "dimming_curve")
	if ok {
		if err = drv.SelectDimmingCurve(g.GetAddr(), ParseDimmingCurve(dimmingCurve)); err != nil {
			log.Print("SetDimmingCurve error: ", err)
		}
	}

	fadeRateStr, ok := GetStateV[string](stateKV, "fade_rate")
	if ok {
		if err = drv.SetFadeRate(g.GetAddr(), MustParseByte(fadeRateStr)); err != nil {
			log.Print("SetFadeRate error: ", err)
		}
	}

	transitionStr, ok := GetStateV[string](stateKV, "transition")
	if ok {
		ftfr, err := drv.QueryFadeTimeFadeRate(g.GetAddr())
		if err != nil {
			return err
		}
		if ftfr.FadeTime != 0 {
			if err = drv.SetFadeTime(g.GetAddr(), 0); err != nil {
				return err
			}
		}
		transition := MustParseFloat64(transitionStr)
		param := ExtendedFadeTime{}
		if transition <= 1.6 {
			param.ExtendedFadeTimeMultiplier = ExtendedFadeTimeMultiplier100Millisecond
			param.ExtendedFadeTimeBase = int(transition * 10)
		} else {
			param.ExtendedFadeTimeMultiplier = ExtendedFadeTimeMultiplier1Second
			param.ExtendedFadeTimeBase = int(transition)
		}
		err = drv.SetExtendedFadeTime(g.GetAddr(), param)
		if err != nil {
			log.Print("SetExtendedFadeTime error: ", err)
		}
	}
	return g.UpdateState(stateKV)
}

func (g *GearShortEntity) QueryState() error {
	if g.Init {
		return nil
	}
	drv := GetDalDriver()
	mem, err := drv.ReadMemoryLocation(g.GetAddr(), 0)
	if err != nil {
		return err
	}
	_ = mem
	level, err := drv.QueryActualLevel(g.GetAddr())
	if err != nil {
		return err
	}
	g.State.Brightness = level
	if level != 0 {
		g.State.LastNonZeroBrightness = g.State.Brightness
		g.State.State = "ON"
	} else {
		g.State.State = "OFF"
	}
	g.State.ColorMode = "color_temp"
	ftfr, err := drv.QueryFadeTimeFadeRate(g.GetAddr())
	if err != nil {
		return err
	}
	if ftfr.FadeTime != 0 {
		g.State.Transition = ftfr.FadeTime
	} else {
		eft, err := drv.QueryExtendedFadeTime(g.GetAddr())
		if err != nil {
			return err
		}
		if eft.Millisecond() != 0 {
			g.State.Transition = float64(eft.Millisecond()) / 1000.0
		}
	}
	g.State.FadeRate = ftfr.FadeRate
	dimmingCurve, err := drv.QueryDimmingCurve(g.GetAddr())
	if err != nil {
		return err
	}
	g.State.DimmingCurve = dimmingCurve.String()
	// 发布mqtt
	token := GetMqttCli().Publish(fmt.Sprintf("yzz-dali/%s", g.UID), 0, false, g.State.GetBytes())
	if token.Error() != nil {
		return token.Error()
	}
	g.Init = true
	return nil
}
