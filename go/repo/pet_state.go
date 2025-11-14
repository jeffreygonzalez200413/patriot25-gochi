package repo

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/juhun32/patriot25-gochi/go/models"
)

type PetStateRepo struct {
	client    *dynamodb.Client
	tableName string
}

func NewPetStateRepo(client *dynamodb.Client, tableName string) *PetStateRepo {
	return &PetStateRepo{
		client:    client,
		tableName: tableName,
	}
}

func (r *PetStateRepo) GetPetState(ctx context.Context, userID string) (*models.PetState, error) {
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

	var state models.PetState
	if err := attributevalue.UnmarshalMap(out.Item, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *PetStateRepo) UpsertPetState(ctx context.Context, state *models.PetState) error {
	if state.LastInteractionAt == 0 {
		state.LastInteractionAt = time.Now().UnixMilli()
	}
	item, err := attributevalue.MarshalMap(state)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	})
	return err
}
