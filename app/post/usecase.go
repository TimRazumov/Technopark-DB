package post

import (
	"github.com/TimRazumov/Technopark-DB/app/models"
)

type UseCase interface {
	Create(postKey string, pst *[]models.Post) *models.Error
}
