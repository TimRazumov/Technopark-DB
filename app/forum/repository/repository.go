package repository

import (
	"github.com/TimRazumov/Technopark-DB/app/forum"
	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"net/http"
	"strconv"
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

func (repository *Repository) GetUsersBySlug(slug string, queryString models.QueryString) *[]models.User {
	sinceCond := ``
	if queryString.Since != "" {
		sinceCond = ` AND nickname `
		if queryString.Desc {
			sinceCond += `< `
		} else {
			sinceCond += `> `
		}
		sinceCond += `'` + queryString.Since + `'`
	}
	sortCond := ` ORDER BY nickname`
	if queryString.Desc {
		sortCond += ` DESC`
	}
	limitCond := ``
	if queryString.Limit > 0 {
		limitCond += ` LIMIT ` + strconv.Itoa(queryString.Limit)
	}
	res, err := repository.DB.Query(`SELECT nickname, fullname, email, about FROM users`+
		` JOIN user_forum USING(nickname) WHERE slug = $1`+sinceCond+sortCond+limitCond, slug)
	if err != nil {
		return nil
	}
	defer res.Close()
	usrs := make([]models.User, 0)
	for res.Next() {
		var tmpUser models.User
		err = res.Scan(&tmpUser.NickName, &tmpUser.FullName, &tmpUser.Email, &tmpUser.About)
		if err != nil {
			return nil
		}
		usrs = append(usrs, tmpUser)
	}
	return &usrs
}

func (repository *Repository) GetThreadsBySlug(slug string, queryString models.QueryString) *[]models.Thread {
	sinceCond := ``
	if queryString.Since != "" {
		sinceCond = ` AND created `
		if queryString.Desc {
			sinceCond += `<= `
		} else {
			sinceCond += `>= `
		}
		sinceCond += `'` + queryString.Since + `'`
	}
	sortCond := ` ORDER BY created`
	if queryString.Desc {
		sortCond += ` DESC`
	}
	limitCond := ``
	if queryString.Limit > 0 {
		limitCond += ` LIMIT ` + strconv.Itoa(queryString.Limit)
	}
	res, err := repository.DB.Query(`SELECT * FROM threads WHERE forum = $1`+sinceCond+sortCond+limitCond, slug)
	if err != nil {
		return nil
	}
	defer res.Close()
	thrd := make([]models.Thread, 0)
	for res.Next() {
		var tmpThrd models.Thread
		nullSlug := &pgtype.Varchar{}
		err = res.Scan(&tmpThrd.ID, &tmpThrd.Title, &tmpThrd.Author, &tmpThrd.Forum,
			&tmpThrd.Message, &tmpThrd.Votes, nullSlug, &tmpThrd.Created)
		if err != nil {
			return nil
		}
		tmpThrd.Slug = nullSlug.String
		thrd = append(thrd, tmpThrd)
	}
	return &thrd
}
