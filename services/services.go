package services

import (
	"tyre-match-backend/db/repositories"
	"tyre-match-backend/pkg/dependencies"
)

type Services struct {
	TyreModel       *TyreModelService
	TyreImpression  *TyreImpressionService
	ImageProcessing *ImageProcessingService
	File            *FileService
}

func NewServices(repos *repositories.Repos, dependencies *dependencies.DependencyService) *Services {
	file := NewFileService(dependencies.GetFileStore())
	imageProcessingService := NewImageProcessingService(repos.TyreImpression, file, repos.File, dependencies.GetFileStore())

	return &Services{
		TyreImpression:  NewTyreImpressionService(repos.TyreImpression, repos.File, dependencies.GetFileStore(), NewUploadedFileValidator(), imageProcessingService),
		TyreModel:       NewTyreModelService(repos.TyreModel, repos.File, dependencies.GetFileStore(), NewUploadedFileValidator()),
		ImageProcessing: imageProcessingService,
		File:            file,
	}
}
