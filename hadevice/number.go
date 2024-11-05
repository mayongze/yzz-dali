package hadevice

import (
	"encoding/json"
	"fmt"
)

type Number struct {
	CommonDefinition
	Min               float64 `json:"min,omitempty"`                 // Optional, default: 1
	Max               float64 `json:"max,omitempty"`                 // Optional, default: 100
	Mode              string  `json:"mode,omitempty"`                // Optional, default: "auto"
	Step              float64 `json:"step,omitempty"`                // Optional, default: 1
	UnitOfMeasurement string  `json:"unit_of_measurement,omitempty"` // Optional
}

func (mnc *Number) GetBytes() []byte {
	bs, _ := json.Marshal(mnc)
	return bs
}

func NewNumber(device Device, numberName, field string, min, max, step float64, unit string) *Number {
	config := &Number{
		CommonDefinition: CommonDefinition{
			Availability: []Availability{
				{
					Topic:         "yzz-dali/bridge/status",
					ValueTemplate: "{{ value_json.state }}",
				},
			},
			CommandTopic: fmt.Sprintf(`yzz-dali/%s/set/%s`, device.UID, field),
			Device: Device{
				Identifiers:  []string{fmt.Sprintf("yzz-dali_%s", device.UID)},
				Manufacturer: `LTECH`,
				Model:        "SE-12-100-500-W2D",
				Name:         device.Name,
				// ViaDevice:    "dali_bridge_yzz",
			},
			EnabledByDefault: true,
			EntityCategory:   "config",
			Icon:             "mdi:timer",
			Name:             numberName,
			ObjectID:         fmt.Sprintf("%s_%s", device.UID, field),
			UniqueID:         fmt.Sprintf("%s_%s_yzz-dali", device.UID, field),
			Origin: &Origin{
				Name: "yzz-dali",
				Sw:   "1.0.0",
				Url:  "https://yzz-dali.dz163.cn",
			},
			StateTopic:    fmt.Sprintf("yzz-dali/%s", device.UID),
			ValueTemplate: fmt.Sprintf("{{ value_json.%s }}", field),
		},
		Min:               min,
		Max:               max,
		UnitOfMeasurement: unit,
		Step:              step,
		Mode:              "box",
	}
	return config
}
