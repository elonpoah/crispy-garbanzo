package v1

import (
	"crispy-garbanzo/common/response"
	"crispy-garbanzo/global"
	system "crispy-garbanzo/internal/admin/models"
	"crispy-garbanzo/internal/admin/models/request"
	systemRes "crispy-garbanzo/internal/admin/models/response"
	"crispy-garbanzo/internal/admin/service"
	"crispy-garbanzo/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysUserApi struct{}

// Login
// @Tags     Base
// @Summary  用户登录
// @Produce   application/json
// @Param    data  body      systemReq.Login                                             true  "用户名, 密码, 验证码"
// @Success  200   {object}  response.Response{data=systemRes.LoginResponse,msg=string}  "返回包括用户信息,token,过期时间"
// @Router   /base/login [post]
func (b *SysUserApi) Login(c *gin.Context) {
	var l request.Login
	err := c.ShouldBindJSON(&l)

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	u := &system.SysUser{Username: l.Username, Password: l.Password}
	user, err := service.ServiceGroupSys.Login(u)
	if err != nil {
		global.FPG_LOG.Error("登陆失败! 用户名不存在或者密码错误!", zap.Error(err))
		response.FailWithMessage("用户名不存在或者密码错误", c)
		return
	}
	if user.Enable != 1 {
		global.FPG_LOG.Error("登陆失败! 用户被禁止登录!")
		response.FailWithMessage("用户被禁止登录", c)
		return
	}
	token := utils.CreateToken(user.Username, "app", global.FPG_CONFIG.Jwt.Key, global.FPG_CONFIG.Jwt.ExpireTime)
	response.OkWithDetailed(systemRes.LoginResponse{
		User:  *user,
		Token: token,
	}, "登录成功", c)
}
