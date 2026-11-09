package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	appconfig "github.com/yarieldis/guebapi/config"
	"github.com/yarieldis/guebapi/internal/app"
	"github.com/yarieldis/guebapi/internal/repository/user"
)

func main() {
	cfg := appconfig.Load()

	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(cfg.AWS.Region))
	if err != nil {
		log.Fatalf("failed to load AWS configuration: %v", err)
	}

	dynamoClient := dynamodb.NewFromConfig(awsCfg)
	userRepo := user.NewDynamoDBRepository(dynamoClient, cfg.Storage.DynamoDBTable)

	r := app.BuildRouter(cfg, userRepo)

	// Proxy API Gateway HTTP API (payload v2) events into the Gin engine.
	lambda.Start(ginadapter.NewV2(r).ProxyWithContext)
}
