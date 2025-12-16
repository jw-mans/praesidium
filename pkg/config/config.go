package config

import "time"

type HealthConfig struct {
	Ping string `yaml:"ping"`
	HTTP string `yaml:"http"`
}

type ActionCfg struct {
	Run string `yaml:"run,omitempty"`
	Log string `yaml:"log,omitempty"`
}

type Config struct {
	Iface         string        `yaml:"iface"`
	CheckInterval time.Duration `yaml:"check_interval"`
	IPCheckURL    string        `yaml:"ip_check_url"`

	Healthcheck  HealthConfig `yaml:"healthcheck"`
	OnDisconnect []ActionCfg  `yaml:"on_disconnect"`
}

// ApplyDefaults sets default values for Config fields if they are not set
func (c *Config) ApplyDefaults() {
	if c.Iface == "" {
		c.Iface = "wg0"
	}
	if c.CheckInterval == 0 {
		c.CheckInterval = 3 * time.Second
	}
	if c.IPCheckURL == "" {
		c.IPCheckURL = "https://api.ipify.org"
	}
}
