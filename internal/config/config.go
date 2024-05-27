package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func GetAwsConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func GetS3BucketName() string {
	return readEnvVar("AWS_S3_BUCKET_NAME", "super-secret-bucket-name")
}

func GetGrpcPort() string {
	return readEnvVar("GRPC_SERVER_PORT", "50051")
}

func readEnvVar(envVar, suggestion string) string {
	if value, ok := os.LookupEnv(envVar); ok {
		return value
	}
	panic(fmt.Sprintf("env var %s is not set, suggested value: %s", envVar, suggestion))
}
