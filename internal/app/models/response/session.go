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
