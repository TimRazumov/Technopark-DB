package thread

import "github.com/TimRazumov/Technopark-DB/app/models"

type Repository interface {
	Create(thrd *models.Thread) *models.Error
	GetByID(id int) *models.Thread
	GetBySlug(slug string) *models.Thread
	Update(newThrd *models.Thread) *models.Error
	UpdateVote(vt models.Vote) *models.Error
	GetPostsByThread(thrd models.Thread, queryString models.QueryString) *models.Posts
}
