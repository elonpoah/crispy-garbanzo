package config

type Server struct {
	Jwt         Jwt         `yaml:"jwt"`
	Mysql       Mysql       `yaml:"mysql"`
	Redis       Redis       `yaml:"redis"`
	Zap         Zap         `yaml:"zap"`
	Application Application `yaml:"application"`
	Captcha     Captcha     `yaml:"captcha"`
}
