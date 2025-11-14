package models

type PetState struct {
	UserID            string  `dynamodbav:"userId"`
	Mood              string  `dynamodbav:"mood"`        // grumpy | neutral | golden
	Personality       string  `dynamodbav:"personality"` // supportive | sarcastic | chill
	StreakDays        int64   `dynamodbav:"streakDays"`
	LastLoginDay      string  `dynamodbav:"lastLoginDay"`      // e.g. "2025-11-14"
	LastInteractionAt int64   `dynamodbav:"lastInteractionAt"` // unix ms
	CompletionScore   float64 `dynamodbav:"completionScore"`   // optional aggregate
}
