package config

const (
	StackName               = "DeveloperSeriesStack"
	TableName               = "OrdersTable"
	PartitionKey            = "OrderID"
	SortKey                 = "OrderDateTime"
	CreateOrderFunctionName = "CreateOrder"
	CreateOrderCodePath     = "function/insert-record/."
	ReadOrderFunctionName   = "ReadOrder"
	ReadOrderCodePath       = "function/read-record/."
	UpdateOrderFunctionName = "UpdateOrder"
	UpdateOrderCodePath     = "function/update-record/."
	DeleteFunctionName      = "DeleteOrder"
	DeleteOrderCodePath     = "function/delete-record/."
	MemorySize              = 128
	MaxDuration             = 60
	Handler                 = "bootstrap"
)
