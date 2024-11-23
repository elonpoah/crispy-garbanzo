package response

type SysUserResponse struct {
	ID            uint    `json:"uid"`
	Pid           uint    `json:"pid"`
	Username      string  `json:"userName"`
	NickName      string  `json:"nickName"`
	Phone         string  `json:"phone"`
	Email         string  `json:"email"`
	Balance       float64 `json:"balance"`
	FreezeBalance float64 `json:"freezeBalance"`
}

type LoginResponse struct {
	User  SysUserResponse `json:"user"`
	Token string          `json:"token"`
}

type DrawDetailRes struct {
	Username     string  `json:"userName"`
	DrawId       string  `json:"drawId"`
	BonusType    int     `json:"bonusType"`
	Bonus        float64 `json:"bonus"`
	Distribute   float64 `json:"distribute"`
	Count        int     `json:"count"`
	Participants uint    `json:"participants"`
	Status       uint    `json:"status"`
}
