package member

func EntityToResponse(member *Entity) *Response {
	return &Response{
		ID:        member.ID,
		Fullname:  member.Fullname,
		Username:  member.Username,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}
}
