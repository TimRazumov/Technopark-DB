package repository

import (
	"net/http"

	"github.com/TimRazumov/Technopark-DB/app/forum"
	"github.com/TimRazumov/Technopark-DB/app/models"

	"github.com/jackc/pgx"
)

type Repository struct {
	DB *pgx.ConnPool
}

func CreateRepository(db *pgx.ConnPool) forum.Repository {
	return &Repository{DB: db}
}

func (repository *Repository) Create(frm models.Forum) *models.Error {
}

func (repository *Repository) GetBySlug(slug string) *models.Forum {
}
