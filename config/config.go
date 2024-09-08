package config

const (
	StackName         = "DeveloperSeriesStack"
	TaskAFunctionName = "LambdaFunctionTaskA"
	TaskBFunctionName = "LambdaFunctionTaskB"
	MemorySize        = 128
	MaxDuration       = 60
	TaskACodePath     = "function/task-a/."
	TaskBCodePath     = "function/task-b/."
	Handler           = "bootstrap"
)
