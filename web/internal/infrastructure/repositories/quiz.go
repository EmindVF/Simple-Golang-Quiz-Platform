package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"quiz_platform/internal/database"
	"quiz_platform/internal/misc/apperrors"
	"quiz_platform/internal/models"

	"github.com/lib/pq"
)

type SqlQuizRepository struct {
	DBProvider database.SqlDatabaseProvider
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	query :=
		`SELECT
		id, name
		FROM categories`

	rows, err := repo.DBProvider.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, &apperrors.ErrInternal{Message: err.Error()}
	}
	defer rows.Close()

	allCategories := make([]*models.Category, 0)
	for rows.Next() {
		var category models.Category
		err = rows.Scan(
			&category.Id, &category.Name)
		if err != nil {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}

		allCategories = append(allCategories, &category)
	}

	return allCategories, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) GetCategoriesPairs(ctx context.Context, quizIds []int32) ([]int32, []int32, error) {
	query :=
		`SELECT
		quiz_id, category_id
		FROM quiz_categories WHERE quiz_id = ANY($1::int[])`

	rows, err := repo.DBProvider.QueryContext(
		ctx,
		query,
		pq.Array(quizIds),
	)
	if err != nil {
		return nil, nil, &apperrors.ErrInternal{Message: err.Error()}
	}
	defer rows.Close()

	resQuizIds := make([]int32, 0)
	resCategoryIds := make([]int32, 0)
	for rows.Next() {
		var quizId int32
		var categoryId int32
		err = rows.Scan(
			&quizId, &categoryId)
		if err != nil {
			return nil, nil, &apperrors.ErrInternal{Message: err.Error()}
		}

		resQuizIds = append(resQuizIds, quizId)
		resCategoryIds = append(resCategoryIds, categoryId)
	}

	return resQuizIds, resCategoryIds, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) GetAllQuizzes(ctx context.Context, categoryId int32) ([]*models.Quiz, error) {
	var query string
	if categoryId == 0 {
		query = "SELECT id, author_id, title, description, created_at, updated_at FROM quizzes"
	} else {
		query = `
            SELECT q.id, q.author_id, q.title, q.description, q.created_at, q.updated_at 
            FROM quizzes q 
            JOIN quiz_categories qc ON q.id = qc.quiz_id 
            WHERE qc.category_id = $1`
	}
	args := []any{}
	if categoryId != 0 {
		args = append(args, categoryId)
	}
	rows, err := repo.DBProvider.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return nil, &apperrors.ErrInternal{Message: err.Error()}
	}
	defer rows.Close()

	allQuizzes := make([]*models.Quiz, 0)
	for rows.Next() {
		var quiz models.Quiz
		err = rows.Scan(
			&quiz.Id, &quiz.AuthorId, &quiz.Title,
			&quiz.Description, &quiz.CreatedAt, &quiz.UpdatedAt)
		if err != nil {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}

		allQuizzes = append(allQuizzes, &quiz)
	}

	return allQuizzes, nil
}

