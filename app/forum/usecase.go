package forum

import "github.com/TimRazumov/Technopark-DB/app/models"

type UseCase interface {
	Create(frm *models.Forum) *models.Error
	GetBySlug(slug string) *models.Forum
	// GetUsers(slug string, query models.PostsRequestQuery) ([]models.User, *models.Error)
}
