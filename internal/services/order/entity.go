package order

import (
	"time"

	"gorm.io/datatypes"
)

type Entity struct {
	ID          uint64                  `gorm:"column:id;primaryKey;autoIncrement:true"`
	InvoiceID   string                  `gorm:"column:invoice_id;size:256"` // create when order created
	Status      int                     `gorm:"column:status;size:256"`
	MetaFile    datatypes.JSONType[any] `gorm:"column:meta_file"` // fill with file data structure
	Description string                  `gorm:"column:description;size:256"`
	CreatedAt   time.Time               `gorm:"column:created_at;not null;default:current_timestamp"`
	UpdatedAt   time.Time               `gorm:"column:updated_at;not null;default:current_timestamp;autoUpdateTime"`
	MemberID    uint64                  `gorm:"column:member_id;not null"`
}

func (a *Entity) TableName() string {
	return "order"
}
