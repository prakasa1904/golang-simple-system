package member

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewRepository(log *logrus.Logger) *Repository {
	return &Repository{
		Log: log,
	}
}

func (r *Repository) Create(db *gorm.DB, entity *Entity) error {
	return db.Create(entity).Error
}

func (r *Repository) Update(db *gorm.DB, entity *Entity) error {
	return db.Save(entity).Error
}

func (r *Repository) Delete(db *gorm.DB, entity *Entity) error {
	return db.Delete(entity).Error
}

func (r *Repository) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Model(new(Entity)).Where("id = ?", id).Count(&total).Error
	return total, err
}

func (r *Repository) FindById(db *gorm.DB, entity *Entity, id any) error {
	return db.Where("id = ?", id).Take(entity).Error
}

func (r *Repository) CountByUsername(db *gorm.DB, username any) (int64, error) {
	var total int64
	err := db.Model(new(Entity)).Where("username = ?", username).Count(&total).Error
	return total, err
}

func (r *Repository) FindByUsername(db *gorm.DB, user *Entity, username string) error {
	return db.Where("username = ?", username).First(user).Error
}

func (r *Repository) Find(db *gorm.DB, users *[]Entity, filters map[string]string, limit int, order clause.OrderByColumn) error {
	// build query find with dynamic data filter
	buildQuery := db.Select(SelectColumn)

	for key, value := range filters {
		buildQuery = buildQuery.Where(key, value)
	}

	buildQuery = buildQuery.Limit(limit).Order(order).Find(users)

	return buildQuery.Error
}
