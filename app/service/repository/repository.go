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
	res, err := repository.DB.Query(`SELECT * FROM (SELECT COUNT(*) FROM users) AS usr, ` +
		`(SELECT COUNT(*) FROM forums) AS frm, (SELECT COUNT(*) FROM threads) AS thrd, ` +
		`(SELECT COUNT(*) FROM posts) AS pst`)
	if err != nil {
		return nil
	}
	defer res.Close()
	if !res.Next() {
		return nil
	}
	var stat models.Status
	err = res.Scan(&stat.User, &stat.Forum, &stat.Thread, &stat.Post)
	if err != nil {
		return nil
	}
	return &stat
}

func (repository *Repository) Clear() *models.Error {
	res, err := repository.DB.Query(`TRUNCATE TABLE users, forums, user_forum, threads, posts, votes CASCADE`)
	if err != nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	defer res.Close()
	return nil
}
