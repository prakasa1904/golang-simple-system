package group

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

func (r *Repository) GetById(db *gorm.DB, entity *Entity, id any) error {
	return db.Where("id = ?", id).Take(entity).Error
}

func (r *Repository) CountByName(db *gorm.DB, name any) (int64, error) {
	var total int64
	err := db.Model(new(Entity)).Where("name = ?", name).Count(&total).Error
	return total, err
}

func (r *Repository) GetByName(db *gorm.DB, group *Entity, name string) error {
	return db.Where("name = ?", name).First(group).Error
}

func (r *Repository) Find(db *gorm.DB, groups *[]Entity, filters map[string]string, limit int, order clause.OrderByColumn) error {
	// build query find with dynamic data filter
	buildQuery := db.Select(SelectColumn)

	for key, value := range filters {
		buildQuery = buildQuery.Where(key, value)
	}

	buildQuery = buildQuery.Limit(limit).Order(order).Find(groups)

	return buildQuery.Error
}
