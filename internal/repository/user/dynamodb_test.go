package user

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/yarieldis/guebapi/internal/models"
)

// fakeDynamo is a minimal in-memory implementation of dynamoAPI used for tests.
type fakeDynamo struct {
	items map[string]models.User
}

func newFakeDynamo() *fakeDynamo {
	return &fakeDynamo{items: map[string]models.User{}}
}

func sValue(item map[string]types.AttributeValue, key string) string {
	if av, ok := item[key].(*types.AttributeValueMemberS); ok {
		return av.Value
	}
	return ""
}

func (f *fakeDynamo) PutItem(_ context.Context, params *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	u := models.User{
		Username: sValue(params.Item, "username"),
		Password: sValue(params.Item, "password"),
	}
	if _, exists := f.items[u.Username]; exists {
		return nil, &types.ConditionalCheckFailedException{}
	}
	f.items[u.Username] = u
	return &dynamodb.PutItemOutput{}, nil
}

func (f *fakeDynamo) GetItem(_ context.Context, params *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	username := sValue(params.Key, "username")
	u, exists := f.items[username]
	if !exists {
		return &dynamodb.GetItemOutput{}, nil
	}
	return &dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{Value: u.Username},
			"password": &types.AttributeValueMemberS{Value: u.Password},
		},
	}, nil
}

func (f *fakeDynamo) UpdateItem(_ context.Context, params *dynamodb.UpdateItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	username := sValue(params.Key, "username")
	u, exists := f.items[username]
	if !exists {
		return nil, &types.ConditionalCheckFailedException{}
	}
	u.Password = sValue(params.ExpressionAttributeValues, ":pwd")
	f.items[username] = u
	return &dynamodb.UpdateItemOutput{}, nil
}

func newTestRepo() (*DynamoDBRepository, *fakeDynamo) {
	fake := newFakeDynamo()
	return &DynamoDBRepository{client: fake, tableName: "users"}, fake
}

func TestDynamoDBRepository_Create(t *testing.T) {
	repo, _ := newTestRepo()
	ctx := context.Background()

	u := &models.User{Username: "alice", Password: "secret"}

	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := repo.Create(ctx, u); err != ErrUserAlreadyExists {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestDynamoDBRepository_GetByUsername(t *testing.T) {
	repo, fake := newTestRepo()
	ctx := context.Background()

	fake.items["alice"] = models.User{Username: "alice", Password: "secret"}

	got, err := repo.GetByUsername(ctx, "alice")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Username != "alice" || got.Password != "secret" {
		t.Fatalf("unexpected user: %+v", got)
	}

	if _, err := repo.GetByUsername(ctx, "missing"); err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestDynamoDBRepository_UpdatePassword(t *testing.T) {
	repo, fake := newTestRepo()
	ctx := context.Background()

	fake.items["alice"] = models.User{Username: "alice", Password: "old"}

	if err := repo.UpdatePassword(ctx, "alice", "new"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fake.items["alice"].Password != "new" {
		t.Fatalf("expected password to be updated, got %q", fake.items["alice"].Password)
	}

	if err := repo.UpdatePassword(ctx, "missing", "new"); err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestDynamoDBRepository_Exists(t *testing.T) {
	repo, fake := newTestRepo()
	ctx := context.Background()

	fake.items["alice"] = models.User{Username: "alice", Password: "secret"}

	exists, err := repo.Exists(ctx, "alice")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !exists {
		t.Fatal("expected user to exist")
	}

	exists, err = repo.Exists(ctx, "missing")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if exists {
		t.Fatal("expected user not to exist")
	}
}
