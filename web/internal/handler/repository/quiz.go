package repository

import (
	"context"
	"quiz_platform/internal/models"
	"time"
)

var (
	QuizRepositoryInstance QuizRepository
)

type QuizRepository interface {
	// May return ErrInternal or ErrNotFound on failure.
	GetAllCategories(ctx context.Context) ([]*models.Category, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetCategoriesPairs(ctx context.Context, quizIds []int32) ([]int32, []int32, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetAllQuizzes(ctx context.Context, categoryId int32) ([]*models.Quiz, error)

	// May return ErrInternal or ErrInvalidInput on failure.
	AddQuiz(ctx context.Context, title string, desc string, author_id int32) (int32, error)

	// May return ErrInternal or ErrNotFound on failure.
	RemoveQuizCategories(ctx context.Context, quizId int32) error

	// May return ErrInternal or ErrNotFound on failure.
	AddQuizCategories(ctx context.Context, quizId int32, categoryIds []int32) error

	// May return ErrInternal or ErrNotFound on failure.
	AddQuestion(ctx context.Context, quizId int32, text string, qtype string) (int32, error)

	// May return ErrInternal or ErrNotFound on failure.
	AddTextQuestionAnswer(ctx context.Context, qId int32, answer string) (int32, error)

	// May return ErrInternal or ErrNotFound on failure.
	AddChoice(ctx context.Context, qId int32, text string, isCorrect bool) (int32, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetQuiz(ctx context.Context, id int32) (*models.Quiz, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetQuizQuestions(ctx context.Context, id int32) ([]*models.Question, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetTextQuestionAnswer(ctx context.Context, questionId int32) (*models.TextQuestionAnswer, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetChoices(ctx context.Context, questionId int32) ([]*models.Choice, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetCorrectChoice(ctx context.Context, questionId int32) (*models.Choice, error)

	// May return ErrInternal or ErrNotFound on failure.
	DeleteQuiz(ctx context.Context, id int32) error

	// May return ErrInternal or ErrNotFound on failure.
	GetLastParticipationTime(ctx context.Context, userId int32) (*models.QuizParticipationTime, error)

	// May return ErrInternal or ErrNotFound on failure.
	DeleteParticipationTime(ctx context.Context, id int32) error

	// May return ErrInternal or ErrNotFound on failure.
	UpdateParticipationTime(ctx context.Context, id int32, finishTime time.Time) error

	// May return ErrInternal or ErrNotFound on failure.
	AddParticipationTime(ctx context.Context, userId int32, quizId int32, startTime time.Time) (int32, error)

	// May return ErrInternal or ErrNotFound on failure.
	RemoveUserChoiceAnswers(ctx context.Context, questionId int32, userId int32) error

	// May return ErrInternal or ErrNotFound on failure.
	AddUserChoiceAnswer(ctx context.Context, userId int32, questionId int32, choiceId int32) error

	// May return ErrInternal or ErrNotFound on failure.
	RemoveUserTextAnswers(ctx context.Context, questionId int32, userId int32) error

	// May return ErrInternal or ErrNotFound on failure.
	AddUserTextAnswer(ctx context.Context, userId int32, questionId int32, text string) error

	// May return ErrInternal or ErrNotFound on failure.
	UpsertUserScore(ctx context.Context, userId int32, quizId int32, score float32, time time.Time) error

	// May return ErrInternal or ErrNotFound on failure.
	GetQuizStatistics(ctx context.Context, quizIds []int32) ([]*models.QuizStatistics, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetUserScore(ctx context.Context, userId int32, quizId int32) (*models.UserQuizScore, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetQuizParticipationTime(ctx context.Context, userId int32, quizId int32) (*models.QuizParticipationTime, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetUserTextAnswer(ctx context.Context, userId int32, questionId int32) (*models.TextAnswer, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetUserChoiceAnswer(ctx context.Context, userId int32, questionId int32) (*models.ChoiceAnswer, error)

	// May return ErrInternal or ErrNotFound on failure.
	GetChoice(ctx context.Context, id int32) (*models.Choice, error)
}
