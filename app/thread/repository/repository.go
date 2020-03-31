package repository

import (
	"github.com/lib/pq"
	"net/http"
	"strconv"

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
		return models.CreateNotFoundAuthorPost(newThrd.Slug)
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
	res, err := repository.DB.Exec(`UPDATE threads SET title = $1, message = $2 WHERE id = $3`,
		newThrd.Title, newThrd.Message, newThrd.ID)
	if err != nil || res.RowsAffected() == 0 {
		return &models.Error{Code: http.StatusNotFound}
	}
	return nil
}

func (repository *Repository) UpdateVote(vt models.Vote) *models.Error {
	selectRes, err := repository.DB.Query(`SELECT voice FROM votes WHERE nickname = $1 AND thread = $2`,
		vt.NickName, vt.Thread)
	if err == nil && selectRes.Next() {
		defer selectRes.Close()
		var oldVoice int
		err = selectRes.Scan(&oldVoice)
		if err != nil {
			return &models.Error{Code: http.StatusInternalServerError}
		}
		if oldVoice == vt.Voice {
			return nil
		}
		res, err := repository.DB.Exec(`UPDATE votes SET voice = $1 WHERE nickname = $2 AND thread = $3`,
			vt.Voice, vt.NickName, vt.Thread)
		if err != nil || res.RowsAffected() == 0 {
			return &models.Error{Code: http.StatusInternalServerError}
		}
		vt.Voice *= 2
	} else {
		res, err := repository.DB.Exec(`INSERT INTO votes (nickname, voice, thread) VALUES ($1, $2, $3)`,
			vt.NickName, vt.Voice, vt.Thread)
		if err != nil || res.RowsAffected() == 0 {
			return &models.Error{Code: http.StatusInternalServerError}
		}
	}
	res, err := repository.DB.Exec(`UPDATE threads SET votes = votes + $1 WHERE id = $2`, vt.Voice, vt.Thread)
	if err != nil || res.RowsAffected() == 0 {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	return nil
}

func (repository *Repository) GetPostsByThread(thrd models.Thread, options models.QueryString) *[]models.Post {
	var res *pgx.Rows
	var err error
	if options.Sort == "flat" {
		sinceCond := ""
		if options.Since != "" {
			sinceCond = " AND id "
			if options.Desc {
				sinceCond += "< "
			} else {
				sinceCond += "> "
			}
			sinceCond += options.Since
		}
		descCond := " ORDER BY id"
		if options.Desc {
			descCond += " DESC"
		}
		limitCond := ""
		if options.Limit > 0 {
			limitCond += " LIMIT " + strconv.Itoa(options.Limit)
		}
		res, err = repository.DB.Query(`SELECT * FROM posts WHERE thread = $1`+sinceCond+descCond+limitCond, thrd.ID)
	} else if options.Sort == "tree" {
		sinceCond := ""
		if options.Since != "" {
			sinceCond = " AND path "
			if options.Desc {
				sinceCond += "< "
			} else {
				sinceCond += "> "
			}
			sinceCond += "(SELECT path FROM posts WHERE id = " + options.Since + ")"
		}
		descCond := " ORDER BY"
		if options.Desc {
			descCond += " path DESC, id DESC"
		} else {
			descCond += " path, id"
		}
		limitCond := ""
		if options.Limit > 0 {
			limitCond += " LIMIT " + strconv.Itoa(options.Limit)
		}
		res, err = repository.DB.Query(`SELECT * FROM posts WHERE thread = $1`+sinceCond+descCond+limitCond, thrd.ID)
	} else if options.Sort == "parent_tree" {
		sinceCond := ""
		if options.Since != "" {
			sinceCond = " AND path[1] "
			if options.Desc {
				sinceCond += "< "
			} else {
				sinceCond += "> "
			}
			sinceCond += "(SELECT path[1] FROM posts WHERE id = " + options.Since + ")"
		}
		descCondIn := " ORDER BY id"
		descCondOut := " ORDER BY"
		if options.Desc {
			descCondIn += " DESC"
			descCondOut += " path[1] DESC, path, id"
		} else {
			descCondOut += " path"
		}
		limitCond := ""
		if options.Limit > 0 {
			limitCond += " LIMIT " + strconv.Itoa(options.Limit)
		}
		res, err = repository.DB.Query("SELECT * FROM posts WHERE path[1] IN"+
			" (SELECT id FROM posts WHERE thread = $1 AND parent = 0"+sinceCond+descCondIn+limitCond+")"+descCondOut, thrd.ID)
	} else {
		return nil
	}
	if err != nil {
		return nil
	}
	defer res.Close()
	psts := make([]models.Post, 0)
	for res.Next() {
		var tmpPost models.Post
		err = res.Scan(&tmpPost.ID, &tmpPost.Parent, &tmpPost.Author, &tmpPost.Message, &tmpPost.IsEdited,
			&tmpPost.Forum, &tmpPost.Thread, &tmpPost.Created, pq.Array(&tmpPost.Path))
		if err != nil {
			return nil
		}
		psts = append(psts, tmpPost)
	}
	return &psts
}
