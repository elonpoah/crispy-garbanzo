package response

type SysUserResponse struct {
	ID            uint    `json:"uid"`
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
