package group

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UseCase struct {
	DB         *gorm.DB
	Log        *logrus.Logger
	Validate   *validator.Validate
	Repository *Repository
}

func NewUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate, repository *Repository) *UseCase {
	return &UseCase{
		DB:         db,
		Log:        logger,
		Validate:   validate,
		Repository: repository,
	}
}

func (c *UseCase) Create(ctx context.Context, request *RequestPayload) (*Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, err
	}

	total, err := c.Repository.CountByName(tx, request.Name)
	if err != nil {
		c.Log.Warnf("Failed count group from database : %+v", err)
		return nil, err
	}

	if total > 0 {
		c.Log.Warnf("Group already exists : %+v", err)
		return nil, errors.New("group already exists")
	}

	group, err := RequestToEntity(request)
	if err != nil {
		c.Log.Warnf("Failed convert payload to entity : %+v", err)
		return nil, err
	}

	if err := c.Repository.Create(tx, group); err != nil {
		c.Log.Warnf("Failed create group to database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	return GroupToResponse(group), nil
}

func (c *UseCase) Find(ctx context.Context, filters map[string]string, limit int, order clause.OrderByColumn) (*[]Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	groups := new([]Entity)
	err := c.Repository.Find(tx, groups, filters, limit, order)
	if err != nil {
		c.Log.Warnf("Failed count group from database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	// map to response
	var response = new([]Response)
	for _, group := range *groups {
		groupItem := GroupToResponse(&group)
		*response = append(*response, *groupItem)
	}

	return response, nil
}

func (c *UseCase) GetById(ctx context.Context, id string) (*Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	group := new(Entity)
	err := c.Repository.GetById(tx, group, id)
	if err != nil {
		c.Log.Warnf("Failed count group from database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	return GroupToResponse(group), nil
}

func (c *UseCase) Update(ctx context.Context, request *RequestPayload) (*Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, err
	}

	group, err := RequestToEntity(request)
	if err != nil {
		c.Log.Warnf("Failed convert payload to entity : %+v", err)
		return nil, err
	}

	if err := c.Repository.Update(tx, group); err != nil {
		c.Log.Warnf("Failed create group to database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	return GroupToResponse(group), nil
}

func (c *UseCase) Delete(ctx context.Context, request *RequestPayload) (*Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, err
	}

	group, err := RequestToEntity(request)
	if err != nil {
		c.Log.Warnf("Failed convert payload to entity : %+v", err)
		return nil, err
	}

	if err := c.Repository.Delete(tx, group); err != nil {
		c.Log.Warnf("Failed delete group to database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	return GroupToResponse(group), nil
}
