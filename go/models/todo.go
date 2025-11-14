package models

type Todo struct {
	UserID          string `dynamodbav:"userId"`
	TodoID          string `dynamodbav:"todoId"`
	Text            string `dynamodbav:"text"`
	Done            bool   `dynamodbav:"done"`
	CreatedAt       int64  `dynamodbav:"createdAt"`
	UpdatedAt       int64  `dynamodbav:"updatedAt"`
	DueAt           *int64 `dynamodbav:"dueAt,omitempty"`
	CalendarEventID string `dynamodbav:"calendarEventId,omitempty"`
}
