package v1

type ApiGroup struct {
	SysUserApi
	SessionApi
	SystemApi
}

var ApiGroupSys = new(ApiGroup)
