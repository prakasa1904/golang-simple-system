package order

import (
	"context"

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

func (c *UseCase) Create(ctx context.Context, request *CreatePayload) (*Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, err
	}

	// new order
	neworder, err := CreatePayloadToEntity(request)
	if err != nil {
		c.Log.Warnf("Failed convert payload to entity : %+v", err)
		return nil, err
	}

	if err := c.Repository.Create(tx, neworder); err != nil {
		c.Log.Warnf("Failed create order to database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	return EntityToResponse(neworder), nil
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

	// append entity to response
	var orderResponse = new([]Response)
	for _, val := range *orders {
		order := EntityToResponse(&val)
		*orderResponse = append(*orderResponse, *order)
	}

	return orderResponse, nil
}

func (c *UseCase) GetById(ctx context.Context, id any) (*Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	order := new(Entity)
	err := c.Repository.GetById(tx, order, id)
	if err != nil {
		c.Log.Warnf("Failed get order by ID from database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	// convert entity to response
	orderResp := EntityToResponse(order)

	return orderResp, nil
}

func (c *UseCase) Update(ctx context.Context, request *UpdatePayload) (*Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, err
	}

	// generate invoice

	order, err := UpdatePayloadToEntity(request)
	if err != nil {
		c.Log.Warnf("Failed convert payload to entity : %+v", err)
		return nil, err
	}

	// set invoice ID when status is pickedUp by courier
	invoiceID := CreateInvoice(order.ID, order.CreatedAt)
	order.InvoiceID = invoiceID

	if err := c.Repository.Update(tx, order); err != nil {
		c.Log.Warnf("Failed create order to database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	return EntityToResponse(order), nil
}

func (c *UseCase) Delete(ctx context.Context, request *DeletePayload) (*Response, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, err
	}

	order, err := DeletePayloadToEntity(request)
	if err != nil {
		c.Log.Warnf("Failed convert payload to entity : %+v", err)
		return nil, err
	}

	if err := c.Repository.Delete(tx, order); err != nil {
		c.Log.Warnf("Failed delete order to database : %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	return EntityToResponse(order), nil
}
