package v1

type ApiGroup struct {
	SysUserApi
	SessionApi
}

var ApiGroupSys = new(ApiGroup)
