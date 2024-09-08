package config

const (
	StackName           = "DeveloperSeriesStack"
	KinesisFunctionName = "LambdaKinesisFunction"
	MemorySize          = 128
	MaxDuration         = 60
	KinesisCodePath     = "function/kinesis-lambda/."
	Handler             = "bootstrap"
)
