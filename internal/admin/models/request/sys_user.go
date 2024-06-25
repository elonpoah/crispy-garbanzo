package request

type Register struct {
	Username string `json:"userName" binding:"required"`
	Password string `json:"passWord" binding:"required"`
	NickName string `json:"nickName" binding:"required"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

// User login structure
type Login struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}
