package order

var (
	SelectColumn       = "id, invoice_id, status, meta_file, description"
	AllowedFilterQuery = []string{"id", "description", "invoice_id"}
)

// status order, order by first time created until delivered
var (
	StatusCreated            = 0
	StatusCourierOnPick      = 1
	StatusPickedUp           = 2
	StatusShippingWarehouse  = 3
	StatusSending            = 4
	StatusReceivingWarehosue = 5
	StatusDelivered          = 6
)
