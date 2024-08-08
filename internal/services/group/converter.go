package group

import (
	"strconv"
)

func GroupToResponse(group *Entity) *Response {
	return &Response{
		ID:        group.ID,
		Name:      group.Name,
		Status:    group.Status,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	}
}

func RequestToEntity(payload *RequestPayload) (*Entity, error) {
	entity := new(Entity)
	// ignore error on converting ID, insert no require to has ID
	id, _ := strconv.ParseUint(payload.ID, 0, 64)

	status, err := strconv.Atoi(payload.Status)
	if err != nil {
		return entity, err
	}

	entity.ID = id
	entity.Name = payload.Name
	entity.Status = status

	return entity, nil
}
