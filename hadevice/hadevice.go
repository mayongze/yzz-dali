package hadevice

// Device represents the device information
type Device struct {
	ConfigurationURL string `json:"configuration_url,omitempty"`
	// A list of connections of the device to the outside world as a list of tuples [connection_type, connection_identifier].
	// For example the MAC address of a network interface: "connections": [["mac", "02:5b:26:a8:dc:12"]].
	Connections     [][]string `json:"connections,omitempty"`
	Identifiers     []string   `json:"identifiers,omitempty"`
	UID             string     `json:"-"`
	Manufacturer    string     `json:"manufacturer,omitempty"`
	Model           string     `json:"model,omitempty"`
	ModelId         string     `json:"model_id,omitempty"`
	Name            string     `json:"name,omitempty"`
	SerialNumber    string     `json:"serial_number,omitempty"`
	SoftwareVersion string     `json:"sw_version,omitempty"`
	HWVersion       string     `json:"hw_version,omitempty"`
	SuggestedArea   string     `json:"suggested_area,omitempty"`
	// Identifier of a device that routes messages between this device and Home Assistant.
	// Examples of such devices are hubs, or parent devices of a sub-device. This is used to show device topology in Home Assistant.
	ViaDevice string `json:"via_device,omitempty"`
}

// Availability represents the structure for MQTT availability
type Availability struct {
	Topic               string `json:"topic"`                           // Required, MQTT topic to receive availability updates
	PayloadAvailable    string `json:"payload_available,omitempty"`     // Optional, payload to represent the available state
	PayloadNotAvailable string `json:"payload_not_available,omitempty"` // Optional, payload to represent the unavailable state
	ValueTemplate       string `json:"value_template,omitempty"`        // Optional, template to extract availability from the topic
}

type Origin struct {
	Name string `json:"name"`
	Sw   string `json:"sw"`
	Url  string `json:"url"`
}

type CommonDefinition struct {
	Name                 string         `json:"name,omitempty"`      // Optional
	UniqueID             string         `json:"unique_id,omitempty"` // Optional
	ObjectID             string         `json:"object_id,omitempty"` // Optional
	Availability         []Availability `json:"availability,omitempty"`
	AvailabilityTopic    string         `json:"availability_topic,omitempty"`
	AvailabilityMode     string         `json:"availability_mode,omitempty"`
	AvailabilityTemplate string         `json:"availability_template,omitempty"` // Optional
	Device               Device         `json:"device,omitempty"`                // Optional
	DeviceClass          string         `json:"device_class,omitempty"`          // Optional
	EnabledByDefault     bool           `json:"enabled_by_default,omitempty"`    // Optional, default: true
	Encoding             string         `json:"encoding,omitempty"`              // Optional, default: utf-8
	// diagnostic config
	EntityCategory         string  `json:"entity_category,omitempty"`          // Optional
	Icon                   string  `json:"icon,omitempty"`                     // Optional
	JsonAttributesTemplate string  `json:"json_attributes_template,omitempty"` // Optional
	JsonAttributesTopic    string  `json:"json_attributes_topic,omitempty"`    // Optional
	QoS                    int     `json:"qos,omitempty"`                      // Optional (default: 0)
	Retain                 bool    `json:"retain,omitempty"`                   // Optional (default: false)
	Optimistic             bool    `json:"optimistic,omitempty"`               // Optional, default: true if no state_topic defined, else false
	StateTopic             string  `json:"state_topic,omitempty"`              // Optional
	ValueTemplate          string  `json:"value_template,omitempty"`           // Optional
	CommandTopic           string  `json:"command_topic,omitempty"`            // Required
	CommandTemplate        string  `json:"command_template,omitempty"`         // Optional
	PayloadAvailable       string  `json:"payload_available,omitempty"`        // Optional, default: "online"
	PayloadNotAvailable    string  `json:"payload_not_available,omitempty"`    // Optional, default: "offline"
	PayloadReset           string  `json:"payload_reset,omitempty"`            // Optional, default: "None"
	Origin                 *Origin `json:"origin,omitempty"`
}
