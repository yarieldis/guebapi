package user

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/yarieldis/guebapi/internal/models"
)

// dynamoAPI is the subset of the DynamoDB client used by this repository.
// It exists to allow testing with an in-memory fake.
type dynamoAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}

// DynamoDBRepository implements Repository using an AWS DynamoDB table.
type DynamoDBRepository struct {
	client    dynamoAPI
	tableName string
}

// NewDynamoDBRepository creates a DynamoDB-backed user repository.
func NewDynamoDBRepository(client *dynamodb.Client, tableName string) *DynamoDBRepository {
	return &DynamoDBRepository{client: client, tableName: tableName}
}

func key(username string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"username": &types.AttributeValueMemberS{Value: username},
	}
}

// Create stores a new user, failing with ErrUserAlreadyExists if the username is taken.
func (r *DynamoDBRepository) Create(ctx context.Context, u *models.User) error {
	item, err := attributevalue.MarshalMap(u)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(username)"),
	})
	if err != nil {
		var ccfe *types.ConditionalCheckFailedException
		if errors.As(err, &ccfe) {
			return ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

// GetByUsername returns the user with the given username, or ErrUserNotFound.
func (r *DynamoDBRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key:       key(username),
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil || len(out.Item) == 0 {
		return nil, ErrUserNotFound
	}

	var u models.User
	if err := attributevalue.UnmarshalMap(out.Item, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdatePassword updates a user's password, failing with ErrUserNotFound if the user does not exist.
func (r *DynamoDBRepository) UpdatePassword(ctx context.Context, username, password string) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        aws.String(r.tableName),
		Key:              key(username),
		UpdateExpression: aws.String("SET password = :pwd"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pwd": &types.AttributeValueMemberS{Value: password},
		},
		ConditionExpression: aws.String("attribute_exists(username)"),
	})
	if err != nil {
		var ccfe *types.ConditionalCheckFailedException
		if errors.As(err, &ccfe) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

// Exists reports whether a user with the given username exists.
func (r *DynamoDBRepository) Exists(ctx context.Context, username string) (bool, error) {
	_, err := r.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
