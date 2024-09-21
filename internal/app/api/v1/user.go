package v1

import (
	"crispy-garbanzo/common/request"
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
	token := utils.CreateToken(user.ID, "app", global.FPG_CONFIG.Jwt.Key, global.FPG_CONFIG.Jwt.ExpireTime)
	userInfo := systemRes.SysUserResponse{
		ID:            user.ID,
		Pid:           user.Pid,
		Username:      user.Username,
		NickName:      user.NickName,
		Phone:         user.Phone,
		Email:         user.Email,
		Balance:       user.Balance,
		FreezeBalance: user.FreezeBalance,
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
	user := &system.SysUser{Username: r.Username, NickName: r.NickName, Password: r.Password, Phone: r.Phone, Email: r.Email, Pid: r.Pid}
	userReturn, err := service.ServiceGroupSys.Register(*user)
	if err != nil {
		global.FPG_LOG.Error("注册失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	userInfo := systemRes.SysUserResponse{
		ID:            userReturn.ID,
		Pid:           user.Pid,
		Username:      userReturn.Username,
		NickName:      userReturn.NickName,
		Phone:         userReturn.Phone,
		Email:         userReturn.Email,
		Balance:       userReturn.Balance,
		FreezeBalance: userReturn.FreezeBalance,
	}
	token := utils.CreateToken(userReturn.ID, "app", global.FPG_CONFIG.Jwt.Key, global.FPG_CONFIG.Jwt.ExpireTime)
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
// @Router    /api/change/password [post]
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
	err = service.ServiceGroupSys.ChangePassword(uid, req.Password, req.NewPassword)
	if err != nil {
		global.FPG_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

// GetUserInfo
// @Tags      用户中心
// @Summary   用户信息
// @Security  ApiKeyAuth
// @Produce  application/json
// @Success   200   {object}  response.Response{data=systemRes.LoginResponse, msg=string}  "返回包括用户信息"
// @Router    /api/user/info [get]
func (b *SysUserApi) GetUserInfo(c *gin.Context) {
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userInfo, err := service.ServiceGroupSys.GetUserInfo(uid)
	if err != nil {
		global.FPG_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	result := systemRes.SysUserResponse{
		ID:            userInfo.ID,
		Username:      userInfo.Username,
		NickName:      userInfo.NickName,
		Phone:         userInfo.Phone,
		Email:         userInfo.Email,
		Balance:       userInfo.Balance,
		FreezeBalance: userInfo.FreezeBalance,
	}
	response.OkWithDetailed(result, "获取成功", c)
}

// GetUserDepositList
// @Tags      用户中心
// @Summary   充值记录
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  query      systemReq.UserDepositRecordReq                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/deposit/history [get]
func (b *SysUserApi) GetUserDepositList(c *gin.Context) {
	var pageInfo systemReq.UserDepositRecordReq
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := service.ServiceGroupSys.GetUserDepositList(pageInfo, uid)
	if err != nil {
		global.FPG_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// Deposit
// @Tags      用户中心
// @Summary   充值
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body     systemReq.UserDepositReq                                        true  "参数"
// @Success   200   {object}  response.Response{data=interface{},msg=string}  "获取地址"
// @Router    /api/user/deposit [post]
func (b *SysUserApi) Deposit(c *gin.Context) {
	var req systemReq.UserDepositReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.Uid = uid
	result, err := service.ServiceGroupSys.Deposit(req)
	if err != nil {
		global.FPG_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
}

// Withdraw
// @Tags      用户中心
// @Summary   提现
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body     systemReq.UserWithdrawReq                                        true  "参数"
// @Success   200   {object}  response.Response{data=interface{},msg=string}  "获取地址"
// @Router    /api/user/withdraw [post]
func (b *SysUserApi) Withdraw(c *gin.Context) {
	var req systemReq.UserWithdrawReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.Uid = uid
	err = service.ServiceGroupSys.Withdraw(req)
	if err != nil {
		global.FPG_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("获取成功", c)
}

// GetUserWithdrawList
// @Tags      用户中心
// @Summary   提现记录
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  query      request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/withdraw/history [get]
func (b *SysUserApi) GetUserWithdrawList(c *gin.Context) {
	var pageInfo systemReq.UserWithdrawRecordReq
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := service.ServiceGroupSys.GetUserWithdrawList(pageInfo, uid)
	if err != nil {
		global.FPG_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetUserFreeSpinList
// @Tags      用户中心
// @Summary   免费旋转记录
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  query      request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/freeSpin/history [get]
func (b *SysUserApi) GetUserFreeSpinList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := service.ServiceGroupSys.GetUserFreeSpinList(pageInfo, uid)
	if err != nil {
		global.FPG_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
