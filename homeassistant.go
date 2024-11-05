package homeassistant

import "encoding/json"

// MQTTDevice represents the device information
type MQTTDevice struct {
	ConfigurationURL string     `json:"configuration_url,omitempty"`
	Connections      [][]string `json:"connections,omitempty"`
	Identifiers      []string   `json:"identifiers,omitempty"`
	Manufacturer     string     `json:"manufacturer,omitempty"`
	Model            string     `json:"model,omitempty"`
	Name             string     `json:"name,omitempty"`
	SerialNumber     string     `json:"serial_number,omitempty"`
	SoftwareVersion  string     `json:"sw_version,omitempty"`
	// via_device
}

// Availability represents the structure for MQTT availability
type Availability struct {
	Topic               string `json:"topic"`                           // Required, MQTT topic to receive availability updates
	PayloadAvailable    string `json:"payload_available,omitempty"`     // Optional, payload to represent the available state
	PayloadNotAvailable string `json:"payload_not_available,omitempty"` // Optional, payload to represent the unavailable state
	ValueTemplate       string `json:"value_template,omitempty"`        // Optional, template to extract availability from the topic
}

type MQTTNumberConfig struct {
	Availability           []Availability `yaml:"availability,omitempty"`             // Optional
	AvailabilityTopic      string         `yaml:"availability_topic,omitempty"`       // Optional
	AvailabilityMode       string         `yaml:"availability_mode,omitempty"`        // Optional, default: latest
	CommandTemplate        string         `yaml:"command_template,omitempty"`         // Optional
	CommandTopic           string         `yaml:"command_topic"`                      // Required
	Device                 MQTTDevice     `yaml:"device,omitempty"`                   // Optional
	DeviceClass            string         `yaml:"device_class,omitempty"`             // Optional
	EnabledByDefault       bool           `yaml:"enabled_by_default,omitempty"`       // Optional, default: true
	Encoding               string         `yaml:"encoding,omitempty"`                 // Optional, default: utf-8
	EntityCategory         string         `yaml:"entity_category,omitempty"`          // Optional
	Icon                   string         `yaml:"icon,omitempty"`                     // Optional
	JSONAttributesTemplate string         `yaml:"json_attributes_template,omitempty"` // Optional
	JSONAttributesTopic    string         `yaml:"json_attributes_topic,omitempty"`    // Optional
	Min                    float64        `yaml:"min,omitempty"`                      // Optional, default: 1
	Max                    float64        `yaml:"max,omitempty"`                      // Optional, default: 100
	Mode                   string         `yaml:"mode,omitempty"`                     // Optional, default: "auto"
	Name                   string         `yaml:"name,omitempty"`                     // Optional
	ObjectID               string         `yaml:"object_id,omitempty"`                // Optional
	Optimistic             bool           `yaml:"optimistic,omitempty"`               // Optional, default: true if no state_topic, else false
	PayloadReset           string         `yaml:"payload_reset,omitempty"`            // Optional, default: "None"
	QoS                    int            `yaml:"qos,omitempty"`                      // Optional, default: 0
	Retain                 bool           `yaml:"retain,omitempty"`                   // Optional, default: false
	StateTopic             string         `yaml:"state_topic,omitempty"`              // Optional
	Step                   float64        `yaml:"step,omitempty"`                     // Optional, default: 1
	UniqueID               string         `yaml:"unique_id,omitempty"`                // Optional
	UnitOfMeasurement      string         `yaml:"unit_of_measurement,omitempty"`      // Optional
	ValueTemplate          string         `yaml:"value_template,omitempty"`           // Optional
}

func (mnc *MQTTNumberConfig) GetBytes() []byte {
	bs, _ := json.Marshal(mnc)
	return bs
}

