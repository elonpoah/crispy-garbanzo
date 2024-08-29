package router

type RouterGroup struct {
	BaseApiRouter
	UserApiRouter
}

var AppRouterGroup = new(RouterGroup)
