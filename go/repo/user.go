package repo

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/juhun32/patriot25-gochi/go/models"
)

type UserRepo struct {
	client    *dynamodb.Client
	tableName string
}

func NewUserRepo(client *dynamodb.Client, tableName string) *UserRepo {
	return &UserRepo{
		client:    client,
		tableName: tableName,
	}
}

func (r *UserRepo) UpsertUser(ctx context.Context, user *models.User) error {
	now := time.Now().UnixMilli()
	if user.CreatedAt == 0 {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
		// You *could* add a ConditionExpression if you want
		// to distinguish between create/update, but not required
	})
	return err
}

func (r *UserRepo) GetUser(ctx context.Context, userID string) (*models.User, error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		"userId": userID,
	})
	if err != nil {
		return nil, err
	}

	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}

	var user models.User
	if err := attributevalue.UnmarshalMap(out.Item, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