// May return ErrInternal or ErrInvalidInput on failure.
func (repo *SqlQuizRepository) AddQuiz(ctx context.Context, title string, desc string, author_id int32) (int32, error) {
	var id int32
	err := repo.DBProvider.QueryRowContext(
		ctx,
		"INSERT INTO quizzes (author_id, title, description) VALUES ($1, $2, $3) RETURNING id",
		author_id, title, desc).Scan(&id)
	if err != nil {
		return 0, &apperrors.ErrInternal{Message: err.Error()}
	}
	return id, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) RemoveQuizCategories(ctx context.Context, quizId int32) error {
	_, err := repo.DBProvider.ExecContext(
		ctx,
		`DELETE FROM quiz_categories WHERE quiz_id = $1`,
		quizId)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	return nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) AddQuizCategories(ctx context.Context, quizId int32, categoryIds []int32) error {
	if len(categoryIds) == 0 {
		return nil
	}

	query := strings.Builder{}
	_, err := query.Write(
		[]byte(
			`INSERT INTO quiz_categories
		(quiz_id, category_id)
		VALUES `))
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	values := []any{}

	placeHolderArray := make([]any, 2)
	for i := int32(1); i <= 2; i++ {
		placeHolderArray[i-1] = i
	}

	query.Write([]byte(fmt.Sprintf(" ($%d, $%d) ",
		placeHolderArray...)))
	values = append(values,
		quizId, categoryIds[0])

	for i := 1; i < len(categoryIds); i++ {
		for j := range 2 {
			placeHolderArray[j] = i*2 + j + 1
		}
		query.Write([]byte(fmt.Sprintf(", ($%d, $%d) ",
			placeHolderArray...)))
		values = append(values,
			quizId, categoryIds[i])
	}

	_, err = repo.DBProvider.ExecContext(
		ctx,
		query.String(),
		values...,
	)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	return nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) AddQuestion(ctx context.Context, quizId int32, text string, qtype string) (int32, error) {
	var id int32
	err := repo.DBProvider.QueryRowContext(
		ctx,
		"INSERT INTO questions (quiz_id, question_text, question_type) VALUES ($1, $2, $3) RETURNING id",
		quizId, text, qtype).Scan(&id)
	if err != nil {
		return 0, &apperrors.ErrInternal{Message: err.Error()}
	}
	return id, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) AddTextQuestionAnswer(ctx context.Context, qId int32, answer string) (int32, error) {
	var id int32
	err := repo.DBProvider.QueryRowContext(
		ctx,
		"INSERT INTO text_question_answers (question_id, right_answer) VALUES ($1, $2) RETURNING question_id",
		qId, answer).Scan(&id)
	if err != nil {
		return 0, &apperrors.ErrInternal{Message: err.Error()}
	}
	return id, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) AddChoice(ctx context.Context, qId int32, text string, isCorrect bool) (int32, error) {
	var id int32
	err := repo.DBProvider.QueryRowContext(
		ctx,
		"INSERT INTO choices (question_id, choice_text, is_correct) VALUES ($1, $2, $3) RETURNING id",
		qId, text, isCorrect).Scan(&id)
	if err != nil {
		return 0, &apperrors.ErrInternal{Message: err.Error()}
	}
	return id, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) GetQuiz(ctx context.Context, id int32) (*models.Quiz, error) {
	quiz := &models.Quiz{}
	err := repo.DBProvider.QueryRowContext(
		ctx,
		`SELECT id, author_id, title, description, created_at, updated_at
		FROM quizzes 
		WHERE id = $1 `,
		id).Scan(
		&quiz.Id, &quiz.AuthorId, &quiz.Title,
		&quiz.Description, &quiz.CreatedAt, &quiz.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperrors.ErrNotFound{Message: "content not found"}
		} else {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}
	}

	return quiz, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) GetQuizQuestions(ctx context.Context, id int32) ([]*models.Question, error) {
	query :=
		`SELECT
		id, quiz_id, question_text, question_type
		FROM questions WHERE quiz_id = $1`

	rows, err := repo.DBProvider.QueryContext(
		ctx,
		query,
		id,
	)
	if err != nil {
		return nil, &apperrors.ErrInternal{Message: err.Error()}
	}
	defer rows.Close()

	allQuestions := make([]*models.Question, 0)
	for rows.Next() {
		var question models.Question
		err = rows.Scan(
			&question.Id, &question.QuizId, &question.QuestionText, &question.QuestionType)
		if err != nil {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}

		allQuestions = append(allQuestions, &question)
	}

	return allQuestions, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) GetTextQuestionAnswer(ctx context.Context, questionId int32) (*models.TextQuestionAnswer, error) {
	answer := &models.TextQuestionAnswer{}
	err := repo.DBProvider.QueryRowContext(
		ctx,
		`SELECT question_id, right_answer
		FROM text_question_answers 
		WHERE question_id = $1 `,
		questionId).Scan(
		&answer.QuestionId, &answer.RightAnswer)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperrors.ErrNotFound{Message: "content not found"}
		} else {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}
	}
	return answer, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) GetChoices(ctx context.Context, questionId int32) ([]*models.Choice, error) {
	query :=
		`SELECT
		id, question_id, choice_text, is_correct
		FROM choices WHERE question_id = $1`

	rows, err := repo.DBProvider.QueryContext(
		ctx,
		query,
		questionId,
	)
	if err != nil {
		return nil, &apperrors.ErrInternal{Message: err.Error()}
	}
	defer rows.Close()

	allChoices := make([]*models.Choice, 0)
	for rows.Next() {
		var choice models.Choice
		err = rows.Scan(
			&choice.Id, &choice.QuestionId, &choice.ChoiceText, &choice.IsCorrect)
		if err != nil {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}

		allChoices = append(allChoices, &choice)
	}

	return allChoices, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) DeleteQuiz(ctx context.Context, id int32) error {
	_, err := repo.DBProvider.ExecContext(
		ctx,
		`DELETE FROM quizzes WHERE
		id = $1`, id)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	return nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) GetCorrectChoice(ctx context.Context, questionId int32) (*models.Choice, error) {
	choice := &models.Choice{}
	err := repo.DBProvider.QueryRowContext(
		ctx,
		`SELECT id, question_id, choice_text, is_correct
		FROM choices 
		WHERE question_id = $1 AND is_correct = true `,
		questionId).Scan(
		&choice.Id, &choice.QuestionId, &choice.ChoiceText,
		&choice.IsCorrect)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperrors.ErrNotFound{Message: "content not found"}
		} else {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}
	}

	return choice, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) GetLastParticipationTime(ctx context.Context, userId int32) (*models.QuizParticipationTime, error) {
	partTime := &models.QuizParticipationTime{}
	err := repo.DBProvider.QueryRowContext(
		ctx,
		`SELECT id, participation_number, user_id, quiz_id, started_at, finished_at
		FROM quiz_participation_times
		WHERE user_id = $1
		AND participation_number = (
			SELECT MAX(participation_number)
			FROM quiz_participation_times
			WHERE user_id = $1
		);`,
		userId).Scan(
		&partTime.Id, &partTime.ParticipationNumber, &partTime.UserId,
		&partTime.QuizId, &partTime.StartedAt, &partTime.FinishedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperrors.ErrNotFound{Message: "content not found"}
		} else {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}
	}

	return partTime, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) DeleteParticipationTime(ctx context.Context, id int32) error {
	_, err := repo.DBProvider.ExecContext(
		ctx,
		`DELETE FROM quiz_participation_times WHERE
		id = $1`, id)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	return nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) UpdateParticipationTime(ctx context.Context, id int32, finishTime time.Time) error {
	res, err := repo.DBProvider.ExecContext(
		ctx,
		`UPDATE quiz_participation_times SET
		finished_at = $1
		WHERE id = $2`,
		finishTime, id)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}
	if rowsAffected == 0 {
		return &apperrors.ErrNotFound{Message: "news not found"}
	}

	return nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) AddParticipationTime(ctx context.Context, userId int32, quizId int32, startTime time.Time) (int32, error) {
	var id int32
	err := repo.DBProvider.QueryRowContext(
		ctx,
		`INSERT INTO quiz_participation_times (user_id, quiz_id, started_at, participation_number)
		VALUES (
			$1,
			$2,
			$3,
			COALESCE((SELECT MAX(participation_number) FROM quiz_participation_times WHERE user_id = $1), 0) + 1
		)
		RETURNING id;`,
		userId, quizId, startTime).Scan(&id)
	if err != nil {
		return 0, &apperrors.ErrInternal{Message: err.Error()}
	}
	return id, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) RemoveUserChoiceAnswers(ctx context.Context, questionId int32, userId int32) error {
	_, err := repo.DBProvider.ExecContext(
		ctx,
		`DELETE FROM choice_answers WHERE
		question_id = $1 AND user_id = $2`, questionId, userId)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	return nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) AddUserChoiceAnswer(ctx context.Context, userId int32, questionId int32, choiceId int32) error {
	var id int32
	err := repo.DBProvider.QueryRowContext(
		ctx,
		"INSERT INTO choice_answers (user_id, question_id, choice_id) VALUES ($1, $2, $3) RETURNING user_id",
		userId, questionId, choiceId).Scan(&id)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}
	return nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) RemoveUserTextAnswers(ctx context.Context, questionId int32, userId int32) error {
	_, err := repo.DBProvider.ExecContext(
		ctx,
		`DELETE FROM text_answers WHERE
		question_id = $1 AND user_id = $2`, questionId, userId)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	return nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) AddUserTextAnswer(ctx context.Context, userId int32, questionId int32, text string) error {
	var id int32
	err := repo.DBProvider.QueryRowContext(
		ctx,
		"INSERT INTO text_answers (user_id, question_id, text_answer) VALUES ($1, $2, $3) RETURNING user_id",
		userId, questionId, text).Scan(&id)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}
	return nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) UpsertUserScore(ctx context.Context, userId int32, quizId int32, score float32, time time.Time) error {
	res, err := repo.DBProvider.ExecContext(
		ctx,
		`INSERT INTO user_quiz_scores
		(user_id, quiz_id, score, last_update_time)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, quiz_id)
		DO UPDATE SET
		score = $3, last_update_time = $4`,
		userId, quizId, score, time,
	)
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &apperrors.ErrInternal{Message: err.Error()}
	}
	if rowsAffected < 1 {
		return &apperrors.ErrNotFound{Message: "task result not found"}
	}

	return nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) GetQuizStatistics(ctx context.Context, quizIds []int32) ([]*models.QuizStatistics, error) {
	query :=
		`SELECT
		quiz_id, total_attempts, average_score, average_completion_time, last_update_time
		FROM quiz_statistics WHERE quiz_id = ANY($1::int[])`

	rows, err := repo.DBProvider.QueryContext(
		ctx,
		query,
		pq.Array(quizIds),
	)
	if err != nil {
		return nil, &apperrors.ErrInternal{Message: err.Error()}
	}
	defer rows.Close()

	allStats := make([]*models.QuizStatistics, 0)
	for rows.Next() {
		var stats models.QuizStatistics
		err = rows.Scan(
			&stats.QuizId, &stats.TotalAttempts, &stats.AverageScore,
			&stats.AverageCompletionTime, &stats.LastUpdateTime)
		if err != nil {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}

		allStats = append(allStats, &stats)
	}

	return allStats, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) GetUserScore(ctx context.Context, userId int32, quizId int32) (*models.UserQuizScore, error) {
	score := &models.UserQuizScore{}
	err := repo.DBProvider.QueryRowContext(
		ctx,
		`SELECT
		user_id, quiz_id, score, last_update_time
		FROM user_quiz_scores WHERE user_id = $1 AND quiz_id = $2`,
		userId, quizId).Scan(
		&score.UserId, &score.QuizId, &score.Score, &score.LastUpdateTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperrors.ErrNotFound{Message: "content not found"}
		} else {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}
	}

	return score, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) GetQuizParticipationTime(ctx context.Context, userId int32, quizId int32) (*models.QuizParticipationTime, error) {
	choice := &models.QuizParticipationTime{}
	err := repo.DBProvider.QueryRowContext(
		ctx,
		`SELECT
			user_id, quiz_id, started_at, finished_at, participation_number
		FROM
			quiz_participation_times
		WHERE
			user_id = $1 AND quiz_id = $2 AND finished_at IS NOT NULL
		ORDER BY
			participation_number DESC
		LIMIT 1`,
		userId, quizId).Scan(
		&choice.UserId, &choice.QuizId, &choice.StartedAt, &choice.FinishedAt, &choice.ParticipationNumber)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperrors.ErrNotFound{Message: "content not found"}
		} else {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}
	}

	return choice, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) GetUserTextAnswer(ctx context.Context, userId int32, questionId int32) (*models.TextAnswer, error) {
	choice := &models.TextAnswer{}
	err := repo.DBProvider.QueryRowContext(
		ctx,
		`SELECT
		user_id, question_id, text_answer
		FROM text_answers WHERE user_id = $1 AND question_id = $2`,
		userId, questionId).Scan(
		&choice.UserId, &choice.QuestionId, &choice.TextAnswer)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperrors.ErrNotFound{Message: "content not found"}
		} else {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}
	}

	return choice, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) GetUserChoiceAnswer(ctx context.Context, userId int32, questionId int32) (*models.ChoiceAnswer, error) {
	choice := &models.ChoiceAnswer{}
	err := repo.DBProvider.QueryRowContext(
		ctx,
		`SELECT
		user_id, question_id, choice_id
		FROM choice_answers WHERE user_id = $1 AND question_id = $2`,
		userId, questionId).Scan(
		&choice.UserId, &choice.QuestionId, &choice.ChoiceId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperrors.ErrNotFound{Message: "content not found"}
		} else {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}
	}

	return choice, nil
}

// May return ErrInternal or ErrNotFound on failure.
func (repo *SqlQuizRepository) GetChoice(ctx context.Context, id int32) (*models.Choice, error) {
	choice := &models.Choice{}
	err := repo.DBProvider.QueryRowContext(
		ctx,
		`SELECT
		id, question_id, choice_text, is_correct
		FROM choices WHERE id = $1`,
		id).Scan(
		&choice.Id, &choice.QuestionId, &choice.ChoiceText, &choice.IsCorrect)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &apperrors.ErrNotFound{Message: "content not found"}
		} else {
			return nil, &apperrors.ErrInternal{Message: err.Error()}
		}
	}

	return choice, nil
}
