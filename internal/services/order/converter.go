package order

func OrderToResponse(order *Entity) *Response {
	// create your own converter if entity and response data require to has their own data structure
	return &Response{
		ID:          order.ID,
		InvoiceID:   order.InvoiceID,
		Description: order.Description,
		MetaFile:    order.MetaFile,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}
