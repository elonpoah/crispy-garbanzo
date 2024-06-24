package config

type Application struct {
	Mode          string `mapstructure:"mode" yaml:"mode"`
	Name          string `mapstructure:"name" yaml:"name"`
	Port          string `mapstructure:"port" yaml:"port"`
	Readtimeout   int    `mapstructure:"readtimeout" yaml:"readtimeout"`
	Writertimeout int    `mapstructure:"writertimeout" yaml:"writertimeout"`
}
