package order

import (
	"context"
	"errors"

	"github.com/devetek/golang-webapp-boilerplate/internal/services/order"
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

func (c *UseCase) Create(ctx context.Context, request *CreateRequest) (*Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, err
	}

	total, err := c.Repository.CountByDescription(tx, request.Description)
	if err != nil {
		c.Log.Warnf("Failed count order from database : %+v", err)
		return nil, err
	}

	if total > 0 {
		c.Log.Warnf("Order already exists : %+v", err)
		return nil, errors.New("order already exists")
	}

	// new order
	order := &Entity{
		Description: request.Description,
		MetaFile:    request.MetaFile,
		Status:      order.StatusCreated,
	}

	if err := c.Repository.Create(tx, order); err != nil {
		c.Log.Warnf("Failed create order to database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	return OrderToResponse(order), nil
}

func (c *UseCase) Find(ctx context.Context, filters map[string]string, limit int, order clause.OrderByColumn) (*[]Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	orders := new([]Entity)
	err := c.Repository.Find(tx, orders, filters, limit, order)
	if err != nil {
		c.Log.Warnf("Failed count order from database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	// map to response
	var ordersResp = new([]Response)
	for _, order := range *orders {
		orderItem := OrderToResponse(&order)
		*ordersResp = append(*ordersResp, *orderItem)
	}

	return ordersResp, nil
}
