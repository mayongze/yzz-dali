package hadevice

import (
	"encoding/json"
	"fmt"
)

type Select struct {
	CommonDefinition
	Options []string `json:"options"` // Required
}

func (mlc *Select) GetBytes() []byte {
	bs, _ := json.Marshal(mlc)
	return bs
}

func NewSelect(device Device, selectName, field string, options []string) *Select {
	config := &Select{
		CommonDefinition: CommonDefinition{
			Availability: []Availability{
				{
					Topic:         "yzz-dali/bridge/status",
					ValueTemplate: "{{ value_json.state }}",
				},
			},
			CommandTopic: fmt.Sprintf(`yzz-dali/%s/set/%s`, device.UID, field),
			// CommandTemplate: `{{ state_attr('this','options').index(states('this')) }}`,
			// https://community.home-assistant.io/t/mqtt-select-option-command-template-by-index/647094
			Device:           device,
			EnabledByDefault: true,
			EntityCategory:   "config",
			Icon:             "mdi:tune",
			Name:             selectName,
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
		Options: options,
	}
	return config
}
