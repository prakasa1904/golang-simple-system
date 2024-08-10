package member

import (
	"github.com/prakasa1904/panji-express/internal/services/group"
	"golang.org/x/crypto/bcrypt"
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

	entity.ID = payload.ID
	entity.Fullname = payload.Fullname
	entity.Username = payload.Username
	entity.Email = payload.Email
	entity.Password = payload.Password
	entity.GroupID = payload.GroupID

	return entity, nil
}

func DeletePayloadToEntity(payload *DeletePayload) (*Entity, error) {
	entity := new(Entity)

	entity.ID = payload.ID

	return entity, nil
}

func CreatePassword(plainPassword string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(password), err
}
