package bandinputs

type RegisterInput struct {
	Title       string  `json:"title" validate:"required,min=2"`
	Description string  `json:"description" validate:"required,min=8"`
	Logo        *string `json:"logo" validate:"omitempty,url"`
}

type UpdateInviteInput struct {
	Status string `json:"status"` // "accepted", "denied"
}

type UpdateMemberInput struct {
	Role string `json:"role"` // "member", "admin"
}
