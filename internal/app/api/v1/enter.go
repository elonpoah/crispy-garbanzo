package v1

type ApiGroup struct {
	SysUserApi
	SessionApi
	SystemApi
	DrawApi
}

var ApiGroupSys = new(ApiGroup)
