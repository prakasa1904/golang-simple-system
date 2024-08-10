package order

import (
	"github.com/prakasa1904/panji-express/internal/services/member"
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

func (r *Repository) GetById(db *gorm.DB, entity *Entity, id any) error {
	return db.Joins("Member", func(db *gorm.DB) *gorm.DB {
		return db.Select(member.SelectColumn)
	}).Where("`order`.`id` = ?", id).Take(entity).Error
}

func (r *Repository) CountByMemberID(db *gorm.DB, memberID any) (int64, error) {
	var total int64
	err := db.Model(new(Entity)).Where("member_id = ?", memberID).Count(&total).Error
	return total, err
}

func (r *Repository) FindByMemberID(db *gorm.DB, order *Entity, memberID any) error {
	return db.Joins("Member", func(db *gorm.DB) *gorm.DB {
		return db.Select(member.SelectColumn)
	}).Where("`order`.`member_id` = ?", memberID).Take(order).Error
}

func (r *Repository) Find(db *gorm.DB, orders *[]Entity, filters map[string]string, limit int, order clause.OrderByColumn) error {
	// Add debug from db connection instance instead
	buildQuery := db.Debug()
	// build query find with dynamic data filter
	buildQuery = buildQuery.Select(SelectColumnWithJoin)

	// preload Group (from another query), because Entity data structure require to has Group struct
	// TODO: operate with single query
	buildQuery = buildQuery.Joins("Member", func(db *gorm.DB) *gorm.DB {
		return db.Select(member.SelectColumn)
	})

	for key, value := range filters {
		buildQuery = buildQuery.Where(key, value)
	}

	buildQuery = buildQuery.Limit(limit).Order(order).Find(orders)

	return buildQuery.Error
}
