package service

import "github.com/TimRazumov/Technopark-DB/app/models"

type Repository interface {
	Get() *models.Status
	Clear() *models.Error
}
