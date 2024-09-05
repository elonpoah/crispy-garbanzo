package v1

import (
	"crispy-garbanzo/common/request"
	"crispy-garbanzo/common/response"
	"crispy-garbanzo/global"
	systemReq "crispy-garbanzo/internal/app/models/request"
	"crispy-garbanzo/internal/app/service"

	"crispy-garbanzo/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SessionApi struct{}

// GetHomeRecommand
// @Tags      活动中心
// @Summary   首页推荐
// @Produce  application/json
// @Success   200   {object}  response.Response{data=system.ActivitySession, msg=string}
// @Router    /api/home/recommand [post]
func (b *SessionApi) GetHomeRecommand(c *gin.Context) {
	result, err := service.ServiceGroupSys.GetHomeRecommand()
	if err != nil {
		global.FPG_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
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
		response.FailWithMessage(err.Error(), c)
		return
	}
	result, err := service.ServiceGroupSys.GetSessionById(req.Id)
	if err != nil {
		global.FPG_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "获取成功", c)
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
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = service.ServiceGroupSys.BuySessionTicket(req.Id, uid)
	if err != nil {
		global.FPG_LOG.Error("购买失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("成功", c)
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
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	result, err := service.ServiceGroupSys.CheckSession(req.Id, uid)
	if err != nil {
		global.FPG_LOG.Error("查询场次失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "成功", c)
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
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := service.ServiceGroupSys.GetSessionList(req)
	if err != nil {
		global.FPG_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// GetGameHistory
// @Tags      活动中心
// @Summary   抽奖记录
// @Produce   application/json
// @Param     data  body      request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/game/history [post]
func (b *SessionApi) GetGameHistory(c *gin.Context) {
	var req request.PageInfo
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid, err := utils.GetUserID(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := service.ServiceGroupSys.GetGameHistory(req, uid)
	if err != nil {
		global.FPG_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}
