package forum

import "github.com/TimRazumov/Technopark-DB/app/models"

type Repository interface {
	Create(frm models.Forum) *models.Error
	GetBySlug(slug string) *models.Forum
	GetUsersBySlug(slug string, queryString models.QueryString) *models.Users
	GetThreadsBySlug(slug string, queryString models.QueryString) *models.Threads
}
