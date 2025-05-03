package ses

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/spf13/viper"
)

func InitSes(ctx context.Context) (*sesv2.Client, error) {
	defaultConfig, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(viper.GetString("aws.region")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(viper.GetString("aws.ses.access_key"), viper.GetString("aws.ses.access_secret_key"), "" )))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	
	client := sesv2.NewFromConfig(defaultConfig)
	return client, nil
}