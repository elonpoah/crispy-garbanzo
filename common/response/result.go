package response

import (
	// "crispy-garbanzo/global"
	"net/http"
	// "strconv"

	"github.com/gin-gonic/gin"
	// "github.com/nicksnyder/go-i18n/v2/i18n"
)

// type Response struct {
// 	Code int         `json:"code"`
// 	Data interface{} `json:"data"`
// 	Msg  string      `json:"msg"`
// }

// type NullMap struct {
// }

// const (
// 	Success             = 0     //成功
// 	Fail                = 10000 //失败
// 	TokenMissing        = 20001 //Token缺失
// 	TokenExpired        = 20002 //Token已过期
// 	TokenError          = 20003 //Token错误
// 	UserNotFound        = 30001 //用户不存在
// 	PasswordError       = 30002 //密码错误
// 	InvalidParameter    = 40001 //参数无效
// 	ObjectNotFound      = 40002 //对象不存在
// 	ObjectExisted       = 40003 //对象已存在
// 	InternalServerError = 50001 //服务器内部错误
// 	QuestFinished       = 60001 //任务已完成
// 	QuestTypeError      = 60002 //任务类型有误
// )

// func Result(code int, data interface{}, c *gin.Context) Response {
// 	status, other := http.StatusOK, "success"
// 	if code != Success {
// 		status = http.StatusInternalServerError
// 		other = "internal server error"
// 	}
// 	if nil == data {
// 		data = NullMap{}
// 	}

// 	lang := c.Request.Header.Get("Access-Language")
// 	localizer := i18n.NewLocalizer(global.FPG_I18N, lang)
// 	msg := localizer.MustLocalize(&i18n.LocalizeConfig{
// 		DefaultMessage: &i18n.Message{
// 			ID:    strconv.Itoa(code),
// 			Other: other,
// 		},
// 	})
// 	res := Response{
// 		code,
// 		data,
// 		msg,
// 	}
// 	c.JSON(status, res)
// 	return res
// }

//	func Failed(code int, c *gin.Context) Response {
//		return Result(code, nil, c)
//	}
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "查询成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func NoAuth(message string, c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		7,
		nil,
		message,
	})
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}
