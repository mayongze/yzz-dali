package homeassistant

import (
	"encoding/json"
	"log"
	"os"
)

type StateKV map[string]interface{}

func GetStateV[T string | int | float64](stateKV StateKV, key string) (T, bool) {
	value, ok := stateKV[key]
	if !ok {
		return *(new(T)), ok
	}
	return value.(T), ok
}

type State interface {
	GetBytes() []byte
	UpdateValue(kv StateKV)
}

type StateStore map[string]interface{}

func (s StateStore) GetLightState(key string) *LightState {
	var state = &LightState{
		LastNonZeroBrightness: 254,
	}
	if v, ok := s[key]; ok {
		switch v.(type) {
		case *LightState:
			return v.(*LightState)
		case map[string]interface{}:
			// json覆盖state
			vv, _ := json.Marshal(v)
			_ = json.Unmarshal(vv, state)
		}
	}
	s[key] = state
	return state
}

func (s StateStore) GetSwitchState(key string) *SwitchState {
	state := &SwitchState{
		L1: "OFF",
		L2: "OFF",
		L3: "OFF",
		L4: "OFF",
	}
	if v, ok := s[key]; ok {
		switch v.(type) {
		case *SwitchState:
			return v.(*SwitchState)
		case map[string]interface{}:
			// json覆盖state
			vv, _ := json.Marshal(v)
			_ = json.Unmarshal(vv, state)
		}
	}
	s[key] = state
	return state
}

var GlobalState = make(StateStore)

func InitGlobalState(stateFile string) {
	bs, err := os.ReadFile(stateFile)
	if err != nil {
		log.Fatal(err)
	}
	if len(bs) == 0 {
		return
	}
	if err = json.Unmarshal(bs, &GlobalState); err != nil {
		log.Fatal(err)
	}
}

func SaveGlobalState(stateFile string) {
	bs, err := json.Marshal(GlobalState)
	if err != nil {
		log.Fatal(err)
	}
	if len(bs) == 0 {
		return
	}
	if err = os.WriteFile(stateFile, bs, 0644); err != nil {
		log.Fatal(err)
	}
}

type RGBColor struct {
	R int     `json:"r"`
	G int     `json:"g"`
	B int     `json:"b"`
	C int     `json:"c"`
	W int     `json:"w"`
	X float64 `json:"x"`
	Y float64 `json:"y"`
	H float64 `json:"h"`
	S float64 `json:"s"`
}

type LightState struct {
	Brightness            int       `json:"brightness,omitempty"`
	LastNonZeroBrightness int       `json:"last_non_zero_brightness,omitempty"`
	ColorMode             string    `json:"color_mode,omitempty"`
	ColorTemp             int       `json:"color_temp,omitempty"`
	Color                 *RGBColor `json:"color,omitempty"`
	Effect                string    `json:"effect,omitempty"`
	State                 string    `json:"state,omitempty"`
	DimmingCurve          string    `json:"dimming_curve,omitempty"`
	FadeRate              int       `json:"fade_rate,omitempty"`
	Transition            float64   `json:"transition,omitempty"`
}

func (l *LightState) UpdateValue(kv StateKV) {
	brightness, ok := GetStateV[float64](kv, "brightness")
	if ok {
		if brightness != 0 {
			l.LastNonZeroBrightness = l.Brightness
		}
		l.Brightness = int(brightness)
	}
	if v, ok := kv["color_mode"]; ok {
		l.ColorMode = v.(string)
	}
	if v, ok := kv["color_temp"]; ok {
		l.ColorTemp = int(v.(float64))
	}
	if v, ok := kv["state"]; ok {
		l.State = v.(string)
	}
	if v, ok := kv["dimming_curve"]; ok {
		l.DimmingCurve = v.(string)
	}
	if v, ok := kv["transition"]; ok {
		l.Transition = MustParseFloat64(v.(string))
	}
	if v, ok := kv["fade_rate"]; ok {
		l.FadeRate = MustParseInt(v.(string))
	}
}

func (l *LightState) GetBytes() []byte {
	bs, _ := json.Marshal(l)
	return bs
}

type SwitchState struct {
	L1      string `json:"l1"`
	StateL1 string `json:"state_l1"`
	L2      string `json:"l2"`
	StateL2 string `json:"state_l2"`
	L3      string `json:"l3"`
	StateL3 string `json:"state_l3"`
	L4      string `json:"l4"`
	StateL4 string `json:"state_l4"`
	Action  string `json:"action"`
}

func (s *SwitchState) GetBytes() []byte {
	bs, _ := json.Marshal(s)
	return bs
}

func (s *SwitchState) UpdateValue(kv StateKV) {
	l1, ok := GetStateV[string](kv, "l1")
	if ok {
		s.L1 = l1
	}
	l2, ok := GetStateV[string](kv, "l2")
	if ok {
		s.L2 = l2
	}
	l3, ok := GetStateV[string](kv, "l3")
	if ok {
		s.L3 = l3
	}
	l4, ok := GetStateV[string](kv, "l4")
	if ok {
		s.L4 = l4
	}
}
