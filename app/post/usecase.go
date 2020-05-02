package post

import (
	"github.com/TimRazumov/Technopark-DB/app/models"
)

type UseCase interface {
	Create(thrdKey string, pst *models.Posts) *models.Error
	GetByID(id int, options models.Related) *models.PostGet
	Update(pst *models.Post) *models.Error
}
