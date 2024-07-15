package member

func UserToResponse(user *Entity) *Response {
	return &Response{
		ID:        user.ID,
		Fullname:  user.Fullname,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
