package response

import (
	"crispy-garbanzo/global"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	ERROR        = -1    //失败
	SUCCESS      = 0     // 成功
	TokenMissing = 20001 //Token缺失
	TokenExpired = 20002 //Token已过期
	TokenError   = 20003 //Token错误

	UserNotFound      = 30001 //用户不存在
	PasswordError     = 30002 //密码错误
	UserPasswordError = 30003 //用户名不存在或者密码错误
	UserLoginForbiden = 30004 //用户被禁止登录
	InvalidUserId     = 30005 //无效的用户ID
	UserNameAlready   = 30006 //用户名已注册
	OldPasswordError  = 30007 //原密码错误
	BalanceNoEnough   = 30008 //余额不足

	InvalidParameter  = 40001 //参数无效
	ObjectNotFound    = 40002 //对象不存在
	ObjectExisted     = 40003 //对象已存在
	NotfoundParameter = 40004 //参数不全或者类型错误

	InternalServerError = 50001 //服务异常
	QuestFinished       = 60001 //任务已完成
	QuestTypeError      = 60002 //任务类型有误

	ChannelCrowded      = 70001 //通道拥挤中，稍后再试
	BalanceNotEnough    = 70002 //余额不足,请充值
	ActivityNotFound    = 70101 //活动不存在
	ActivityGetIn       = 70102 //请须知，只能参加一次
	ActivityFullIn      = 70103 //当前场次参与人数已满
	FreeSpinAlreadyJoin = 70104 //已抽奖，请查看抽奖记录
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Messagei18(msgCode int, c *gin.Context) string {
	lang := c.Request.Header.Get("Access-Language")
	localizer := i18n.NewLocalizer(global.FPG_I18N, lang)
	msg := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    strconv.Itoa(msgCode),
			Other: "internal server error",
		},
	})
	return msg
}

func Result(code int, data interface{}, msgCode int, c *gin.Context) {
	msg := Messagei18(msgCode, c)
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, 0, c)
}

func OkWithMessage(message int, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, 0, c)
}

func OkWithDetailed(data interface{}, message int, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, -1, c)
}

func FailWithMessage(message int, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func NoAuth(message int, c *gin.Context) {
	msg := Messagei18(message, c)
	c.JSON(http.StatusUnauthorized, Response{
		7,
		nil,
		msg,
	})
}

func FailWithDetailed(data interface{}, message int, c *gin.Context) {
	Result(ERROR, data, message, c)
}
