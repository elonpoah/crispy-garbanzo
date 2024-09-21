package initialize

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func I18n() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFile("message/en-US.toml")
	bundle.LoadMessageFile("message/id-ID.toml")
	bundle.LoadMessageFile("message/ja-JP.toml")
	bundle.LoadMessageFile("message/ko-KR.toml")
	bundle.LoadMessageFile("message/pt-BR.toml")
	bundle.LoadMessageFile("message/th-TH.toml")
	bundle.LoadMessageFile("message/vi-VN.toml")
	bundle.LoadMessageFile("message/zh-CN.toml")
	return bundle
}
