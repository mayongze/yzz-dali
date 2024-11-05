package homeassistant

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type SysConfig struct {
	Mqtt    MqttConfig              `yaml:"mqtt"`
	Serial  string                  `yaml:"serial"`
	Devices map[string]DeviceConfig `yaml:"devices"`
	Groups  map[string]GroupConfig  `yaml:"groups"`
}

type MqttConfig struct {
	Server    string `yaml:"server"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	BaseTopic string `yaml:"baseTopic"`
}

type DeviceConfig struct {
	Type       DeviceType `yaml:"type"`
	Addr       int        `yaml:"addr,omitempty"`
	Name       string     `yaml:"name,omitempty"`
	Transition int        `yaml:"transition,omitempty"`
}

type GroupConfig struct {
	DeviceConfig `yaml:",inline"`
	Devices      []string `yaml:"devices"`
}

var Cfg *SysConfig
var CfgFile string

func InitSysConfig(file string) *SysConfig {
	CfgFile = file
	Cfg = &SysConfig{}
	bs, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	if err = yaml.Unmarshal(bs, Cfg); err != nil {
		log.Fatal(err)
	}
	return Cfg
}

func SaveSysConfig() {
	file := CfgFile
	if len(file) == 0 {
		log.Fatal("config file not set")
	}
	bs, err := yaml.Marshal(Cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err = os.WriteFile(file, bs, 0644); err != nil {
		log.Fatal(err)
	}
}
