package config

const (
	StackName               = "DeveloperSeriesStack"
	TableName               = "OrdersTable"
	PartitionKey            = "OrderID"
	SortKey                 = "OrderDateTime"
	CreateOrderFunctionName = "CreateOrder"
	CreateOrderCodePath     = "functions/insert-record/."
	ReadOrderFunctionName   = "ReadOrder"
	ReadOrderCodePath       = "functions/read-record/."
	UpdateOrderFunctionName = "UpdateOrder"
	UpdateOrderCodePath     = "functions/update-record/."
	DeleteFunctionName      = "DeleteOrder"
	DeleteOrderCodePath     = "functions/delete-record/."
	MemorySize              = 128
	MaxDuration             = 60
	Handler                 = "bootstrap"
)
