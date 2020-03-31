package repository

import (
	"net/http"
	"time"

	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/post"

	"github.com/jackc/pgx"
	"github.com/lib/pq"
)

type Repository struct {
	DB *pgx.ConnPool
}

func CreateRepository(db *pgx.ConnPool) post.Repository {
	return &Repository{DB: db}
}

func (repository *Repository) Create(thrd models.Thread, posts *[]models.Post) *models.Error {
	if posts == nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	transact, err := repository.DB.Begin()
	if err != nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	defer transact.Rollback()
	parents := make(map[int]models.Post)
	for _, pst := range *posts {
		if _, has := parents[pst.Parent]; !has && pst.Parent != 0 {
			parentPost := repository.GetByID(pst.Parent)
			if parentPost == nil || parentPost.Thread != thrd.ID {
				return models.CreateConflictPost()
			}
			parents[pst.Parent] = *parentPost
		}
	}
	postRows, err := transact.Query(`SELECT nextval(pg_get_serial_sequence('posts', 'id')) FROM generate_series(1, $1)`, len(*posts))
	if err != nil {
		return &models.Error{Code: http.StatusNotFound}
	}
	defer postRows.Close()
	var postsID []int
	for postRows.Next() {
		var id int
		err = postRows.Scan(&id)
		if err != nil {
			return &models.Error{Code: http.StatusInternalServerError}
		}
		postsID = append(postsID, id)
	}
	if len(postsID) != len(*posts) {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	timeNow := time.Now()
	for idx, pst := range *posts {
		pst.ID = postsID[idx]
		pst.Forum = thrd.Forum
		pst.Thread = thrd.ID
		pst.Created = timeNow
		pst.Path = append(parents[pst.Parent].Path, int64(postsID[idx]))
		(*posts)[idx] = pst
		res, err := transact.Exec(`INSERT INTO posts (id, parent, author, message, forum, thread, created, path) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			pst.ID, pst.Parent, pst.Author, pst.Message, pst.Forum, pst.Thread, pst.Created, pq.Array(pst.Path))
		if err != nil || res.RowsAffected() == 0 {
			return models.CreateNotFoundAuthorPost(pst.Author)
		}
		_, err = transact.Exec(`INSERT INTO user_forum (nickname, slug) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			pst.Author, thrd.Forum)
		if err != nil {
			return &models.Error{Code: http.StatusInternalServerError}
		}
	}
	res, err := transact.Exec(`UPDATE forums SET posts = posts + $1 WHERE slug = $2`, len(*posts), thrd.Forum)
	if err != nil || res.RowsAffected() == 0 || transact.Commit() != nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	return nil
}

func (repository *Repository) GetByID(id int) *models.Post {
	res, err := repository.DB.Query(`SELECT * FROM posts WHERE id = $1`, id)
	if err != nil {
		return nil
	}
	defer res.Close()
	for !res.Next() {
		return nil
	}
	var pst models.Post
	err = res.Scan(&pst.ID, &pst.Parent, &pst.Author, &pst.Message, &pst.IsEdited,
		&pst.Forum, &pst.Thread, &pst.Created, pq.Array(&pst.Path))
	if err != nil {
		return nil
	}
	return &pst
}

func (repository *Repository) Update(pst *models.Post) *models.Error {
	if pst == nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	realPst := repository.GetByID(pst.ID)
	if realPst == nil {
		return &models.Error{Code: http.StatusNotFound}
	}
	if pst.Message == "" || pst.Message == realPst.Message {
		*pst = *realPst
		return nil
	}
	message := pst.Message
	*pst = *realPst
	pst.Message = message
	pst.IsEdited = true
	_, err := repository.DB.Exec("UPDATE posts SET message = $1, is_edited = $2 WHERE id = $3",
		pst.Message, pst.IsEdited, pst.ID)
	if err != nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	return nil
}
