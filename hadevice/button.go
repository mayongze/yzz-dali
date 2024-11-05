package hadevice

import (
	"encoding/json"
	"fmt"
)

type Button struct {
	CommonDefinition
	PayloadPress string `json:"payload_press"`
}

func (mlc *Button) GetBytes() []byte {
	bs, _ := json.Marshal(mlc)
	return bs
}

func NewButton(device Device, buttonName, field string, payloadPress string) *Button {
	config := &Button{
		CommonDefinition: CommonDefinition{
			CommandTopic:     fmt.Sprintf(`yzz-dali/%s/set/%s`, device.UID, field),
			Device:           device,
			EnabledByDefault: true,
			// DeviceClass:      "restart",
			Name:     buttonName,
			ObjectID: fmt.Sprintf("%s_btn_%s_%s", device.UID, field, payloadPress),
			UniqueID: fmt.Sprintf("%s_btn_%s_%s_yzz-dali", device.UID, field, payloadPress),
			Origin: &Origin{
				Name: "yzz-dali",
				Sw:   "1.0.0",
				Url:  "https://yzz-dali.dz163.cn",
			},
		},
		PayloadPress: payloadPress,
	}
	return config
}
