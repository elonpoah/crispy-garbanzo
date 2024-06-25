package response

type SysUserResponse struct {
	Username string `json:"userName"`
	NickName string `json:"nickName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type LoginResponse struct {
	User  SysUserResponse `json:"user"`
	Token string          `json:"token"`
}
