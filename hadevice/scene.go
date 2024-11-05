package hadevice

import (
	"encoding/json"
	"fmt"
)

type Scene struct {
	CommonDefinition
	PayloadOn string `json:"payload_on,omitempty"`
}

func (mlc *Scene) GetBytes() []byte {
	bs, _ := json.Marshal(mlc)
	return bs
}

func NewScene(device Device, sceneName, field string, payloadOn string) *Scene {
	config := &Scene{
		CommonDefinition: CommonDefinition{
			CommandTopic:     fmt.Sprintf(`yzz-dali/%s/set/%s`, device.UID, field),
			Device:           device,
			EnabledByDefault: true,
			DeviceClass:      "scene",
			Name:             sceneName,
			ObjectID:         fmt.Sprintf("%s_scene_%s", device.UID, field),
			UniqueID:         fmt.Sprintf("%s_scene_%s_yzz-dali", device.UID, field),
			Origin: &Origin{
				Name: "yzz-dali",
				Sw:   "1.0.0",
				Url:  "https://yzz-dali.dz163.cn",
			},
		},
		PayloadOn: payloadOn,
	}
	return config
}
