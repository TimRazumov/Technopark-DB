package post

import "github.com/TimRazumov/Technopark-DB/app/models"

type Repository interface {
	Create(thrd models.Thread, posts *[]models.Post) *models.Error
}
