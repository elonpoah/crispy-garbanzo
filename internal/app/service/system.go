package service

import (
	"crispy-garbanzo/common/response"
	"crispy-garbanzo/global"
	system "crispy-garbanzo/internal/app/models"
	"encoding/json"
)

type SystemService struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetPlatformSetting
//@description: 获取平台配置
//@return: errCode int

func (userService *UserService) GetPlatformSetting() (setting map[string]interface{}, errCode int) {
	var u system.WebStting
	err := global.FPG_DB.Where("key_name", system.AppWithdrawSetting).First(&u).Error
	if err != nil {
		return nil, response.InternalServerError
	}
	var finance map[string]interface{}

	err = json.Unmarshal(u.Info, &finance)
	if err != nil {
		return nil, response.InternalServerError
	}

	setting = make(map[string]interface{})

	setting["finance"] = finance

	return setting, response.SUCCESS
}
