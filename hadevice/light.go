package hadevice

import (
	"encoding/json"
	"fmt"
)

// Light represents the configuration for an MQTT JSON light
type Light struct {
	CommonDefinition
	Schema              string   `json:"schema,omitempty"`                // Must be set to "json"
	Brightness          bool     `json:"brightness,omitempty"`            // Optional, default: false
	BrightnessScale     int      `json:"brightness_scale,omitempty"`      // Optional, default: 255
	Effect              bool     `json:"effect,omitempty"`                // Optional, default: false
	EffectList          []string `json:"effect_list,omitempty"`           // Optional
	FlashTimeShort      int      `json:"flash_time_short,omitempty"`      // Optional, default: 2
	FlashTimeLong       int      `json:"flash_time_long,omitempty"`       // Optional, default: 10
	MaxMireds           int      `json:"max_mireds,omitempty"`            // Optional
	MinMireds           int      `json:"min_mireds,omitempty"`            // Optional
	SupportedColorModes []string `json:"supported_color_modes,omitempty"` // Optional
	WhiteScale          int      `json:"white_scale,omitempty"`           // Optional, default: 255
}

func (mlc *Light) GetBytes() []byte {
	bs, _ := json.Marshal(mlc)
	return bs
}

func NewLight(device Device, lightName string) *Light {
	config := &Light{
		CommonDefinition: CommonDefinition{
			Availability: []Availability{
				{
					Topic:         "yzz-dali/bridge/status",
					ValueTemplate: "{{ value_json.state }}",
				},
			},
			StateTopic:   fmt.Sprintf("yzz-dali/%s", device.UID),
			CommandTopic: fmt.Sprintf(`yzz-dali/%s/set`, device.UID),
			Device:       device,
			Name:         lightName,
			ObjectID:     fmt.Sprintf("%s_%s", device.UID, "light"),
			UniqueID:     fmt.Sprintf("%s_%s_yzz-dali", device.UID, "light"),
			Origin: &Origin{
				Name: "yzz-dali",
				Sw:   "1.0.0",
				Url:  "https://yzz-dali.dz163.cn",
			},
		},
		Schema:              "json",
		Brightness:          true,
		BrightnessScale:     254,
		MaxMireds:           370,
		MinMireds:           153,
		SupportedColorModes: []string{"color_temp"},
	}
	return config
}
