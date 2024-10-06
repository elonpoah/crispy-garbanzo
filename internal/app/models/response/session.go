package response

type CheckSessionStautsResponse struct {
	IsGot bool `json:"isGot"`
}

type UserSummaryResponse struct {
	SessionCount int64 `json:"sessionCount"`
	FreeCount    int64 `json:"freeCount"`
}

type InviteSessionResponse struct {
	Registrations int   `json:"registrations"`
	Participates  int64 `json:"participates"`
}
type InviteConfig struct {
	Daily InviteItemConfig `json:"daily" binding:"required"`
	Week  InviteItemConfig `json:"week" binding:"required"`
	Month InviteItemConfig `json:"month" binding:"required"`
}
type InviteItemConfig struct {
	Bonus        float64 `json:"bonus" binding:"required"`
	Count        uint    `json:"count" binding:"required"`
	Participants uint    `json:"participants" binding:"required"`
	Enable       uint    `json:"enable" binding:"required,oneof=1 2"` // 1开启2关闭
}