type MQTTSelectConfig struct {
	Availability           []Availability `json:"availability,omitempty"`             // Optional, array of availability structures
	AvailabilityTopic      string         `json:"availability_topic,omitempty"`       // Optional
	AvailabilityMode       string         `json:"availability_mode,omitempty"`        // Optional (default: "latest")
	AvailabilityTemplate   string         `json:"availability_template,omitempty"`    // Optional
	CommandTemplate        string         `json:"command_template,omitempty"`         // Optional
	CommandTopic           string         `json:"command_topic"`                      // Required
	Device                 MQTTDevice     `json:"device,omitempty"`                   // Optional
	EnabledByDefault       bool           `json:"enabled_by_default,omitempty"`       // Optional (default: true)
	Encoding               string         `json:"encoding,omitempty"`                 // Optional (default: "utf-8")
	EntityCategory         string         `json:"entity_category,omitempty"`          // Optional
	Icon                   string         `json:"icon,omitempty"`                     // Optional
	JSONAttributesTemplate string         `json:"json_attributes_template,omitempty"` // Optional
	JSONAttributesTopic    string         `json:"json_attributes_topic,omitempty"`    // Optional
	Name                   string         `json:"name,omitempty"`                     // Optional, default: "MQTT JSON Light"
	ObjectID               string         `json:"object_id,omitempty"`                // Optional
	Optimistic             bool           `json:"optimistic,omitempty"`               // Optional, default: true if no state_topic defined, else false
	Options                []string       `json:"options"`                            // Required
	QoS                    int            `json:"qos,omitempty"`                      // Optional (default: 0)
	Retain                 bool           `json:"retain,omitempty"`                   // Optional (default: false)
	StateTopic             string         `json:"state_topic,omitempty"`              // Optional
	UniqueID               string         `json:"unique_id,omitempty"`                // Optional
	ValueTemplate          string         `json:"value_template,omitempty"`           // Optional
}

func (mlc *MQTTSelectConfig) GetBytes() []byte {
	bs, _ := json.Marshal(mlc)
	return bs
}

// MQTTLightConfig represents the configuration for an MQTT JSON light
type MQTTLightConfig struct {
	Availability           []Availability `json:"availability,omitempty"`             // Optional, array of availability structures
	Schema                 string         `json:"schema,omitempty"`                   // Must be set to "json"
	CommandTopic           string         `json:"command_topic"`                      // Required
	StateTopic             string         `json:"state_topic,omitempty"`              // Optional
	AvailabilityMode       string         `json:"availability_mode,omitempty"`        // Optional
	AvailabilityTemplate   string         `json:"availability_template,omitempty"`    // Optional
	Brightness             bool           `json:"brightness,omitempty"`               // Optional, default: false
	BrightnessScale        int            `json:"brightness_scale,omitempty"`         // Optional, default: 255
	Device                 MQTTDevice     `json:"device,omitempty"`                   // Optional
	EnabledByDefault       bool           `json:"enabled_by_default,omitempty"`       // Optional, default: true
	Encoding               string         `json:"encoding,omitempty"`                 // Optional, default: utf-8
	EntityCategory         string         `json:"entity_category,omitempty"`          // Optional
	Effect                 bool           `json:"effect,omitempty"`                   // Optional, default: false
	EffectList             []string       `json:"effect_list,omitempty"`              // Optional
	FlashTimeLong          int            `json:"flash_time_long,omitempty"`          // Optional, default: 10
	FlashTimeShort         int            `json:"flash_time_short,omitempty"`         // Optional, default: 2
	Icon                   string         `json:"icon,omitempty"`                     // Optional
	JsonAttributesTemplate string         `json:"json_attributes_template,omitempty"` // Optional
	JsonAttributesTopic    string         `json:"json_attributes_topic,omitempty"`    // Optional
	MaxMireds              int            `json:"max_mireds,omitempty"`               // Optional
	MinMireds              int            `json:"min_mireds,omitempty"`               // Optional
	Name                   string         `json:"name,omitempty"`                     // Optional, default: "MQTT JSON Light"
	ObjectID               string         `json:"object_id,omitempty"`                // Optional
	Optimistic             bool           `json:"optimistic,omitempty"`               // Optional, default: true if no state_topic defined, else false
	PayloadAvailable       string         `json:"payload_available,omitempty"`        // Optional, default: "online"
	PayloadNotAvailable    string         `json:"payload_not_available,omitempty"`    // Optional, default: "offline"
	QoS                    int            `json:"qos,omitempty"`                      // Optional, default: 0
	Retain                 bool           `json:"retain,omitempty"`                   // Optional, default: false
	SupportedColorModes    []string       `json:"supported_color_modes,omitempty"`    // Optional
	UniqueID               string         `json:"unique_id,omitempty"`                // Optional
	WhiteScale             int            `json:"white_scale,omitempty"`              // Optional, default: 255
}

func (mlc *MQTTLightConfig) GetBytes() []byte {
	bs, _ := json.Marshal(mlc)
	return bs
}
