package post

import "github.com/TimRazumov/Technopark-DB/app/models"

type Repository interface {
	Create(thrd models.Thread, posts *models.Posts) *models.Error
	GetByID(id int) *models.Post
	Update(pst *models.Post) *models.Error
}
