package config

type Captcha struct {
	KeyLong            int `mapstructure:"key-long" yaml:"key-long"`
	ImageWidth         int `mapstructure:"img-width" yaml:"img-width"`
	ImageHeight        int `mapstructure:"img-height" yaml:"img-height"`
	OpenCaptcha        int `mapstructure:"open-captcha" yaml:"open-captcha"`
	OpenCaptchaTimeOut int `mapstructure:"open-captcha-timeout" yaml:"open-captcha-timeout"`
}
