package service

import (
	"context"
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
	ctx := context.Background()
	setting = make(map[string]interface{})
	results, err := global.FPG_REIDS.MGet(ctx, system.AppWithdrawSetting, system.AppInviteSetting).Result()
	if err != nil {
		return nil, response.InternalServerError
	}

	// 分别检查 finance 和 invite 的数据
	var finance map[string]interface{}
	var invite map[string]interface{}
	// 处理提现设置
	if results[0] != nil {
		// Redis 中存在，反序列化
		err := json.Unmarshal([]byte(results[0].(string)), &finance)
		if err != nil {
			return nil, response.InternalServerError
		}
		setting["finance"] = finance
	}

	// 处理邀请设置
	if results[1] != nil {
		// Redis 中存在，反序列化
		err := json.Unmarshal([]byte(results[1].(string)), &invite)
		if err != nil {
			return nil, response.InternalServerError
		}
		setting["invite"] = invite
	}
	return setting, response.SUCCESS
}
