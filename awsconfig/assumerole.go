package awsconfig

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

const (
	errParseIAMRoleArn      = "Cannot parse passed IAM Role ARN"
	errStsGetCallerIdentity = "Cannot determine caller identity of passed aws.Config"
)

// NewAssumeRoleConf returns an aws.Config configured to assume the given roleArn
// using auto-refreshing credentials and optional AssumeRoleOptions.
func NewAssumeRoleConf(
	ctx context.Context,
	cfg aws.Config,
	roleArn string,
	opts ...func(*stscreds.AssumeRoleOptions),
) (aws.Config, error) {
	// Validate role ARN
	if _, err := arn.Parse(roleArn); err != nil {
		return aws.Config{}, fmt.Errorf("%v: %w", errParseIAMRoleArn, err)
	}

	// Create STS client from base config
	stsClient := sts.NewFromConfig(cfg)
	_, err := stsClient.GetCallerIdentity(ctx, nil)
	if err != nil {
		return aws.Config{}, fmt.Errorf("%v: %w", errStsGetCallerIdentity, err)
	}

	// Construct assume-role provider
	provider := stscreds.NewAssumeRoleProvider(stsClient, roleArn, opts...)

	// Wrap in auto-refreshing cache
	cached := aws.NewCredentialsCache(provider)

	// Return a copy of the config with assumed credentials
	newCfg := cfg.Copy()
	newCfg.Credentials = cached
	return newCfg, nil
}

// WithRoleSessionName sets the session name
func WithRoleSessionName(name string) func(*stscreds.AssumeRoleOptions) {
	return func(o *stscreds.AssumeRoleOptions) {
		o.RoleSessionName = name
	}
}

// WithDuration sets the session duration
func WithDuration(duration time.Duration) func(*stscreds.AssumeRoleOptions) {
	return func(o *stscreds.AssumeRoleOptions) {
		o.Duration = duration
	}
}

// WithExternalID sets the external ID
func WithExternalID(externalID string) func(*stscreds.AssumeRoleOptions) {
	return func(o *stscreds.AssumeRoleOptions) {
		o.ExternalID = aws.String(externalID)
	}
}

// WithPolicy sets an inline session policy
func WithPolicy(policy string) func(*stscreds.AssumeRoleOptions) {
	return func(o *stscreds.AssumeRoleOptions) {
		o.Policy = aws.String(policy)
	}
}

// WithPolicyArns sets managed policy ARNs
func WithPolicyArns(arns []string) func(*stscreds.AssumeRoleOptions) {
	var inputPolicyARNs []types.PolicyDescriptorType
	for _, arn := range arns {
		inputPolicyARNs = append(
			inputPolicyARNs,
			types.PolicyDescriptorType{
				Arn: aws.String(arn),
			},
		)
	}
	return func(o *stscreds.AssumeRoleOptions) {
		o.PolicyARNs = inputPolicyARNs
	}
}

// WithSourceIdentity sets the source identity
func WithSourceIdentity(id string) func(*stscreds.AssumeRoleOptions) {
	return func(o *stscreds.AssumeRoleOptions) {
		o.SourceIdentity = aws.String(id)
	}
}

// WithTags attaches session tags
func WithTags(tags map[string]string) func(*stscreds.AssumeRoleOptions) {
	var inputTags []types.Tag
	for key, value := range tags {
		inputTags = append(inputTags, types.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		})
	}
	return func(o *stscreds.AssumeRoleOptions) {
		o.Tags = inputTags
	}
}

// WithTransitiveTagKeys specifies transitive tag keys
func WithTransitiveTagKeys(keys []string) func(*stscreds.AssumeRoleOptions) {
	return func(o *stscreds.AssumeRoleOptions) {
		o.TransitiveTagKeys = keys
	}
}

// WithMFA sets the MFA serial number and token provider
func WithMFA(serial string, tokenProvider func() (string, error)) func(*stscreds.AssumeRoleOptions) {
	return func(o *stscreds.AssumeRoleOptions) {
		o.SerialNumber = aws.String(serial)
		o.TokenProvider = tokenProvider
	}
}
