package service

type ServiceGroup struct {
	UserService
	SessionService
	SystemService
}

var ServiceGroupSys = new(ServiceGroup)
