package config

import "errors"

var (
	ErrMissingAWSAccessKey = errors.New("AWS_ACCESS_KEY_ID is required")
	ErrMissingAWSSecretKey = errors.New("AWS_SECRET_ACCESS_KEY is required")
	ErrMissingAWSRegion    = errors.New("AWS_REGION is required")
	ErrMissingAWSBucket    = errors.New("AWS_BUCKET_NAME is required")
)
