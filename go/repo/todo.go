package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/juhun32/patriot25-gochi/go/models"
)

type TodoRepo struct {
	client    *dynamodb.Client
	tableName string
}

func NewTodoRepo(client *dynamodb.Client, tableName string) *TodoRepo {
	return &TodoRepo{
		client:    client,
		tableName: tableName,
	}
}

func (r *TodoRepo) CreateTodo(ctx context.Context, userID, todoID, text string, dueAt *int64) (*models.Todo, error) {
	now := time.Now().UnixMilli()
	todo := &models.Todo{
		UserID:    userID,
		TodoID:    todoID,
		Text:      text,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
		DueAt:     dueAt,
	}

	item, err := attributevalue.MarshalMap(todo)
	if err != nil {
		return nil, err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           &r.tableName,
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(userId) AND attribute_not_exists(todoId)"),
	})
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *TodoRepo) ListTodos(ctx context.Context, userID string) ([]models.Todo, error) {
	keyCond, err := attributevalue.MarshalMap(map[string]string{
		"userId": userID,
	})
	if err != nil {
		return nil, err
	}

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:     &r.tableName,
		KeyConditions: map[string]types.Condition{"userId": {ComparisonOperator: types.ComparisonOperatorEq, AttributeValueList: []types.AttributeValue{keyCond["userId"]}}},
	})
	if err != nil {
		return nil, err
	}

	var todos []models.Todo
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &todos); err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *TodoRepo) UpdateTodoDone(ctx context.Context, userID, todoID string, done bool) error {
	now := time.Now().UnixMilli()

	key, err := attributevalue.MarshalMap(map[string]string{
		"userId": userID,
		"todoId": todoID,
	})
	if err != nil {
		return err
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        &r.tableName,
		Key:              key,
		UpdateExpression: aws.String("SET done = :d, updatedAt = :u"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":d": &types.AttributeValueMemberBOOL{Value: done},
			":u": &types.AttributeValueMemberN{Value: fmt.Sprint(now)},
		},
	})
	return err
}

func (r *TodoRepo) DeleteTodo(ctx context.Context, userID, todoID string) error {
	key, err := attributevalue.MarshalMap(map[string]string{
		"userId": userID,
		"todoId": todoID,
	})
	if err != nil {
		return err
	}

	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &r.tableName,
		Key:       key,
	})
	return err
}
