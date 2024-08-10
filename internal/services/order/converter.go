package order

import (
	"fmt"
	"strconv"
	"time"

	"github.com/prakasa1904/panji-express/internal/services/member"
	"gorm.io/datatypes"
)

func EntityToResponse(order *Entity) *Response {
	// create your own converter if entity and response data require to has their own data structure
	return &Response{
		ID:          order.ID,
		InvoiceID:   order.InvoiceID,
		Description: order.Description,
		MetaFile:    order.MetaFile,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt.String(),
		UpdatedAt:   order.UpdatedAt.String(),
		MemberID:    order.MemberID,
		Member:      *member.EntityToResponse(&order.Member),
	}
}

func CreatePayloadToEntity(payload *CreatePayload) (*Entity, error) {
	entity := new(Entity)

	entity.MemberID = payload.MemberID
	entity.Description = payload.Description
	entity.MetaFile = datatypes.NewJSONType(payload.MetaFile)

	return entity, nil
}

func UpdatePayloadToEntity(payload *UpdatePayload) (*Entity, error) {
	entity := new(Entity)

	// generate invoice when courier pickedUp
	// TODO: this method is not complete yet, need improvement
	if payload.Status == StatusPickedUp {
		createDate, err := time.Parse(time.RFC3339, payload.CreatedAt)
		if err != nil {
			return entity, err
		}

		entity.InvoiceID = CreateInvoice(payload.ID, createDate)
	}

	entity.ID = payload.ID
	entity.Description = payload.Description
	entity.MetaFile = datatypes.NewJSONType(payload.MetaFile)
	entity.Status = payload.Status
	entity.MemberID = payload.MemberID

	return entity, nil
}

func DeletePayloadToEntity(payload *DeletePayload) (*Entity, error) {
	entity := new(Entity)

	id, err := strconv.ParseUint(payload.ID, 0, 64)
	if err != nil {
		return entity, err
	}

	entity.ID = id

	return entity, nil
}

func CreateInvoice(orderID uint64, createDate time.Time) string {
	// invoice formula
	return fmt.Sprintf("INV/PA/%d/%d/%d/%d", orderID, createDate.Day(), createDate.Month(), createDate.Year())
}
