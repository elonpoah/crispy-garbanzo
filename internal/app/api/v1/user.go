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
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	u := &system.SysUser{Username: l.Username, Password: l.Password}
	user, errCode := service.ServiceGroupSys.Login(u)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	if user.Enable != 1 {
		response.FailWithMessage(response.UserLoginForbiden, c)
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
	}, response.SUCCESS, c)
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
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	user := &system.SysUser{Username: r.Username, NickName: r.NickName, Password: r.Password, Phone: r.Phone, Email: r.Email, Pid: r.Pid}
	userReturn, errCode := service.ServiceGroupSys.Register(*user)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
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
	}, response.SUCCESS, c)
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
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	err = utils.Verify(req, utils.ChangePasswordVerify)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(response.InvalidUserId, c)
		return
	}
	errCode := service.ServiceGroupSys.ChangePassword(uid, req.Password, req.NewPassword)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.Ok(c)
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
		response.FailWithMessage(response.InvalidUserId, c)
		return
	}
	userInfo, errCode := service.ServiceGroupSys.GetUserInfo(uid)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
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
	response.OkWithDetailed(result, response.SUCCESS, c)
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
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(response.InvalidUserId, c)
		return
	}
	list, total, errCode := service.ServiceGroupSys.GetUserDepositList(pageInfo, uid)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, response.SUCCESS, c)
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
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(response.InvalidUserId, c)
		return
	}
	req.Uid = uid
	result, errCode := service.ServiceGroupSys.Deposit(req)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(result, response.SUCCESS, c)
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
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(response.InvalidUserId, c)
		return
	}
	req.Uid = uid
	errCode := service.ServiceGroupSys.Withdraw(req)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.Ok(c)
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
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(response.InvalidUserId, c)
		return
	}
	list, total, errCode := service.ServiceGroupSys.GetUserWithdrawList(pageInfo, uid)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, response.SUCCESS, c)
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
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(response.InvalidUserId, c)
		return
	}
	list, total, errCode := service.ServiceGroupSys.GetUserFreeSpinList(pageInfo, uid)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, response.SUCCESS, c)
}

// MakeDraw
// @Tags      用户中心
// @Summary   创建抽奖
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body     systemReq.UserMakeDrawReq       true  "参数"
// @Success   200   {object}  response.Response{data=interface{},msg=string}  "无"
// @Router    /api/draw/make [post]
func (b *SysUserApi) MakeDraw(c *gin.Context) {
	var req systemReq.UserMakeDrawReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(response.InvalidUserId, c)
		return
	}
	if req.BonusType == 2 && req.Bonus < 0.01 {
		response.FailWithMessage(response.LessThanMinDrawBonus, c)
		return
	}
	if req.BonusType == 1 && req.Bonus/float64(req.Count) < 0.01 {
		response.FailWithMessage(response.LessThanMinDrawBonus, c)
		return
	}
	req.Uid = uid
	key, errCode := service.ServiceGroupSys.MakeDraw(req)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(key, response.SUCCESS, c)
}

// GetDrawById
// @Tags      用户中心
// @Summary   用户抽奖详情
// @Produce  application/json
// @Param     data  body      systemReq.DrawDetailReq true  "场次ID"
// @Success   200   {object}  response.Response{data=system.MemberDraw, msg=string}
// @Router    /api/draw/detail [post]
func (b *SysUserApi) GetDrawById(c *gin.Context) {
	var req systemReq.DrawDetailReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	info, errCode := service.ServiceGroupSys.GetDrawById(req.Id)
	result := systemRes.DrawDetailRes{
		Username:     info.Username,
		DrawId:       info.DrawId,
		BonusType:    info.BonusType,
		Bonus:        info.Bonus,
		Distribute:   info.Distribute,
		Count:        info.Count,
		Participants: info.Participants,
		Status:       info.Status,
	}
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(result, response.SUCCESS, c)
}

// ClaimDrawById
// @Tags      用户中心
// @Summary   用户参与抽奖
// @Produce  application/json
// @Security  ApiKeyAuth
// @Param     data  body      systemReq.DrawDetailReq true  "场次ID"
// @Success   200   {object}  response.Response{data=interface{}, msg=string}
// @Router    /api/draw/join [post]
func (b *SysUserApi) ClaimDrawById(c *gin.Context) {
	var req systemReq.DrawDetailReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(response.InvalidUserId, c)
		return
	}
	bonus, errCode := service.ServiceGroupSys.ClaimDrawById(req.Id, uid)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(bonus, response.SUCCESS, c)
}

// GetUserDrawList
// @Tags      用户中心
// @Summary   我发起的抽奖
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  query      request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/draw/history [get]
func (b *SysUserApi) GetUserDrawList(c *gin.Context) {
	var pageInfo systemReq.UserMakeDrawRecordReq
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(response.InvalidUserId, c)
		return
	}
	list, total, errCode := service.ServiceGroupSys.GetUserDrawList(pageInfo, uid)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, response.SUCCESS, c)
}
