package models

import (
	"time"
)

const (
	MANAGE_NEWS_PERM = 1 << iota
	VIEW_ACTIONS_PERM
	MANAGE_USERS_PERM
	MANAGE_QUIZZES_PERM
)

type User struct {
	Id           int32     `json:"id" db:"id"`
	UserName     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Role struct {
	Id          int32     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Permissions int64     `json:"permissions" db:"permissions"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type UserRole struct {
	UserId int32 `json:"user_id" db:"user_id"`
	RoleId int32 `json:"role_id" db:"role_id"`
}

type UserAction struct {
	Id                int32     `json:"id" db:"id"`
	UserId            *int32    `json:"user_id" db:"user_id"`
	ActionDescription string    `json:"action_description" db:"action_description"`
	ActionTime        time.Time `json:"action_time" db:"action_time"`
}

type Category struct {
	Id   int32  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Quiz struct {
	Id          int32     `json:"id" db:"id"`
	AuthorId    *int32    `json:"author_id" db:"author_id"`
	Title       string    `json:"title" db:"title"`
	Description *string   `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type QuizCategory struct {
	QuizId     int32 `json:"quiz_id" db:"quiz_id"`
	CategoryId int32 `json:"category_id" db:"category_id"`
}

type Question struct {
	Id           int32  `json:"id" db:"id"`
	QuizId       int32  `json:"quiz_id" db:"quiz_id"`
	QuestionText string `json:"question_text" db:"question_text"`
	QuestionType string `json:"question_type" db:"question_type"`
}

type TextQuestionAnswer struct {
	QuestionId  int32  `json:"question_id" db:"question_id"`
	RightAnswer string `json:"right_answer" db:"right_answer"`
}

type Choice struct {
	Id         int32  `json:"id" db:"id"`
	QuestionId int32  `json:"question_id" db:"question_id"`
	ChoiceText string `json:"choice_text" db:"choice_text"`
	IsCorrect  bool   `json:"is_correct" db:"is_correct"`
}

type ChoiceQuestionAnswer struct {
	QuestionId    int32 `json:"question_id" db:"question_id"`
	RightChoiceId int32 `json:"right_choice_id" db:"right_choice_id"`
}

type ChoiceAnswer struct {
	UserId     int32     `json:"user_id" db:"user_id"`
	QuestionId int32     `json:"question_id" db:"question_id"`
	ChoiceId   *int32    `json:"choice_id" db:"choice_id"`
	AnsweredAt time.Time `json:"answered_at" db:"answered_at"`
}

type TextAnswer struct {
	UserId     int32     `json:"user_id" db:"user_id"`
	QuestionId int32     `json:"question_id" db:"question_id"`
	TextAnswer *string   `json:"text_answer" db:"text_answer"`
	AnsweredAt time.Time `json:"answered_at" db:"answered_at"`
}

type UserQuizParticipation struct {
	UserId             int32 `json:"user_id" db:"user_id"`
	QuizId             int32 `json:"quiz_id" db:"quiz_id"`
	ParticipationCount int   `json:"participation_count" db:"participation_count"`
}

type UserQuizScore struct {
	UserId         int32     `json:"user_id" db:"user_id"`
	QuizId         int32     `json:"quiz_id" db:"quiz_id"`
	Score          float64   `json:"score" db:"score"`
	LastUpdateTime time.Time `json:"last_update_time" db:"last_update_time"`
}

type QuizParticipationTime struct {
	Id                  int32      `json:"id" db:"id"`
	ParticipationNumber int        `json:"participation_number" db:"participation_number"`
	UserId              int32      `json:"user_id" db:"user_id"`
	QuizId              int32      `json:"quiz_id" db:"quiz_id"`
	StartedAt           time.Time  `json:"started_at" db:"started_at"`
	FinishedAt          *time.Time `json:"finished_at" db:"finished_at"`
}

type QuizStatistics struct {
	QuizId                int32     `json:"quiz_id" db:"quiz_id"`
	TotalAttempts         int       `json:"total_attempts" db:"total_attempts"`
	AverageScore          float64   `json:"average_score" db:"average_score"`
	AverageCompletionTime *string   `json:"average_completion_time" db:"average_completion_time"`
	LastUpdateTime        time.Time `json:"last_update_time" db:"last_update_time"`
}

type News struct {
	Id        int32     `json:"id" db:"id"`
	AuthorId  int32     `json:"author_id" db:"author_id"`
	Title     string    `json:"title" db:"title"`
	NewsText  string    `json:"news_text" db:"news_text"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
