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
	res, err := repository.DB.Exec(`INSERT INTO forums (slug, title, usr, posts, threads) VALUES ($1, $2, $3, $4, $5)`,
		frm.Slug, frm.Title, frm.User, frm.Posts, frm.Threads)
	if err != nil || res.RowsAffected() == 0 {
		return &models.Error{Code: http.StatusConflict}
	}
	return nil
}

func (repository *Repository) GetBySlug(slug string) *models.Forum {
	res, err := repository.DB.Query(`SELECT * FROM forums WHERE slug = $1`, slug)
	if err != nil {
		return nil
	}
	defer res.Close()
	if !res.Next() {
		return nil
	}
	var frm models.Forum
	err = res.Scan(&frm.Slug, &frm.Title, &frm.User, &frm.Posts, &frm.Threads)
	if err != nil {
		return nil
	}
	return &frm
}
