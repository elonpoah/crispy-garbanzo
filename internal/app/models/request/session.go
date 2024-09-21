package request

import (
	"crispy-garbanzo/common/request"
)

type SessionDetailReq struct {
	Id uint `json:"id" binding:"required"`
}
type CheckInviteDutyReq struct {
	Type int `json:"type" binding:"required"` // 1:daily 2:weekly 3:monthly
}
type SessionListReq struct {
	request.PageInfo
	Type int `json:"type" binding:"required" ` // 1:hight bonus 2:hight rate 3:hot"
}

type GameHistoryReq struct {
	request.PageInfo
	Status uint `json:"status"`
}
