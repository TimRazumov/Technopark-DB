package thread

import "github.com/TimRazumov/Technopark-DB/app/models"

type UseCase interface {
	Create(thrd *models.Thread) *models.Error
	GetByID(id int) *models.Thread
	GetBySlug(slug string) *models.Thread
	Update(newThrd *models.Thread) *models.Error
}
