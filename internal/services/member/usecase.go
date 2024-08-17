package member

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

	total, err := c.Repository.CountByUsername(tx, request.Username)
	if err != nil {
		c.Log.Warnf("Failed count member from database : %+v", err)
		return nil, err
	}

	if total > 0 {
		c.Log.Warnf("Member already exists : %+v", err)
		return nil, errors.New("member already exists")
	}

	if request.Password != "" {
		strongPass, err := CreatePassword(request.Password)
		if err != nil {
			c.Log.Warnf("Failed to generate bcrypt hash : %+v", err)
			return nil, err
		}

		// encrypt password
		request.Password = strongPass
		request.ConfirmPassword = strongPass
	}

	// new member
	newmember, err := RequestPayloadToEntity(request)
	if err != nil {
		c.Log.Warnf("Failed convert payload to entity : %+v", err)
		return nil, err
	}

	if err := c.Repository.Create(tx, newmember); err != nil {
		c.Log.Warnf("Failed create member to database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	return EntityToResponse(newmember), nil
}

func (c *UseCase) Find(ctx context.Context, filters map[string]string, limit int, order clause.OrderByColumn) (*[]Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	members := new([]Entity)
	err := c.Repository.Find(tx, members, filters, limit, order)
	if err != nil {
		c.Log.Warnf("Failed count member from database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	// append entity to response
	var memberResponse = new([]Response)
	for _, val := range *members {
		member := EntityToResponse(&val)
		*memberResponse = append(*memberResponse, *member)
	}

	return memberResponse, nil
}

func (c *UseCase) GetById(ctx context.Context, id any) (*Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	member := new(Entity)
	err := c.Repository.GetById(tx, member, id)
	if err != nil {
		c.Log.Warnf("Failed get member by ID from database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	// convert entity to response
	memberResp := EntityToResponse(member)

	return memberResp, nil
}

func (c *UseCase) GetByGroupName(ctx context.Context, groupName any) (*Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	member := new(Entity)
	err := c.Repository.GetByGroupName(tx, member, groupName)
	if err != nil {
		c.Log.Warnf("Failed get member by group name from database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	// convert entity to response
	memberResp := EntityToResponse(member)

	return memberResp, nil
}

func (c *UseCase) Update(ctx context.Context, request *RequestPayload) (*Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, err
	}

	// update password if exist
	if request.Password != "" {
		strongPass, err := CreatePassword(request.Password)
		if err != nil {
			c.Log.Warnf("Failed to generate bcrypt hash : %+v", err)
			return nil, err
		}

		// encrypt password
		request.Password = strongPass
		request.ConfirmPassword = strongPass
	}

	member, err := RequestPayloadToEntity(request)
	if err != nil {
		c.Log.Warnf("Failed convert payload to entity : %+v", err)
		return nil, err
	}

	if err := c.Repository.Update(tx, member); err != nil {
		c.Log.Warnf("Failed create member to database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	return EntityToResponse(member), nil
}

func (c *UseCase) Delete(ctx context.Context, request *DeletePayload) (*Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, err
	}

	member, err := DeletePayloadToEntity(request)
	if err != nil {
		c.Log.Warnf("Failed convert payload to entity : %+v", err)
		return nil, err
	}

	if err := c.Repository.Delete(tx, member); err != nil {
		c.Log.Warnf("Failed delete member to database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	return EntityToResponse(member), nil
}
