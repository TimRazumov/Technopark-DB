package repository

import (
	"net/http"

	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/service"

	"github.com/jackc/pgx"
)

type Repository struct {
	DB *pgx.ConnPool
}

func CreateRepository(db *pgx.ConnPool) service.Repository {
	return &Repository{DB: db}
}

func (repository *Repository) Get() *models.Status {
	return nil
}

func (repository *Repository) Clear() *models.Error {
	res, err := repository.DB.Query("TRUNCATE TABLE users, forums users_forum CASCADE")
	if err != nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	defer res.Close()
	return nil
}
