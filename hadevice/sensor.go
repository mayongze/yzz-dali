package hadevice

import (
	"encoding/json"
	"fmt"
)

type Sensor struct {
	CommonDefinition
}

func (s *Sensor) GetBytes() []byte {
	bs, _ := json.Marshal(s)
	return bs
}

func NewSensor(device Device, sensorName, field string) *Sensor {
	config := &Sensor{
		CommonDefinition: CommonDefinition{
			Availability: []Availability{
				{
					Topic:         "yzz-dali/bridge/status",
					ValueTemplate: "{{ value_json.state }}",
				},
			},
			CommandTopic:     fmt.Sprintf(`yzz-dali/%s/set/%s`, device.UID, field),
			Device:           device,
			EnabledByDefault: true,
			EntityCategory:   "diagnostic",
			Icon:             "mdi:gesture-double-tap",
			Name:             sensorName,
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
	}
	return config
}
