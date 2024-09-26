package v1

import (
	"crispy-garbanzo/common/response"
	"crispy-garbanzo/internal/app/service"

	"github.com/gin-gonic/gin"
)

type SystemApi struct{}

// SetWithdrawConfig
// @Tags      公共中心
// @Summary   获取平台配置
// @accept    application/json
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "平台配置"
// @Router    /app/platform/setting [get]
func (b *SystemApi) GetPlatformSetting(c *gin.Context) {
	info, errCode := service.ServiceGroupSys.GetPlatformSetting()
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(info, response.SUCCESS, c)
}
