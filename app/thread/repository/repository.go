package repository

import (
	"net/http"

	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/thread"

	"github.com/jackc/pgx"
)

type Repository struct {
	DB *pgx.ConnPool
}

func CreateRepository(db *pgx.ConnPool) thread.Repository {
	return &Repository{DB: db}
}

func (repository *Repository) Create(thrd *models.Thread) *models.Error {
	if thrd == nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	transact, err := repository.DB.Begin()
	if err != nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	defer transact.Rollback()
	if thrd.Slug != "" {
		err = transact.QueryRow(`INSERT INTO threads (title, author, forum, message, slug, created) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			thrd.Title, thrd.Author, thrd.Forum, thrd.Message, thrd.Slug, thrd.Created).Scan(&thrd.ID)
	} else {
		err = transact.QueryRow(`INSERT INTO threads (title, author, forum, message, created) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			thrd.Title, thrd.Author, thrd.Forum, thrd.Message, thrd.Created).Scan(&thrd.ID)
	}
	if err != nil {
		return &models.Error{Code: http.StatusConflict}
	}
	err = transact.Commit()
	if err != nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	return nil
}

func (repository *Repository) GetByID(id int) *models.Thread {
	res, err := repository.DB.Query(`SELECT * FROM threads WHERE id = $1`, id)
	if err != nil {
		return nil
	}
	defer res.Close()
	if !res.Next() {
		return nil
	}
	var thrd models.Thread
	err = res.Scan(&thrd.ID, &thrd.Title, &thrd.Author, &thrd.Forum, &thrd.Message, &thrd.Votes, &thrd.Slug, &thrd.Created)
	if err != nil {
		return nil
	}
	return &thrd
}

func (repository *Repository) GetBySlug(slug string) *models.Thread {
	if slug == "" {
		return nil
	}
	res, err := repository.DB.Query(`SELECT * FROM threads WHERE slug = $1`, slug)
	if err != nil {
		return nil
	}
	defer res.Close()
	if !res.Next() {
		return nil
	}
	var thrd models.Thread
	err = res.Scan(&thrd.ID, &thrd.Title, &thrd.Author, &thrd.Forum, &thrd.Message, &thrd.Votes, &thrd.Slug, &thrd.Created)
	if err != nil {
		return nil
	}
	return &thrd
}

func (repository *Repository) Update(newThrd *models.Thread) *models.Error {
	if newThrd == nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	var realThrd *models.Thread
	if newThrd.ID != -1 {
		realThrd = repository.GetByID(newThrd.ID)
	} else {
		realThrd = repository.GetBySlug(newThrd.Slug)
	}
	if realThrd == nil {
		return &models.Error{Code: http.StatusNotFound}
	}
	newTitle := newThrd.Title
	newMessage := newThrd.Message
	*newThrd = *realThrd
	if newTitle == "" && newMessage == "" {
		return nil
	}
	if newTitle != "" {
		newThrd.Title = newTitle
	}
	if newMessage != "" {
		newThrd.Message = newMessage
	}
	res, err := repository.DB.Exec("UPDATE threads SET title = $1, message = $2 WHERE id = $3",
		newThrd.Title, newThrd.Message, newThrd.ID)
	if err != nil || res.RowsAffected() == 0 {
		return &models.Error{Code: http.StatusNotFound}
	}
	return nil
}
