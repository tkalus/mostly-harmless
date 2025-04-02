package awsconfig

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// CustomFunctionProvider implements the aws.CredentialsProvider interface
type CustomFunctionProvider struct {
	retrieve func(ctx context.Context) (aws.Credentials, error)
}

// NewCustomFunctionProvider
// NewCustomFunctionProvider initializes a new CustomFunctionProviderinstance and returns aws.CredentialsProvider interface.
func NewCustomFunctionProvider(
	retrieve func(ctx context.Context) (aws.Credentials, error),
) (aws.CredentialsProvider, error) {
	provider := &CustomFunctionProvider{
		retrieve: retrieve,
	}
	return provider, nil
}

// Retrieve implements the aws.CredentialsProvider interface method
func (p *CustomFunctionProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return p.retrieve(ctx)
}

// NewCustomFunctionConf initializes a new CustomFunctionConf instance and returns aws.Config interface.
func NewCustomFunctionConf(
	_ context.Context,
	cfg aws.Config,
	retrieve func(ctx context.Context) (aws.Credentials, error),
) (aws.Config, error) {
	credProvider, err := NewCustomFunctionProvider(retrieve)
	if err != nil {
		return aws.Config{}, err
	}
	credentials := aws.NewCredentialsCache(
		credProvider,
		func(options *aws.CredentialsCacheOptions) {
			options.ExpiryWindow = 5 * time.Minute
		},
	)

	config := cfg.Copy()
	config.Credentials = credentials
	return config, nil
}
