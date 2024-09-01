package v1

import (
	"crispy-garbanzo/common/response"
	"crispy-garbanzo/global"
	system "crispy-garbanzo/internal/app/models"
	systemReq "crispy-garbanzo/internal/app/models/request"
	systemRes "crispy-garbanzo/internal/app/models/response"
	"crispy-garbanzo/internal/app/service"
	"crispy-garbanzo/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysUserApi struct{}

// Login
// @Tags     用户中心
// @Summary  登录
// @Produce   application/json
// @Param    data  body      systemReq.Login                                             true  "用户名, 密码"
// @Success  200   {object}  response.Response{data=systemRes.LoginResponse,msg=string}  "返回包括用户信息,token"
// @Router   /api/login [post]
func (b *SysUserApi) Login(c *gin.Context) {
	var l systemReq.Login
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
	userInfo := systemRes.SysUserResponse{
		ID:       user.ID,
		Username: user.Username,
		NickName: user.NickName,
		Phone:    user.Phone,
		Email:    user.Email,
	}
	response.OkWithDetailed(systemRes.LoginResponse{
		User:  userInfo,
		Token: token,
	}, "登录成功", c)
}

// Register
// @Tags     用户中心
// @Summary  注册账号
// @Produce   application/json
// @Param    data  body      systemReq.Register                                            true  "用户名, 昵称, 密码,"
// @Success  200   {object}  response.Response{data=systemRes.LoginResponse,msg=string}  "用户注册账号,返回包括用户信息 token"
// @Router   /api/register [post]
func (b *SysUserApi) Register(c *gin.Context) {
	var r systemReq.Register
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user := &system.SysUser{Username: r.Username, NickName: r.NickName, Password: r.Password, Phone: r.Phone, Email: r.Email}
	userReturn, err := service.ServiceGroupSys.Register(*user)
	if err != nil {
		global.FPG_LOG.Error("注册失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	userInfo := systemRes.SysUserResponse{
		ID:       userReturn.ID,
		Username: userReturn.Username,
		NickName: userReturn.NickName,
		Phone:    userReturn.Phone,
		Email:    userReturn.Email,
	}
	token := utils.CreateToken(user.Username, "app", global.FPG_CONFIG.Jwt.Key, global.FPG_CONFIG.Jwt.ExpireTime)
	response.OkWithDetailed(systemRes.LoginResponse{
		User:  userInfo,
		Token: token,
	}, "注册成功", c)
}

// ChangePassword
// @Tags      用户中心
// @Summary   修改密码
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body      systemReq.ChangePasswordReq    true  "原密码, 新密码"
// @Success   200   {object}  response.Response{msg=string}  "用户修改密码"
// @Router    /api/changePassword [post]
func (b *SysUserApi) ChangePassword(c *gin.Context) {
	var req systemReq.ChangePasswordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.ChangePasswordVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	u := &system.SysUser{Username: uid, Password: req.Password}
	_, err = service.ServiceGroupSys.ChangePassword(u, req.NewPassword)
	if err != nil {
		global.FPG_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("修改成功", c)
}
