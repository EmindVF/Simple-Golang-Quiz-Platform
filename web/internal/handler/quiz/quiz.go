package quiz

import (
	"context"
	"fmt"
	"net/http"
	"quiz_platform/internal/handler/repository"
	"quiz_platform/internal/middleware"
	"quiz_platform/internal/misc/apperrors"
	"quiz_platform/internal/utility"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Choice struct {
	Id         int32  `json:"-"`
	QuestionId int32  `json:"-"`
	Text       string `json:"text" binding:"required"`
	IsCorrect  bool   `json:"is_correct"`
}

type Question struct {
	Id          int32    `json:"-"`
	Text        string   `json:"text" binding:"required"`
	Type        string   `json:"type" binding:"required,oneof=choice text"`
	Choices     []Choice `json:"choices,omitempty"`
	RightAnswer string   `json:"right_answer,omitempty"`
}

type Quiz struct {
	Id          int32      `json:"-"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description" binding:"required"`
	Categories  []string   `json:"categories" binding:"required"`
	Questions   []Question `json:"questions" binding:"required,dive"`

	TotalAttempts int32  `json:"-"`
	AverageScore  string `json:"-"`
	AverageTime   string `json:"-"`
}

type Submission struct {
	QuizId  int32             `json:"quiz_id"`
	Answers map[string]string `json:"answers"`
}

type QuizResult struct {
	Title       string
	Description string
	Score       string
	Time        string
	Questions   []AnsweredQuestion `json:"questions" binding:"required,dive"`
}

type AnsweredQuestion struct {
	Text        string `json:"text" binding:"required"`
	Type        string `json:"type" binding:"required,oneof=choice text"`
	RightAnswer string `json:"right_answer,omitempty"`
	UserAnswer  string
	IsCorrect   bool
}

func formatDuration(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	milliseconds := int(d.Milliseconds()) % 1000

	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d:%04d", hours, minutes, seconds, milliseconds)
}

func QuizCreatePostHandler(c *gin.Context) {

	var (
		authorId int32
		quiz     Quiz
	)
	data, ok := c.Get("sessionData")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	if sessionData, ok := data.(*middleware.SessionData); ok {
		authorId = sessionData.UserId
	}
	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	categoryIds := make([]int32, len(quiz.Categories))
	for i, v := range quiz.Categories {
		categoryId, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		categoryIds[i] = int32(categoryId)
	}
	if len(categoryIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
		return
	}

	ctx := context.Background()
	err := repository.TransactionManager.Run(ctx, func(ctx context.Context) error {

		quizId, err := repository.QuizRepositoryInstance.
			AddQuiz(ctx, quiz.Title, quiz.Description, authorId)
		if err != nil {
			return err
		}

		err = repository.QuizRepositoryInstance.
			AddQuizCategories(ctx, quizId, categoryIds)
		if err != nil {
			return err
		}

		if len(quiz.Questions) == 0 {
			return fmt.Errorf("invalid question count")
		}

		for i, v := range quiz.Questions {
			quiz.Questions[i].Id, err = repository.QuizRepositoryInstance.
				AddQuestion(ctx, quizId, v.Text, v.Type)
			if err != nil {
				return err
			}
		}

		for _, v := range quiz.Questions {
			if v.Type == "text" {
				_, err = repository.QuizRepositoryInstance.
					AddTextQuestionAnswer(ctx, v.Id, v.RightAnswer)
				if err != nil {
					return err
				}
			} else {
				if len(v.Choices) == 0 {
					return fmt.Errorf("invalid choice count")
				}
				for _, c := range v.Choices {
					_, err = repository.QuizRepositoryInstance.
						AddChoice(ctx, v.Id, c.Text, c.IsCorrect)
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/quiz")
}

func QuizCreateFormGetHandler(c *gin.Context) {
	ctx := context.Background()
	categories, err := repository.QuizRepositoryInstance.GetAllCategories(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	baseHInterface, _ := c.Get("BaseH")
	baseH, _ := baseHInterface.(*gin.H)
	c.HTML(http.StatusOK, "quiz_form.html", utility.MergeMaps(*baseH, gin.H{
		"title":      "Quizzes",
		"categories": categories}))
}

func QuizParticipationPostHandler(c *gin.Context) {

	var (
		userId     int32
		quizId     int32
		submission Submission
		err        error
	)
	data, ok := c.Get("sessionData")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	if sessionData, ok := data.(*middleware.SessionData); ok {
		userId = sessionData.UserId
	}
	if err := c.ShouldBindJSON(&submission); err != nil {
		println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	quizId = submission.QuizId

	questionMap := make(map[int32]string)
	for key, value := range submission.Answers {
		questionId, err := strconv.ParseInt(key, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		questionMap[int32(questionId)] = value
	}

	ctx := context.Background()
	err = repository.TransactionManager.Run(ctx, func(ctx context.Context) error {

		partTime, err := repository.QuizRepositoryInstance.GetLastParticipationTime(ctx, userId)
		if _, ok := err.(*apperrors.ErrNotFound); ok {
			return fmt.Errorf("no participation")
		} else if err == nil {
			if partTime.QuizId != int32(quizId) {
				println(quizId)
				println(partTime.QuizId)
				return fmt.Errorf("invalid quiz")
			}
			err := repository.QuizRepositoryInstance.UpdateParticipationTime(ctx, partTime.Id, time.Now().UTC())
			if err != nil {
				return err
			}
		} else {
			return err
		}

		questionModels, err := repository.QuizRepositoryInstance.
			GetQuizQuestions(ctx, int32(quizId))
		if err != nil {
			return err
		}
		amountOfQuestions := float32(len(questionModels))
		rightAnswers := float32(0)

		for _, q := range questionModels {
			if answ, ok := questionMap[q.Id]; ok {
				if q.QuestionType != "text" {
					choiceId, err := strconv.ParseInt(answ, 10, 32)
					if err != nil {
						return err
					}
					correctChoice, err := repository.QuizRepositoryInstance.GetCorrectChoice(ctx, q.Id)
					if err != nil {
						return err
					}
					if correctChoice.Id == int32(choiceId) {
						rightAnswers += 1
					}
					err = repository.QuizRepositoryInstance.RemoveUserChoiceAnswers(ctx, q.Id, userId)
					if err != nil {
						return err
					}
					err = repository.QuizRepositoryInstance.AddUserChoiceAnswer(ctx, userId, q.Id, int32(choiceId))
					if err != nil {
						return err
					}
				} else {
					correctText, err := repository.QuizRepositoryInstance.GetTextQuestionAnswer(ctx, q.Id)
					if err != nil {
						return err
					}
					if correctText.RightAnswer == answ {
						rightAnswers += 1
					}
					err = repository.QuizRepositoryInstance.RemoveUserTextAnswers(ctx, q.Id, userId)
					if err != nil {
						return err
					}
					err = repository.QuizRepositoryInstance.AddUserTextAnswer(ctx, userId, q.Id, answ)
					if err != nil {
						return err
					}
				}
			}
		}

		err = repository.QuizRepositoryInstance.
			UpsertUserScore(ctx, userId, int32(quizId), rightAnswers/amountOfQuestions, time.Now().UTC())
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/quiz")
}

func QuizParticipationFormGetHandler(c *gin.Context) {
	var userId int32
	data, ok := c.Get("sessionData")
	if ok {
		if sessionData, ok := data.(*middleware.SessionData); ok {
			userId = sessionData.UserId
		}
	}

	i, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := int32(i)

	ctx := context.Background()

	var quiz Quiz
	err = repository.TransactionManager.Run(ctx, func(ctx context.Context) error {

		partTime, err := repository.QuizRepositoryInstance.GetLastParticipationTime(ctx, userId)
		if _, ok := err.(*apperrors.ErrNotFound); ok {
			_, err := repository.QuizRepositoryInstance.AddParticipationTime(ctx, userId, id, time.Now().UTC())
			if err != nil {
				return err
			}
		} else if err == nil {
			if partTime.FinishedAt == nil && partTime.QuizId != id {
				err := repository.QuizRepositoryInstance.DeleteParticipationTime(ctx, partTime.Id)
				if err != nil {
					return err
				}
				_, err = repository.QuizRepositoryInstance.AddParticipationTime(ctx, userId, id, time.Now().UTC())
				if err != nil {
					return err
				}
			} else if partTime.FinishedAt != nil {
				_, err = repository.QuizRepositoryInstance.AddParticipationTime(ctx, userId, id, time.Now().UTC())
				if err != nil {
					return err
				}
			}
		} else {
			return err
		}

		quizModel, err := repository.QuizRepositoryInstance.
			GetQuiz(ctx, id)
		if err != nil {
			return err
		}
		quiz.Id = quizModel.Id
		quiz.Title = quizModel.Title
		quiz.Description = *quizModel.Description

		questionModels, err := repository.QuizRepositoryInstance.
			GetQuizQuestions(ctx, id)
		if err != nil {
			return err
		}
		quiz.Questions = make([]Question, len(questionModels))

		for i, v := range questionModels {
			quiz.Questions[i].Id = v.Id
			quiz.Questions[i].Text = v.QuestionText
			quiz.Questions[i].Type = v.QuestionType

			if v.QuestionType != "text" {
				choiceModels, err := repository.QuizRepositoryInstance.
					GetChoices(ctx, v.Id)
				if err != nil {
					return err
				}
				quiz.Questions[i].Choices = make([]Choice, len(choiceModels))
				for j, choice := range choiceModels {
					quiz.Questions[i].Choices[j].Id = choice.Id
					quiz.Questions[i].Choices[j].Text = choice.ChoiceText
					quiz.Questions[i].Choices[j].QuestionId = v.Id
				}
			}
		}

		return nil
	})
	if err != nil {
		println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	baseHInterface, _ := c.Get("BaseH")
	baseH, _ := baseHInterface.(*gin.H)
	c.HTML(http.StatusOK, "quiz_participation.html", utility.MergeMaps(*baseH, gin.H{
		"title": "Quizzes",
		"quiz":  quiz}))
}

func QuizIndexGetHandler(c *gin.Context) {
	ctx := context.Background()
	categories, err := repository.QuizRepositoryInstance.GetAllCategories(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	categoryMap := make(map[int32]string)
	for _, v := range categories {
		categoryMap[v.Id] = v.Name
	}

	categoryId := int32(0)
	if c.Query("category_id") != "" {
		i, err := strconv.ParseInt(c.Query("category_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		categoryId = int32(i)
	}

	quizzes, err := repository.QuizRepositoryInstance.GetAllQuizzes(ctx, categoryId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	frontQuizzes := make([]Quiz, 0, len(quizzes))
	quizIds := make([]int32, 0, len(quizzes))
	quizMap := make(map[int32]*Quiz)
	for _, q := range quizzes {
		quizIds = append(quizIds, q.Id)
		frontQuizzes = append(frontQuizzes, Quiz{
			Id:           q.Id,
			Title:        q.Title,
			Description:  *q.Description,
			AverageTime:  "00:00:00",
			AverageScore: "0%",
			Categories:   make([]string, 0),
		})
		quizMap[q.Id] = &frontQuizzes[len(frontQuizzes)-1]
	}

	rQuizIds, rCategoryIds, err := repository.QuizRepositoryInstance.
		GetCategoriesPairs(ctx, quizIds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, qId := range rQuizIds {
		quizMap[qId].Categories = append(quizMap[qId].Categories, categoryMap[rCategoryIds[i]])
	}

	quizStatistics, err := repository.QuizRepositoryInstance.GetQuizStatistics(ctx, quizIds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, s := range quizStatistics {
		qp := quizMap[s.QuizId]
		qp.AverageScore = fmt.Sprintf("%.2f%%", s.AverageScore*100)
		qp.AverageTime = *s.AverageCompletionTime
		qp.TotalAttempts = int32(s.TotalAttempts)
	}

	baseHInterface, _ := c.Get("BaseH")
	baseH, _ := baseHInterface.(*gin.H)
	c.HTML(http.StatusOK, "quiz_list.html", utility.MergeMaps(*baseH, gin.H{
		"title":            "Quizzes",
		"categories":       categories,
		"current_category": categoryId,
		"quizzes":          frontQuizzes}))
}

func QuizResultGetHandler(c *gin.Context) {
	var (
		userId int32
		quizId int32
	)
	data, ok := c.Get("sessionData")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	if sessionData, ok := data.(*middleware.SessionData); ok {
		userId = sessionData.UserId
	}
	i, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quizId = int32(i)

	ctx := context.Background()
	userScore, err := repository.QuizRepositoryInstance.
		GetUserScore(ctx, userId, quizId)
	if _, ok := err.(*apperrors.ErrNotFound); ok {
		c.Redirect(http.StatusFound, "/quiz")
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quizModel, err := repository.QuizRepositoryInstance.
		GetQuiz(ctx, quizId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quizPartModel, err := repository.QuizRepositoryInstance.
		GetQuizParticipationTime(ctx, userId, quizId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quizResult := QuizResult{
		Title:       quizModel.Title,
		Description: *quizModel.Description,
		Score:       fmt.Sprintf("%.2f%%", userScore.Score*100),
		Time:        formatDuration(quizPartModel.FinishedAt.Sub(quizPartModel.StartedAt)),
	}

	questionModels, err := repository.QuizRepositoryInstance.
		GetQuizQuestions(ctx, quizId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quizResult.Questions = make([]AnsweredQuestion, len(questionModels))

	for i, v := range questionModels {
		quizResult.Questions[i].Text = v.QuestionText
		quizResult.Questions[i].Type = v.QuestionType

		if v.QuestionType != "text" {
			correctChoice, err := repository.QuizRepositoryInstance.GetCorrectChoice(ctx, v.Id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			userChoice, err := repository.QuizRepositoryInstance.GetUserChoiceAnswer(ctx, userId, v.Id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			choice, err := repository.QuizRepositoryInstance.GetChoice(ctx, *userChoice.ChoiceId)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			quizResult.Questions[i].IsCorrect = correctChoice.Id == *userChoice.ChoiceId
			quizResult.Questions[i].RightAnswer = correctChoice.ChoiceText
			quizResult.Questions[i].UserAnswer = choice.ChoiceText
		} else {
			correctText, err := repository.QuizRepositoryInstance.GetTextQuestionAnswer(ctx, v.Id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			userText, err := repository.QuizRepositoryInstance.GetUserTextAnswer(ctx, userId, v.Id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			quizResult.Questions[i].IsCorrect = *userText.TextAnswer == correctText.RightAnswer
			quizResult.Questions[i].RightAnswer = correctText.RightAnswer
			quizResult.Questions[i].UserAnswer = *userText.TextAnswer
		}
	}

	baseHInterface, _ := c.Get("BaseH")
	baseH, _ := baseHInterface.(*gin.H)
	c.HTML(http.StatusOK, "quiz_my_stats.html", utility.MergeMaps(*baseH, gin.H{
		"title": "Quiz Result",
		"quiz":  quizResult}))
}

func QuizDeletePostHandler(c *gin.Context) {
	i, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := int32(i)

	ctx := context.Background()
	err = repository.QuizRepositoryInstance.DeleteQuiz(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/quiz")
}
