package services

import "example.com/simple-api/models"

type MailService interface {
	CreateMail(*models.Mail) error
	GetMail(*int) (*models.Mail, error)
	GetAll() ([]*models.Mail, error)
	UpdateMail(*int, *models.Mail) error
	DeleteMail(*int) error
}
