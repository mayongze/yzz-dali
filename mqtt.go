package homeassistant

import (
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"strings"
	"time"
)

var (
	mqttCli mqtt.Client
	// 注册订阅函数,方便断线重连
	onConnectList []func(client mqtt.Client) = make([]func(client mqtt.Client), 0)
)

func GetMqttCli() mqtt.Client {
	return mqttCli
}

func RegisterOnConnect(f func(client mqtt.Client)) {
	onConnectList = append(onConnectList, f)
}

func InitMQTTClient(config MqttConfig) mqtt.Client {
	clientID := "go_dali_client"
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Server).SetClientID(clientID).SetUsername(config.User).SetPassword(config.Password).
		SetKeepAlive(60 * time.Second).SetPingTimeout(time.Second).SetAutoReconnect(true).
		SetConnectRetry(true).SetConnectRetryInterval(3 * time.Second)
	opts.SetWill("yzz-dali/bridge/status", `{"state": "offline"}`, 1, true)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Print("Connected to MQTT server")
		// 发布Birth消息
		if token := client.Publish("yzz-dali/bridge/status", 1, true, `{"state": "online"}`); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
		for _, onConnect := range onConnectList {
			onConnect(client)
		}
	})
	mqttCli = mqtt.NewClient(opts)
	if token := mqttCli.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	return mqttCli
}

func NewNumberDiscoveryConfig(deviceUID, deviceName string, numberName, field string, min, max, step float64, unit string) *MQTTNumberConfig {
	config := &MQTTNumberConfig{
		Availability: []Availability{
			{
				Topic:         "yzz-dali/bridge/status",
				ValueTemplate: "{{ value_json.state }}",
			},
		},
		CommandTopic: fmt.Sprintf(`yzz-dali/%s/set/%s`, deviceUID, field),
		Device: MQTTDevice{
			Identifiers:     []string{fmt.Sprintf("yzz-dali-%s", deviceUID)},
			Manufacturer:    "yzz-dali",
			Model:           "Light ZS",
			Name:            deviceName,
			SoftwareVersion: "yzz 1.0.0",
		},
		EnabledByDefault:  true,
		EntityCategory:    "config",
		Icon:              "mdi:tune",
		Name:              numberName,
		Min:               min,
		Max:               max,
		UnitOfMeasurement: unit,
		Step:              step,
		ObjectID:          fmt.Sprintf("%s_number_%s", deviceUID, field),
		StateTopic:        fmt.Sprintf("yzz-dali/%s", deviceUID),
		UniqueID:          fmt.Sprintf("%s_number_%s_yzz-dali", deviceUID, field),
		ValueTemplate:     fmt.Sprintf("{{ value_json.%s }}", field),
	}
	return config
}

func SubscriptionCommandSet(client mqtt.Client) {
	// 订阅设备控制命令
	client.Subscribe("yzz-dali/+/set", 0, func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := msg.Payload()
		stateKV := StateKV{}
		_ = json.Unmarshal(payload, &stateKV)
		log.Print("Received message from topic: ", topic)
		log.Print("Message: ", string(payload))
		log.Print("State: ", stateKV)

		uid := strings.Split(topic, "/")[1]
		device := GlobalDeviceList[uid]
		err := device.HandleCommand(stateKV)
		if err != nil {
			log.Print("HandleCommand error: ", err)
		}
	})

	client.Subscribe("yzz-dali/+/set/+", 0, func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := msg.Payload()
		log.Print("Received message from topic: ", topic)
		log.Print("Message: ", string(payload))

		uid := strings.Split(topic, "/")[1]
		filed := strings.Split(topic, "/")[3]
		stateKV := StateKV{
			filed: string(payload),
		}
		if strings.HasPrefix(uid, "strdance") {
			device := GlobalStrDanceDeviceList[uid]
			err := device.HandleCommand(stateKV)
			if err != nil {
				log.Print("HandleCommand error: ", err)
			}
		} else {
			device := GlobalDeviceList[uid]
			err := device.HandleCommand(stateKV)
			if err != nil {
				log.Print("HandleCommand error: ", err)
			}
		}
	})
}
