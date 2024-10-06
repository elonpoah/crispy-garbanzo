package v1

import (
	"crispy-garbanzo/common/response"
	systemReq "crispy-garbanzo/internal/app/models/request"
	"crispy-garbanzo/internal/app/service"

	"crispy-garbanzo/utils"

	"github.com/gin-gonic/gin"
)

type SessionApi struct{}

// GetHomeRecommand
// @Tags      活动中心
// @Summary   首页推荐
// @Produce  application/json
// @Success   200   {object}  response.Response{data=system.ActivitySession, msg=string}
// @Router    /api/home/recommand [post]
func (b *SessionApi) GetHomeRecommand(c *gin.Context) {
	result, errCode := service.ServiceGroupSys.GetHomeRecommand()
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(result, response.SUCCESS, c)
}

// GetSessionById
// @Tags      活动中心
// @Summary   活动场次详情
// @Produce  application/json
// @Param     data  body      systemReq.SessionDetailReq true  "场次ID"
// @Success   200   {object}  response.Response{data=system.ActivitySession, msg=string}
// @Router    /api/session/detail [post]
func (b *SessionApi) GetSessionById(c *gin.Context) {
	var req systemReq.SessionDetailReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	result, errCode := service.ServiceGroupSys.GetSessionById(req.Id)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(result, response.SUCCESS, c)
}

// BuySessionTicket
// @Tags      活动中心
// @Summary   购买入场券
// @Produce  application/json
// @Security  ApiKeyAuth
// @Param     data  body      systemReq.SessionDetailReq true  "场次ID"
// @Success   200   {object}  response.Response{data=system.ActivitySession, msg=string}
// @Router    /api/session/ticket [post]
func (b *SessionApi) BuySessionTicket(c *gin.Context) {
	var req systemReq.SessionDetailReq
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
	errCode := service.ServiceGroupSys.BuySessionTicket(req.Id, uid)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.Ok(c)
}

// CheckSession
// @Tags      活动中心
// @Summary   是否已购入场券
// @Produce  application/json
// @Security  ApiKeyAuth
// @Param     data  body      systemReq.SessionDetailReq true  "场次ID"
// @Success   200   {object}  response.Response{data=bool, msg=string}
// @Router    /api/session/check [post]
func (b *SessionApi) CheckSession(c *gin.Context) {
	var req systemReq.SessionDetailReq
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
	result, errCode := service.ServiceGroupSys.CheckSession(req.Id, uid)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(result, response.SUCCESS, c)
}

// GetSessionList
// @Tags      活动中心
// @Summary   场次列表
// @Produce   application/json
// @Param     data  body      systemReq.SessionListReq                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/session/list [post]
func (b *SessionApi) GetSessionList(c *gin.Context) {
	var req systemReq.SessionListReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	err = utils.Verify(req, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	list, total, errCode := service.ServiceGroupSys.GetSessionList(req)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, response.SUCCESS, c)
}

// GetGameHistory
// @Tags      活动中心
// @Summary   抽奖记录
// @Produce   application/json
// @Param     data  body      systemReq.GameHistoryReq                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/game/history [post]
func (b *SessionApi) GetGameHistory(c *gin.Context) {
	var req systemReq.GameHistoryReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	err = utils.Verify(req, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(response.NotfoundParameter, c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(response.InvalidUserId, c)
		return
	}
	list, total, errCode := service.ServiceGroupSys.GetGameHistory(req, uid)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, response.SUCCESS, c)
}

// GetUserSummary
// @Tags      活动中心
// @Summary   活动统计
// @Security  ApiKeyAuth
// @Produce  application/json
// @Success   200   {object}  response.Response{data=response.UserSummaryResponse, msg=string}  "返回包括用户信息"
// @Router    /api/session/summary [get]
func (b *SessionApi) GetUserSummary(c *gin.Context) {
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(response.InvalidUserId, c)
		return
	}
	result, errCode := service.ServiceGroupSys.GetUserSummary(uid)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(result, response.SUCCESS, c)
}

// CheckInviteDuty
// @Tags      活动中心
// @Summary   邀请注册活动
// @Produce  application/json
// @Security  ApiKeyAuth
// @Param     data  body      systemReq.CheckInviteDutyReq true  "场次ID"
// @Success   200   {object}  response.Response{data=bool, msg=string}
// @Router    /api/free/inviteInfo [post]
func (b *SessionApi) CheckInviteDuty(c *gin.Context) {
	var req systemReq.CheckInviteDutyReq
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
	result, errCode := service.ServiceGroupSys.CheckInviteDuty(req.Type, uid)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(result, response.SUCCESS, c)
}

// StartInviteSpin
// @Tags      活动中心
// @Summary   邀请注册活动抽奖
// @Produce  application/json
// @Security  ApiKeyAuth
// @Param     data  body      systemReq.CheckInviteDutyReq true  "场次ID"
// @Success   200   {object}  response.Response{data=bool, msg=string}
// @Router    /api/free/inviteSpin [post]
func (b *SessionApi) StartInviteSpin(c *gin.Context) {
	var req systemReq.CheckInviteDutyReq
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
	result, errCode := service.ServiceGroupSys.CheckInviteDuty(req.Type, uid)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	bonus, errCode := service.ServiceGroupSys.StartInviteSpin(req.Type, uid, result)
	if errCode != response.SUCCESS {
		response.FailWithMessage(errCode, c)
		return
	}
	response.OkWithDetailed(bonus, response.SUCCESS, c)
}
