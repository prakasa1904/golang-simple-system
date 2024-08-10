package group

func EntityToResponse(group *Entity) *Response {
	// convert uint64 to int64 to resolve bug with plush (ejs template)
	var id int64 = int64(group.ID)

	return &Response{
		ID:        id,
		Name:      group.Name,
		Status:    group.Status,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	}
}

func RequestToEntity(payload *RequestPayload) (*Entity, error) {
	entity := new(Entity)

	entity.ID = payload.ID
	entity.Name = payload.Name
	entity.Status = payload.Status

	return entity, nil
}
