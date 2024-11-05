package hadevice

import (
	"encoding/json"
	"fmt"
)

type Switch struct {
	CommonDefinition
	PayloadOff string `json:"payload_off"`
	PayloadOn  string `json:"payload_on"`
}

func (s *Switch) GetBytes() []byte {
	bs, _ := json.Marshal(s)
	return bs
}

func NewSwitch(device Device, switchName, field string, payloadOn, payloadOff string) *Switch {
	config := &Switch{
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
			Icon:             "mdi:light-switch",
			Name:             switchName,
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
		PayloadOff: payloadOff,
		PayloadOn:  payloadOn,
	}
	return config
}
