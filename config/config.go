package config

const (
	StackName    = "DeveloperSeriesStack"
	FunctionName = "ProcessingSQSFunction"
	MemorySize   = 128
	MaxDuration  = 60
	CodePath     = "function/."
	Handler      = "bootstrap"
)
