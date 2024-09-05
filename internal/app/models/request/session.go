package request

import (
	"crispy-garbanzo/common/request"
)

type SessionDetailReq struct {
	Id int `json:"id" binding:"required"`
}

type SessionListReq struct {
	request.PageInfo
	Type int `json:"type" binding:"required" ` // 1:hight bonus 2:hight rate 3:hot"
}
