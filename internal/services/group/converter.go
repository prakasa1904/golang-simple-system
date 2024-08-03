package group

func GroupToResponse(group *Entity) *Response {
	return &Response{
		ID:        group.ID,
		Name:      group.Name,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	}
}
