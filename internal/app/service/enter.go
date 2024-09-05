package service

type ServiceGroup struct {
	UserService
	SessionService
}

var ServiceGroupSys = new(ServiceGroup)
