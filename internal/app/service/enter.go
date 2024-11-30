package service

type ServiceGroup struct {
	UserService
	SessionService
	SystemService
	DrawService
}

var ServiceGroupSys = new(ServiceGroup)
