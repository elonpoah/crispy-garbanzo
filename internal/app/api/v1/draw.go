package v1

import (
	"crispy-garbanzo/common/response"
	systemReq "crispy-garbanzo/internal/app/models/request"
	systemRes "crispy-garbanzo/internal/app/models/response"
	"crispy-garbanzo/internal/app/service"
	"crispy-garbanzo/utils"

	"github.com/gin-gonic/gin"
)

type DrawApi struct{}

// MakeDraw
// @Tags      抽奖中心
// @Summary   创建抽奖
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body     systemReq.UserMakeDrawReq       true  "参数"
// @Success   200   {object}  response.Response{data=interface{},msg=string}  "无"
// @Router    /api/draw/make [post]
func (b *DrawApi) MakeDraw(c *gin.Context) {
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
// @Tags      抽奖中心
// @Summary   用户抽奖详情
// @Produce  application/json
// @Param     data  body      systemReq.DrawDetailReq true  "场次ID"
// @Success   200   {object}  response.Response{data=system.MemberDraw, msg=string}
// @Router    /api/draw/detail [post]
func (b *DrawApi) GetDrawById(c *gin.Context) {
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
// @Tags      抽奖中心
// @Summary   用户参与抽奖
// @Produce  application/json
// @Security  ApiKeyAuth
// @Param     data  body      systemReq.DrawDetailReq true  "场次ID"
// @Success   200   {object}  response.Response{data=interface{}, msg=string}
// @Router    /api/draw/join [post]
func (b *DrawApi) ClaimDrawById(c *gin.Context) {
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
// @Tags      抽奖中心
// @Summary   我发起的抽奖
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  query      request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/draw/history [get]
func (b *DrawApi) GetUserDrawList(c *gin.Context) {
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
