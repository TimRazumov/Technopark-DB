package forum

import "github.com/TimRazumov/Technopark-DB/app/models"

type UseCase interface {
	Create(frm *models.Forum) *models.Error
	GetBySlug(slug string) *models.Forum
	GetUsersBySlug(slug string, queryString models.QueryString) *[]models.User
	GetThreadsBySlug(slug string, queryString models.QueryString) *[]models.Thread
}
