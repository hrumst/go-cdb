package config

import "time"

type AppConfig struct {
	Engine  AppConfigEngine  `yaml:"engine"`
	Network AppConfigNetwork `yaml:"network"`
	Logging AppConfigLogging `yaml:"logging"`
}

type AppConfigEngine struct {
	Type string `yaml:"type"`
}

type AppConfigNetwork struct {
	Address        string        `yaml:"address"`
	MaxConnections int64         `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	AcceptTimeout  time.Duration `yaml:"accept_timeout"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

type AppConfigLogging struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}
