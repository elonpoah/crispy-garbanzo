package router

type RouterGroup struct {
	BaseApiRouter
	UserApiRouter
	SessionApiRouter
	DrawApiRouter
}

var AppRouterGroup = new(RouterGroup)
