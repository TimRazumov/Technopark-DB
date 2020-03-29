package post

import (
	"github.com/TimRazumov/Technopark-DB/app/models"
)

type UseCase interface {
	Create(thrdKey string, pst *[]models.Post) *models.Error
}
