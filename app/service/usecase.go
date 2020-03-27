package service

import "github.com/TimRazumov/Technopark-DB/app/models"

type UseCase interface {
	Get() *models.Status
	Clear() *models.Error
}
