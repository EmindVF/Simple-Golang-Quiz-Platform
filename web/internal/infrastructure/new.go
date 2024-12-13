package infrastructure

import (
	"quiz_platform/internal/database"
	"quiz_platform/internal/handler/repository"
	"quiz_platform/internal/infrastructure/repositories"
)

func NewSqlUserRepository(db database.SqlDatabaseProvider) repository.UserRepository {
	return &repositories.SqlUserRepository{
		DBProvider: db,
	}
}

func NewSqlNewsRepository(db database.SqlDatabaseProvider) repository.NewsRepository {
	return &repositories.SqlNewsRepository{
		DBProvider: db,
	}
}

func NewSqlActionsRepository(db database.SqlDatabaseProvider) repository.ActionsRepository {
	return &repositories.SqlActionsRepository{
		DBProvider: db,
	}
}

func NewSqlRoleRepository(db database.SqlDatabaseProvider) repository.RoleRepository {
	return &repositories.SqlRoleRepository{
		DBProvider: db,
	}
}

func NewSqlQuizRepository(db database.SqlDatabaseProvider) repository.QuizRepository {
	return &repositories.SqlQuizRepository{
		DBProvider: db,
	}
}
