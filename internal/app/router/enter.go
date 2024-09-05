package router

type RouterGroup struct {
	BaseApiRouter
	UserApiRouter
	SessionApiRouter
}

var AppRouterGroup = new(RouterGroup)
