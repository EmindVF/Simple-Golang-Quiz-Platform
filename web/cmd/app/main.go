package main

import (
	"fmt"
	"html/template"

	"quiz_platform/internal/infrastructure"
	"quiz_platform/internal/middleware"
	"quiz_platform/internal/misc/config"
	"quiz_platform/internal/misc/formatters"
	"quiz_platform/internal/misc/logger"
	"quiz_platform/internal/misc/templates"
	"quiz_platform/internal/misc/transaction"
	"quiz_platform/internal/models"

	"quiz_platform/internal/handler/actions"
	"quiz_platform/internal/handler/auth"
	"quiz_platform/internal/handler/misc"
	"quiz_platform/internal/handler/news"
	"quiz_platform/internal/handler/quiz"
	"quiz_platform/internal/handler/repository"
	"quiz_platform/internal/handler/users"

	"quiz_platform/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Init config
	err := config.ReadGlobalConfig("./config/config.json")
	if err != nil {
		panic(fmt.Sprintf("cannot read config: %v", err.Error()))
	}

	// Init logger
	err = logger.InitLogger("./logs")
	if err != nil {
		panic(fmt.Sprintf("cannot init config: %v", err.Error()))
	}
	defer logger.CleanLogger()

	// Init database provider
	sqlProvider :=
		database.NewPostgresSqlDatabaseProvider()

	// Init transaction manager
	repository.TransactionManager =
		transaction.NewTransactionManager(sqlProvider.GetDb())

	// Init repositories
	repository.UserRepositoryInstance =
		infrastructure.NewSqlUserRepository(sqlProvider)

	repository.NewsRepositoryInstance =
		infrastructure.NewSqlNewsRepository(sqlProvider)

	repository.ActionsRepositoryInstance =
		infrastructure.NewSqlActionsRepository(sqlProvider)

	repository.RoleRepositoryInstance =
		infrastructure.NewSqlRoleRepository(sqlProvider)

	repository.QuizRepositoryInstance =
		infrastructure.NewSqlQuizRepository(sqlProvider)

	r := gin.Default()

	// Set funcs
	r.SetFuncMap(template.FuncMap{
		"formatDate": formatters.FormatDate,
		"bitwiseAnd": formatters.BitwiseAnd,
	})

	// Init templates
	files, err := templates.LoadTemplates("templates")
	if err != nil {
		panic(fmt.Sprintf("cannot load templates: %v", err.Error()))
	}
	r.LoadHTMLFiles(files...)
	r.Static("/static", "./static")
	r.Use(middleware.GetTokenMiddleware())

	// Set routes
	r.GET("/", misc.IndexHandler)

	// Auth
	r.GET("/register", auth.RegisterGetHandler)
	r.POST("/register", auth.RegisterPostHandler)
	r.GET("/login", auth.LoginGetHandler)
	r.POST("/login", auth.LoginPostHandler)
	r.POST("/logout", middleware.RequirePermissionMiddleware(0), misc.IndexHandler)

	// News
	r.GET("/news", news.NewsListGetHandler)
	r.GET("/news/:id", news.NewsViewGetHandler)
	r.GET("/news/new", middleware.RequirePermissionMiddleware(models.MANAGE_NEWS_PERM), news.NewsCreateFormGetHandler)
	r.POST("/news", middleware.RequirePermissionMiddleware(models.MANAGE_NEWS_PERM), news.NewsCreatePostHandler)
	r.GET("/news/:id/edit", middleware.RequirePermissionMiddleware(models.MANAGE_NEWS_PERM), news.NewsEditFormGetHandler)
	r.POST("/news/:id", middleware.RequirePermissionMiddleware(models.MANAGE_NEWS_PERM), news.NewsEditPostHandler)
	r.POST("/news/:id/delete", middleware.RequirePermissionMiddleware(models.MANAGE_NEWS_PERM), news.NewsDeletePostHandler)

	// Actions
	r.GET("/actions", middleware.RequirePermissionMiddleware(models.VIEW_ACTIONS_PERM), actions.ActionsListGetHandler)

	// User Management
	r.GET("/users", middleware.RequirePermissionMiddleware(models.MANAGE_USERS_PERM), users.UsersListGetHandler)
	r.GET("/users/:id/edit", middleware.RequirePermissionMiddleware(models.MANAGE_USERS_PERM), users.UserEditFormGetHandler)
	r.POST("/users/:id/edit", middleware.RequirePermissionMiddleware(models.MANAGE_USERS_PERM), users.UserEditFormPostHandler)
	r.POST("/users/:id/delete", middleware.RequirePermissionMiddleware(models.MANAGE_USERS_PERM), users.UserDeletePostHandler)

	// Quiz
	r.GET("/quiz", middleware.RequirePermissionMiddleware(0), quiz.QuizIndexGetHandler)
	r.GET("/quiz/create", middleware.RequirePermissionMiddleware(0), quiz.QuizCreateFormGetHandler)
	r.POST("/quiz/create", middleware.RequirePermissionMiddleware(0), quiz.QuizCreatePostHandler)
	r.GET("/quiz/:id/participate", middleware.RequirePermissionMiddleware(0), quiz.QuizParticipationFormGetHandler)
	r.POST("/quiz/participate", middleware.RequirePermissionMiddleware(0), quiz.QuizParticipationPostHandler)
	r.POST("/quiz/:id/delete", middleware.RequirePermissionMiddleware(models.MANAGE_QUIZZES_PERM), quiz.QuizDeletePostHandler)
	r.GET("/quiz/:id/result", middleware.RequirePermissionMiddleware(0), quiz.QuizResultGetHandler)

	// Run server
	r.Run(fmt.Sprintf(":%d", config.GlobalConfig.App.Port))
}
