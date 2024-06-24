package response

import (
	system "crispy-garbanzo/internal/admin/models"
)

type SysUserResponse struct {
	User system.SysUser `json:"user"`
}

type LoginResponse struct {
	User  system.SysUser `json:"user"`
	Token string         `json:"token"`
}
