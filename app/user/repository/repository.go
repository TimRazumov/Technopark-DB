package repository

import (
	"net/http"

	"github.com/TimRazumov/Technopark-DB/app/models"
	"github.com/TimRazumov/Technopark-DB/app/user"

	"github.com/jackc/pgx"
)

type Repository struct {
	DB *pgx.ConnPool
}

func CreateRepository(db *pgx.ConnPool) user.Repository {
	return &Repository{DB: db}
}

func (repository *Repository) Create(usr models.User) *models.Error {
	res, err := repository.DB.Exec(`INSERT INTO users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4)`,
		usr.NickName, usr.FullName, usr.Email, usr.About)
	if err != nil || res.RowsAffected() == 0 {
		return &models.Error{Code: http.StatusConflict}
	}
	return nil
}

func (repository *Repository) GetByNickName(nickname string) *models.User {
	res, err := repository.DB.Query(`SELECT * FROM users WHERE nickname = $1`, nickname)
	if err != nil {
		return nil
	}
	defer res.Close()
	if !res.Next() {
		return nil
	}
	var usr models.User
	err = res.Scan(&usr.NickName, &usr.FullName, &usr.Email, &usr.About)
	if err != nil {
		return nil
	}
	return &usr
}

func (repository *Repository) GetByEmail(email string) *models.User {
	res, err := repository.DB.Query(`SELECT * FROM users WHERE email = $1`, email)
	if err != nil {
		return nil
	}
	defer res.Close()
	if !res.Next() {
		return nil
	}
	var usr models.User
	err = res.Scan(&usr.NickName, &usr.FullName, &usr.Email, &usr.About)
	if err != nil {
		return nil
	}
	return &usr
}

func (repository *Repository) Update(newUsr *models.User) *models.Error {
	if newUsr == nil {
		return &models.Error{Code: http.StatusInternalServerError}
	}
	realUsr := repository.GetByNickName(newUsr.NickName)
	if realUsr == nil {
		return models.CreateNotFoundUser(newUsr.NickName)
	}
	if newUsr.FullName == "" && newUsr.Email == "" && newUsr.About == "" {
		*newUsr = *realUsr
		return nil
	}
	if newUsr.FullName == "" {
		newUsr.FullName = realUsr.FullName
	}
	if newUsr.Email == "" {
		newUsr.Email = realUsr.Email
	}
	if newUsr.About == "" {
		newUsr.About = realUsr.About
	}
	res, err := repository.DB.Exec("UPDATE users SET fullname = $1, email = $2, about = $3 WHERE nickname = $4",
		newUsr.FullName, newUsr.Email, newUsr.About, newUsr.NickName)
	if err != nil {
		conflictUsr := repository.GetByEmail(newUsr.Email)
		if conflictUsr != nil {
			return models.CreateConflictUser(conflictUsr.NickName)
		}
		return &models.Error{Code: http.StatusInternalServerError}
	}
	if res.RowsAffected() == 0 {
		return &models.Error{Code: http.StatusNotFound}
	}
	return nil
}
