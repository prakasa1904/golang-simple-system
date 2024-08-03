package group

import "time"

type Entity struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement:true"`
	Name      string    `gorm:"column:name;size:256"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:current_timestamp;autoUpdateTime"`
}

func (a *Entity) TableName() string {
	return "order"
}
