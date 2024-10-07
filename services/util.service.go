package services

import gonanoid "github.com/matoous/go-nanoid/v2"

type UtilService interface {
	GenerateRandomID(length int) string
}

type UtilServiceImpl struct{}

func (service *UtilServiceImpl) GenerateRandomID(length int) string {
	// generate random string id using nanoid package
	id, _ := gonanoid.New(length)
	return id
}
