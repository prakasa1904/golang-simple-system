package member

import (
	"strconv"

	"github.com/prakasa1904/panji-express/internal/services/group"
)

func EntityToResponse(member *Entity) *Response {
	groupResponse := group.EntityToResponse(&member.Group)

	// convert uint64 to int64 to resolve bug with plush (ejs template)
	var id int64 = int64(member.ID)
	var gid int64 = int64(member.GroupID)

	return &Response{
		ID:              id,
		Fullname:        member.Fullname,
		Username:        member.Username,
		Email:           member.Email,
		Password:        member.Password,
		ConfirmPassword: member.Password,
		CreatedAt:       member.CreatedAt,
		UpdatedAt:       member.UpdatedAt,
		GroupID:         gid,
		Group:           *groupResponse,
	}
}

func RequestPayloadToEntity(payload *RequestPayload) (*Entity, error) {
	entity := new(Entity)
	// ignore error on converting ID, insert no require to has ID
	id, _ := strconv.ParseUint(payload.ID, 0, 64)
	gid, _ := strconv.ParseUint(payload.GroupID, 0, 64)

	entity.ID = id
	entity.Fullname = payload.Fullname
	entity.Username = payload.Username
	entity.Email = payload.Email
	entity.Password = payload.Password
	entity.GroupID = gid

	return entity, nil
}

func DeletePayloadToEntity(payload *DeletePayload) (*Entity, error) {
	entity := new(Entity)
	// ignore error on converting ID, insert no require to has ID
	id, _ := strconv.ParseUint(payload.ID, 0, 64)

	entity.ID = id

	return entity, nil
}
