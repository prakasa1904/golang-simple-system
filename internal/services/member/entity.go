package member

import "time"

type Entity struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement:true"`
	Fullname  string    `gorm:"column:fullname;size:256"`
	Username  string    `gorm:"column:username;size:256;unique"`
	Email     string    `gorm:"column:email;size:256;unique"`
	Password  string    `gorm:"column:password;size:256"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:current_timestamp;autoUpdateTime"`
	GroupID   uint64    `gorm:"column:group_id;primaryKey;autoIncrement:true"`
}

func (a *Entity) TableName() string {
	return "member"
}
